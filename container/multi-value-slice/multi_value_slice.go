package multi_value_slice

import (
	"sync"

	"github.com/pkg/errors"
	consensus_types "github.com/prysmaticlabs/prysm/v4/consensus-types"
	"github.com/prysmaticlabs/prysm/v4/container/multi-value-slice/interfaces"
)

// MultiValueSlice defines an abstraction over all concrete implementations of the generic Slice.
type MultiValueSlice[O interfaces.Identifiable] interface {
	Len(obj O) int
}

// Value defines a single value along with one or more IDs that share this value.
type Value[V any] struct {
	val V
	ids []interfaces.Id
}

// MultiValue defines a collection of Value items.
type MultiValue[V any] struct {
	Values []*Value[V]
}

// Slice is the main component of the multi-value slice data structure. It has two type parameters:
//   - V comparable - the type of values stored the slice. The constraint is required
//     because certain operations (e.g. updating, appending) have to compare values against each other.
//   - O interfaces.Identifiable - the type of objects sharing the slice. The constraint is required
//     because we need a way to compare objects against each other in order to know which objects
//     values should be accessed.
type Slice[V comparable, O interfaces.Identifiable] struct {
	sharedItems     []V
	individualItems map[uint64]*MultiValue[V]
	appendedItems   []*MultiValue[V]
	cachedLengths   map[interfaces.Id]int
	lock            sync.RWMutex
}

// Init initializes the slice with sensible defaults. Input values are assigned to shared items.
func (s *Slice[V, O]) Init(items []V) {
	s.sharedItems = items
	s.individualItems = map[interfaces.Id]*MultiValue[V]{}
	s.appendedItems = []*MultiValue[V]{}
	s.cachedLengths = map[interfaces.Id]int{}
}

// Len returns the number of items for the input object.
func (s *Slice[V, O]) Len(obj O) int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	l, ok := s.cachedLengths[obj.Id()]
	if !ok {
		return len(s.sharedItems)
	}
	return l
}

// Copy copies items between the source and destination.
func (s *Slice[V, O]) Copy(src O, dst O) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range s.individualItems {
	individualLoop:
		for _, v := range item.Values {
			for _, id := range v.ids {
				if id == src.Id() {
					v.ids = append(v.ids, dst.Id())
					break individualLoop
				}
			}
		}
	}

appendedLoop:
	for _, item := range s.appendedItems {
		found := false
	individualLoop2:
		for _, v := range item.Values {
			for _, id := range v.ids {
				if id == src.Id() {
					found = true
					v.ids = append(v.ids, dst.Id())
					break individualLoop2
				}
			}
		}
		if !found {
			// This is an optimization. If we didn't find an appended item at index i,
			// then all larger indices don't have an appended item for the object either.
			break appendedLoop
		}
	}

	srcLen, ok := s.cachedLengths[src.Id()]
	if ok {
		s.cachedLengths[dst.Id()] = srcLen
	}
}

// Value returns all items for the input object.
func (s *Slice[V, O]) Value(obj O) []V {
	s.lock.RLock()
	defer s.lock.RUnlock()

	result := make([]V, len(s.sharedItems))
	for i, item := range s.sharedItems {
		ind, ok := s.individualItems[uint64(i)]
		if !ok {
			result[i] = item
		} else {
			found := false
		individualLoop:
			for _, v := range ind.Values {
				for _, id := range v.ids {
					if id == obj.Id() {
						result[i] = v.val
						found = true
						break individualLoop
					}
				}
			}
			if !found {
				result[i] = item
			}
		}
	}

	for _, item := range s.appendedItems {
		found := false
	individualLoop2:
		for _, v := range item.Values {
			for _, id := range v.ids {
				if id == obj.Id() {
					found = true
					result = append(result, v.val)
					break individualLoop2
				}
			}
		}
		if !found {
			// This is an optimization. If we didn't find an appended item at index i,
			// then all larger indices don't have an appended item for the object either.
			return result
		}
	}

	return result
}

// At returns the item at the requested index for the input object.
// If the object has an individual value at that index, it will be returned. Otherwise the shared value will be returned.
// If the object has an appended value at that index, it will be returned.
func (s *Slice[V, O]) At(obj O, index uint64) (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if index >= uint64(len(s.sharedItems)+len(s.appendedItems)) {
		var def V
		return def, errors.Wrapf(consensus_types.ErrOutOfBounds, "no item at index %d", index)
	}

	isOriginal := index < uint64(len(s.sharedItems))
	if isOriginal {
		ind, ok := s.individualItems[index]
		if !ok {
			return s.sharedItems[index], nil
		}
		for _, v := range ind.Values {
			for _, id := range v.ids {
				if id == obj.Id() {
					return v.val, nil
				}
			}
		}
		return s.sharedItems[index], nil
	} else {
		item := s.appendedItems[index-uint64(len(s.sharedItems))]
		for _, v := range item.Values {
			for _, id := range v.ids {
				if id == obj.Id() {
					return v.val, nil
				}
			}
		}
		var def V
		return def, errors.Wrapf(consensus_types.ErrOutOfBounds, "no item at index %d", index)
	}
}

