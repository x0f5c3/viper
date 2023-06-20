package encoding

import (

	"sync"

)

type GenericKV[V any] struct {
	key string
	val V
}

func (g *GenericKV[V]) Key() string {
	return g.key
}

func (g *GenericKV[V]) SetKey(key string) {
	g.key = key
}

func (g *GenericKV[V]) Val() V {
	return g.val
}

func (g *GenericKV[V]) SetVal(val V) {
	g.val = val
}

type KVMerge struct {
	m map[string]GenericKV[any]
}

type GenericEncoder interface {
	Encode(v map[string]GenericKV[any]) ([]byte, error)
}

// Encoder encodes the contents of v into a byte representation.
// It's primarily used for encoding a map[string]interface{} into a file format.
type Encoder interface {
	Encode(v map[string]interface{}) ([]byte, error)
}

// NewGenericEncoderRegistry returns a new, initialized EncoderRegistry.
func NewGenericEncoderRegistry() *EncoderRegistry {
	return &EncoderRegistry{
		encoders: make(map[string]Encoder),
	}
}

// RegisterEncoder registers an Encoder for a format.
// Registering a Encoder for an already existing format is not supported.
func (e *EncoderRegistry) RegisterEncoder(format string, enc Encoder) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.encoders[format]; ok {
		return ErrEncoderFormatAlreadyRegistered
	}

	e.encoders[format] = enc

	return nil
}

func (e *EncoderRegistry) Encode(format string, v map[string]interface{}) ([]byte, error) {
	e.mu.RLock()
	encoder, ok := e.encoders[format]
	e.mu.RUnlock()

	if !ok {
		return nil, ErrEncoderNotFound
	}

	return encoder.Encode(v)
}

const (
	// ErrEncoderNotFound is returned when there is no encoder registered for a format.
	ErrEncoderNotFound = encodingError("encoder not found for this format")

	// ErrEncoderFormatAlreadyRegistered is returned when an encoder is already registered for a format.
	ErrEncoderFormatAlreadyRegistered = encodingError("encoder already registered for this format")
)

type GenericEncoderRegistry struct {
	encoders map[string]GenericEncoder
	mu       sync.RWMutex
}

// EncoderRegistry can choose an appropriate Encoder based on the provided format.
type EncoderRegistry struct {
	encoders map[string]Encoder

	mu sync.RWMutex
}

// NewEncoderRegistry returns a new, initialized EncoderRegistry.
func NewEncoderRegistry() *EncoderRegistry {
	return &EncoderRegistry{
		encoders: make(map[string]Encoder),
	}
}

// RegisterEncoder registers an Encoder for a format.
// Registering a Encoder for an already existing format is not supported.
func (e *EncoderRegistry) RegisterEncoder(format string, enc Encoder) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.encoders[format]; ok {
		return ErrEncoderFormatAlreadyRegistered
	}

	e.encoders[format] = enc

	return nil
}

func (e *EncoderRegistry) Encode(format string, v map[string]interface{}) ([]byte, error) {
	e.mu.RLock()
	encoder, ok := e.encoders[format]
	e.mu.RUnlock()

	if !ok {
		return nil, ErrEncoderNotFound
	}

	return encoder.Encode(v)
}
