package jsoniter

import (
	"errors"
	"fmt"
)

const (
	// Missed represent not follow the path
	Missed int = 0

	// FollowThePath represent follow the path
	FollowThePath int = 1
)

var (
	// ErrShouldBeArray should be a array
	ErrShouldBeArray = errors.New("Should be a Array")

	// ErrShouldBeObject should be a Object
	ErrShouldBeObject = errors.New("Should be a Object")
)

// Extract represent a Path with several KeyReader in one object followed by the path
type Extract struct {
	Path       Path
	KeyReaders []KeyReader
	index      int
	status     []int
}

// NewExtract constructor
func NewExtract(path Path, readers ...KeyReader) *Extract {
	return &Extract{Path: path, KeyReaders: readers, index: -1, status: make([]int, len(path))}
}

func (e *Extract) finished() bool {
	for _, reader := range e.KeyReaders {
		if !reader.HasRead() && !reader.MaybeNull() {
			return false
		}
	}
	return true
}

func (e *Extract) scanArray(iter *Iterator, depth, index int) (bool, bool, int, error) {
	if depth < len(e.status) && e.status[depth] == FollowThePath {
		return false, false, 0, nil
	}

	if depth == e.index+1 && depth < len(e.Path) {
		if e.Path[depth].Type() == TypeArrayIndex && e.Path[depth].Index() == index {
			e.index++
			e.status[depth] = FollowThePath

			nextValue := -1
			switch iter.nextToken() {
			case '{':
				nextValue = TypeObject
			case '[':
				nextValue = TypeArray
			default:
				return false, false, 0, ErrInvalidPath
			}

			iter.unreadByte()
			return false, true, nextValue, nil
		}
		return false, false, 0, nil
	}

	allFollowThePath := true
	if depth < len(e.status) {
		for _, status := range e.status[:depth] {
			if status != FollowThePath {
				allFollowThePath = false
				break
			}
		}
	}

	if (allFollowThePath || len(e.Path) == 0) && e.index == len(e.Path)-1 && depth == len(e.Path) {
		read, err := e.extract(iter, INT(index))
		return read, false, 0, err
	}

	return false, false, 0, nil
}

// {"type":"xxx", "payload":{"a":true,"b":3,"c":{"d":"e"}},"f":"g"}
// ["payload", "c"]
// ["payload"]
// bool: readValue, bool: FollowThePath, bool:array or object, error
func (e *Extract) scanObject(iter *Iterator, field string, depth int) (bool, bool, int, error) {
	if depth < len(e.status) && e.status[depth] == FollowThePath {
		return false, false, 0, nil
	}

	// 如果depth在路径上
	if depth == e.index+1 && depth < len(e.Path) {
		if e.Path[depth].Type() != TypeArrayIndex && e.Path[depth].String() == field {
			e.index++
			e.status[depth] = FollowThePath

			return false, true, e.Path[depth].Type(), nil
		}
		return false, false, 0, nil
	}

	allFollowThePath := true
	if depth < len(e.status) {
		for _, status := range e.status[:depth] {
			if status != FollowThePath {
				allFollowThePath = false
				break
			}
		}
	}

	// 最后一个路径,depth是最后一个路径的object
	if (allFollowThePath || len(e.Path) == 0) && e.index == len(e.Path)-1 && depth == len(e.Path) {
		var read bool
		var err error

		if field == "" {
			read, err = e.extract(iter, nil)
		} else {
			read, err = e.extract(iter, STRING(field))
		}

		return read, false, 0, err
	}

	return false, false, 0, nil
}

// bool : 读取value, error: 错误
func (e *Extract) extract(iter *Iterator, field IKey) (bool, error) {
	if field != nil {
		for _, reader := range e.KeyReaders {
			if success, err := reader.Read(iter, field); err != nil {
				return false, err
			} else if success {
				return true, nil
			}
		}
		return false, nil
	}

	left := 0
	for _, reader := range e.KeyReaders {
		if !reader.HasRead() && !reader.MaybeNull() {
			left++
		}
	}

	if left > 0 {
		return false, fmt.Errorf("%d not read", left)
	}

	return false, nil
}

// ExtractMany will extract several path -------------------------------------------
// {"type":"XXX","a":"b",c:[1,3,5,8,{"d":"e","f":[23,"big"]}]}
type ExtractMany struct {
	Iter     *Iterator
	Extracts []*Extract
}

// NewExtractMany constructor
func NewExtractMany(iter *Iterator, extracts ...*Extract) *ExtractMany {
	return &ExtractMany{Iter: iter, Extracts: extracts}
}

// ExtractObject API
func (em *ExtractMany) ExtractObject() error {
	return em.extractObject(0)
}

func (em *ExtractMany) extractObject(depth int) error {
	for {
		field := em.Iter.ReadObject()
		if em.Iter.Error != nil {
			return ErrInvalidPath
		}

		valueRead, isFollowThePath, t, err := em.scanObject(field, depth)
		if err != nil {
			return err
		}

		if isFollowThePath {
			switch t {
			case TypeArray:
				err = em.extractArray(depth + 1)
			case TypeObject:
				err = em.extractObject(depth + 1)
			default:
				panic("should not happen")
			}

			if err != nil {
				return err
			}
		} else if field != "" && !valueRead {
			em.Iter.Skip()
		}

		if field == "" {
			break
		}
	}

	return nil
}

func (em *ExtractMany) scanObject(field string, depth int) (bool, bool, int, error) {
	valueRead := false
	isFollowThePath := false
	typeOfPath := -1

	for _, extract := range em.Extracts {
		if extract.finished() {
			continue
		}

		vRead, followed, t, err := extract.scanObject(em.Iter, field, depth)
		if err != nil {
			return false, false, 0, err
		}
		if vRead {
			valueRead = true
		}

		if followed {
			isFollowThePath = true
			typeOfPath = t
		}
	}

	return valueRead, isFollowThePath, typeOfPath, nil
}

// ExtractArray API
func (em *ExtractMany) ExtractArray() error {
	return em.extractArray(0)
}

func (em *ExtractMany) extractArray(depth int) error {
	index := 0
	for em.Iter.ReadArray() {
		valueRead, isFollowThePath, t, err := em.scanArray(depth, index)
		if err != nil {
			return err
		}

		if isFollowThePath {
			switch t {
			case TypeArray:
				err = em.extractArray(depth + 1)
			case TypeObject:
				err = em.extractObject(depth + 1)
			case TypeArrayIndex:
			default:
				panic("should not happen")
			}

			if err != nil {
				return err
			}
		} else if !valueRead {
			em.Iter.Skip()
		}

		index++
	}

	return nil
}

func (em *ExtractMany) scanArray(depth, index int) (bool, bool, int, error) {
	valueRead := false
	isFollowThePath := false
	typeOfPath := -1

	for _, extract := range em.Extracts {
		if extract.finished() {
			continue
		}

		vRead, followed, t, err := extract.scanArray(em.Iter, depth, index)
		if err != nil {
			return false, false, 0, err
		}

		if vRead {
			valueRead = true
		}

		if followed {
			isFollowThePath = true
			typeOfPath = t
		}
	}

	return valueRead, isFollowThePath, typeOfPath, nil
}
