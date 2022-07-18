package placeholders

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
)

// placeholderManagerInit is used for the lazy initialization of
// placeholderManager.
var placeholderManagerInit sync.Once

// placeholderManager is a singleton used to keep track of created placeholders.
// It is initialized lazily only when it is actually used.
var placeholderManager struct {
	sequenceNumber uint64
	placeholdersMx sync.RWMutex
	placeholders   map[unsafe.Pointer]unsafe.Pointer
}

// allocationWrapper is used to ensure that the allocated object is placed at a
// unique memory location. See the comments inside the implementation of Make
// for more details.
type allocationWrapper[T any] struct {
	obj T
	sn  uint64
}

// TestingCleanup is an interface implemented by *testing.T and *testing.B.
type TestingCleanup interface {
	// Cleanup is used to reclaim the resources allocated to keep track of the
	// placeholders.
	Cleanup(f func())
}

// Make allocates a new zero-valued object of type T and returns an object of
// type TPtr that references the allocated object. TPtr is either *T or has *T
// as its underlying type (i.e, defined as "type Foo *Bar", where TPtr = Foo
// and T = Bar). The returned reference or any other reference to the allocated
// object is considered a placeholder reference: when cmp.Equal or cmp.Diff is
// invoked with the placeholders.Ignore option, a placeholder reference is
// considered to be equal to any other object, regardless of types or values.
//
// Only pointers to the allocated object are considered placeholders and not
// the object itself (see examples below).
//
// The second type parameter (T) can always be inferred from the first one
// (TPtr) and does not need to be explicitly specified.
//
// The only parameter (t) can be of type *testing.T, *testing.B, or any other
// type with a function Cleanup with a similar semantic. It ensures that the
// resources allocated to keep track of the placeholders are eventually
// reclaimed (e.g., when t is of type *testing.T, they are reclaimed upon the
// completion of the test).
//
// Example:
//    helloString := "hello"
//
//    // true
//    cmp.Equal(placeholders.Make[*string](t), &helloString, placeholders.Ignore())
//
//    placeholder := placeholders.Make[*string](t)
//    anotherRef := &(*placeholder)
//
//    // true, any reference to the allocated object is a placeholder
//    cmp.Equal(anotherRef, &helloString, placeholders.Ignore())
//
//    // false, the allocated object itself is not a placeholder
//    cmp.Equal(*placeholder, "hello", placeholders.Ignore())
//
//    type Foo struct {SPtr *string; S string}
//
//    // true, it works with struct fields and embedded types as well!
//    cmp.Equal(Foo{placeholders.Make[*string](t), "world"}, Foo{&helloString, "world"}, placeholders.Ignore())
//
//    // false, non-placeholder fields differ
//    cmp.Equal(Foo{placeholders.Make[*string](t), "earthlings"}, Foo{&helloString, "world"}, placeholders.Ignore())
func Make[TPtr ~*T, T any](t TestingCleanup) TPtr {
	placeholderManagerInit.Do(func() {
		// placeholderManager.placeholdersMx doesn't need initialization.
		placeholderManager.placeholders = make(map[unsafe.Pointer]unsafe.Pointer)
	})

	// The unique sequence number sn in the wrapper ensures that the allocated
	// object lives in its own memory location and the pointer to it is unique.
	//
	// This is done to overcome a weird property of new(T) that the created
	// object may share its location in memory with another zero-valued object
	// of the same type.
	//
	// In the current version of Go compiler, this wrapper is only necessary
	// for T = struct{}, but there is no guarantee that it won't change in the
	// future. Hence, it is safer to always use the wrapper.
	wrapper := &allocationWrapper[T]{
		sn: atomic.AddUint64(&placeholderManager.sequenceNumber, 1),
	}

	newTPtr := &wrapper.obj
	// Note that, untypedPtr is a strong pointer to the object. Notably:
	//   (1) it prevents the underlying object from being garbage-collected;
	//   (2) if the underlying object is moved, untypedPtr will be updated.
	untypedPtr := unsafe.Pointer(newTPtr)

	// The memory allocated for this placeholder will be reclaimed when the test completes.
	t.Cleanup(func() {
		placeholderManager.placeholdersMx.Lock()
		defer placeholderManager.placeholdersMx.Unlock()
		delete(placeholderManager.placeholders, untypedPtr)
	})

	placeholderManager.placeholdersMx.Lock()
	defer placeholderManager.placeholdersMx.Unlock()
	// A pointer to the wrapper is stored in the map to make sure the wrapper
	// is not garbage-collected.
	placeholderManager.placeholders[untypedPtr] = unsafe.Pointer(wrapper)

	return newTPtr
}

// IsPlaceholder returns true iff obj is a reference to an object allocated
// with placeholders.Make.
func IsPlaceholder(obj any) bool {
	if obj == nil {
		return false
	}
	if reflect.ValueOf(obj).Kind() != reflect.Pointer {
		return false
	}

	placeholderManager.placeholdersMx.RLock()
	defer placeholderManager.placeholdersMx.RUnlock()

	_, found := placeholderManager.placeholders[reflect.ValueOf(obj).UnsafePointer()]
	return found
}

// Ignore returns a cmp.Ignore option that ignores all placeholders.
func Ignore() cmp.Option {
	filter := func(a, b any) bool {
		return IsPlaceholder(a) || IsPlaceholder(b)
	}

	return cmp.FilterValues(filter, cmp.Ignore())
}
