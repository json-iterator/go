package jsoniter

type arrayLazyAny struct {
	baseAny
	buf       []byte
	iter      *Iterator
	err       error
	cache     []Any
	remaining []byte
}

func (any *arrayLazyAny) parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
		any.iter = iter
	}
	iter.ResetBytes(any.remaining)
	return iter
}

func (any *arrayLazyAny) fillCacheUntil(target int) Any {
	if any.remaining == nil {
		if target >= len(any.cache) {
			return nil
		}
		return any.cache[target]
	}
	i := len(any.cache)
	if target < i {
		return any.cache[target]
	}
	iter := any.parse()
	if (len(any.remaining) == len(any.buf)) {
		iter.head++
		c := iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			element := iter.readAny(iter)
			any.cache = append(any.cache, element)
			if target == 0 {
				any.remaining = iter.buf[iter.head:]
				return element
			}
			i = 1
		} else {
			any.remaining = nil
			return nil
		}
	}
	for iter.nextToken() == ',' {
		element := iter.readAny(iter)
		any.cache = append(any.cache, element)
		if i == target {
			any.remaining = iter.buf[iter.head:]
			return element
		}
		i++
 	}
	any.remaining = nil
	return nil
}

func (any *arrayLazyAny) fillCache() {
	if any.remaining == nil {
		return
	}
	iter := any.parse()
	if len(any.remaining) == len(any.buf) {
		iter.head++
		c := iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			any.cache = append(any.cache, iter.readAny(iter))
		} else {
			any.remaining = nil
			return
		}
	}
	for iter.nextToken() == ',' {
		any.cache = append(any.cache, iter.readAny(iter))
	}
	any.remaining = nil
	return
}

func (any *arrayLazyAny) LastError() error {
	return nil
}

func (any *arrayLazyAny) ToBool() bool {
	return false
}

func (any *arrayLazyAny) ToInt() int {
	return 0
}

func (any *arrayLazyAny) ToInt32() int32 {
	return 0
}

func (any *arrayLazyAny) ToInt64() int64 {
	return 0
}

func (any *arrayLazyAny) ToFloat32() float32 {
	return 0
}

func (any *arrayLazyAny) ToFloat64() float64 {
	return 0
}

func (any *arrayLazyAny) ToString() string {
	return ""
}

func (any *arrayLazyAny) Get(path ...interface{}) Any {
	if len(path) == 0 {
		return any
	}
	if len(path) == 1 {
		idx := path[0].(int)
		return any.fillCacheUntil(idx)
	} else {
		idx := path[0].(int)
		return any.fillCacheUntil(idx).Get(path[1:]...)
	}
}

func (any *arrayLazyAny) Size() int {
	any.fillCache()
	return len(any.cache)
}


func (any *arrayLazyAny) IterateArray() (func() (Any, bool), bool) {
	remaining := any.remaining
	if len(remaining) == len(any.buf) {
		iter := any.parse()
		iter.head++
		c := iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			v := iter.readAny(iter)
			any.cache = append(any.cache, v)
			remaining = iter.buf[iter.head:]
			any.remaining = remaining
		} else {
			remaining = nil
			any.remaining = nil
			return nil, false
		}
	}
	if len(any.cache) == 0 {
		return nil, false
	}
	arr := any.cache
	nextValue := arr[0]
	i := 1
	return func() (Any, bool) {
		value := nextValue
		if i < len(arr) {
			// read from cache
			nextValue = arr[i]
			i++
			return value, true
		} else {
			// read from buffer
			iter := any.iter
			if iter == nil {
				iter = NewIterator()
				any.iter = iter
			}
			iter.ResetBytes(remaining)
			c := iter.nextToken()
			if c == ',' {
				nextValue = iter.readAny(iter)
				any.cache = append(any.cache, nextValue)
				remaining = iter.buf[iter.head:]
				any.remaining = remaining
				return value, true
			} else {
				remaining = nil
				any.remaining = nil
				return value, false
			}
		}
	}, true
}