// UpdateAt updates the item at the required index for the input object to the passed in value.
func (s *Slice[V, O]) UpdateAt(obj O, index uint64, val V) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if index >= uint64(len(s.sharedItems)+len(s.appendedItems)) {
		return errors.Wrapf(consensus_types.ErrOutOfBounds, "no item at index %d", index)
	}

	isOriginal := index < uint64(len(s.sharedItems))
	if isOriginal {
		ind, ok := s.individualItems[index]
		if ok {
		individualLoop:
			for mvi, v := range ind.Values {
				for idi, id := range v.ids {
					if id == obj.Id() {
						if len(v.ids) == 1 {
							// There is an improvement to be made here. If len(ind.Values) == 1,
							// then after removing the item from the slice s.individualItems[i]
							// will be a useless map entry whose value is an empty slice.
							ind.Values = append(ind.Values[:mvi], ind.Values[mvi+1:]...)
						} else {
							v.ids = append(v.ids[:idi], v.ids[idi+1:]...)
						}
						break individualLoop
					}
				}
			}
		}

		if val == s.sharedItems[index] {
			return nil
		}

		if !ok {
			s.individualItems[index] = &MultiValue[V]{Values: []*Value[V]{{val: val, ids: []uint64{obj.Id()}}}}
		} else {
			newValue := true
			for _, v := range ind.Values {
				if v.val == val {
					v.ids = append(v.ids, obj.Id())
					newValue = false
					break
				}
			}
			if newValue {
				ind.Values = append(ind.Values, &Value[V]{val: val, ids: []uint64{obj.Id()}})
			}
		}
	} else {
		item := s.appendedItems[index-uint64(len(s.sharedItems))]
		found := false
	individualLoop2:
		for vi, v := range item.Values {
			for idi, id := range v.ids {
				if id == obj.Id() {
					found = true
					if len(v.ids) == 1 {
						item.Values = append(item.Values[:vi], item.Values[vi+1:]...)
					} else {
						v.ids = append(v.ids[:idi], v.ids[idi+1:]...)
					}
					break individualLoop2
				}
			}
		}
		if !found {
			return errors.Wrapf(consensus_types.ErrOutOfBounds, "no item at index %d", index)
		}

		newValue := true
		for _, v := range item.Values {
			if v.val == val {
				v.ids = append(v.ids, obj.Id())
				newValue = false
				break
			}
		}
		if newValue {
			item.Values = append(item.Values, &Value[V]{val: val, ids: []uint64{obj.Id()}})
		}
	}

	return nil
}

// Append adds a new item to the input object.
func (s *Slice[V, O]) Append(obj O, val V) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.appendedItems) == 0 {
		s.appendedItems = append(s.appendedItems, &MultiValue[V]{Values: []*Value[V]{{val: val, ids: []uint64{obj.Id()}}}})
		s.cachedLengths[obj.Id()] = len(s.sharedItems) + 1
		return
	}

	for _, item := range s.appendedItems {
		found := false
	individualLoop:
		for _, v := range item.Values {
			for _, id := range v.ids {
				if id == obj.Id() {
					found = true
					break individualLoop
				}
			}
		}
		if !found {
			newValue := true
			for _, v := range item.Values {
				if v.val == val {
					v.ids = append(v.ids, obj.Id())
					newValue = false
					break
				}
			}
			if newValue {
				item.Values = append(item.Values, &Value[V]{val: val, ids: []uint64{obj.Id()}})
			}

			l, ok := s.cachedLengths[obj.Id()]
			if ok {
				s.cachedLengths[obj.Id()] = l + 1
			} else {
				s.cachedLengths[obj.Id()] = len(s.sharedItems) + 1
			}

			return
		}
	}

	s.appendedItems = append(s.appendedItems, &MultiValue[V]{Values: []*Value[V]{{val: val, ids: []uint64{obj.Id()}}}})

	s.cachedLengths[obj.Id()] = s.cachedLengths[obj.Id()] + 1
}

// Detach removes the input object from the multi-value slice.
// What this means in practice is that we remove all individual and appended values for that object and clear the cached length.
func (s *Slice[V, O]) Detach(obj O) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i, ind := range s.individualItems {
	individualLoop:
		for vi, v := range ind.Values {
			for idi, id := range v.ids {
				if id == obj.Id() {
					if len(v.ids) == 1 {
						if len(ind.Values) == 1 {
							delete(s.individualItems, i)
						} else {
							ind.Values = append(ind.Values[:vi], ind.Values[vi+1:]...)
						}
					} else {
						v.ids = append(v.ids[:idi], v.ids[idi+1:]...)
					}
					break individualLoop
				}
			}
		}
	}

appendedLoop:
	for _, item := range s.appendedItems {
		found := false
	individualLoop2:
		for vi, v := range item.Values {
			for idi, id := range v.ids {
				if id == obj.Id() {
					found = true
					if len(v.ids) == 1 {
						item.Values = append(item.Values[:vi], item.Values[vi+1:]...)
					} else {
						v.ids = append(v.ids[:idi], v.ids[idi+1:]...)
					}
					break individualLoop2
				}
			}
		}
		if !found {
			// This is an optimization. If we didn't find an appended item at index i,
			// then all larger indices don't have an appended item for the object either.
			break appendedLoop
		}
	}

	delete(s.cachedLengths, obj.Id())
}