package jsoniter

import (
	"reflect"
	"sync/atomic"
	"unsafe"
)

type Config struct {
	decoderCache unsafe.Pointer
	encoderCache unsafe.Pointer
}

var DEFAULT_CONFIG = &Config{}

func init() {
	initConfig(DEFAULT_CONFIG)
}

func initConfig(cfg *Config) {
	atomic.StorePointer(&cfg.decoderCache, unsafe.Pointer(&map[string]Decoder{}))
	atomic.StorePointer(&cfg.encoderCache, unsafe.Pointer(&map[string]Encoder{}))
}
func (cfg *Config) addDecoderToCache(cacheKey reflect.Type, decoder Decoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.decoderCache)
		cache := *(*map[reflect.Type]Decoder)(ptr)
		copied := map[reflect.Type]Decoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = decoder
		done = atomic.CompareAndSwapPointer(&cfg.decoderCache, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *Config) addEncoderToCache(cacheKey reflect.Type, encoder Encoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.encoderCache)
		cache := *(*map[reflect.Type]Encoder)(ptr)
		copied := map[reflect.Type]Encoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = encoder
		done = atomic.CompareAndSwapPointer(&cfg.encoderCache, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *Config) getDecoderFromCache(cacheKey reflect.Type) Decoder {
	ptr := atomic.LoadPointer(&cfg.decoderCache)
	cache := *(*map[reflect.Type]Decoder)(ptr)
	return cache[cacheKey]
}

func (cfg *Config) getEncoderFromCache(cacheKey reflect.Type) Encoder {
	ptr := atomic.LoadPointer(&cfg.encoderCache)
	cache := *(*map[reflect.Type]Encoder)(ptr)
	return cache[cacheKey]
}

// CleanDecoders cleans decoders registered or cached
func (cfg *Config) CleanDecoders() {
	typeDecoders = map[string]Decoder{}
	fieldDecoders = map[string]Decoder{}
	atomic.StorePointer(&cfg.decoderCache, unsafe.Pointer(&map[string]Decoder{}))
}

// CleanEncoders cleans encoders registered or cached
func (cfg *Config) CleanEncoders() {
	typeEncoders = map[string]Encoder{}
	fieldEncoders = map[string]Encoder{}
	atomic.StorePointer(&cfg.encoderCache, unsafe.Pointer(&map[string]Encoder{}))
}
