/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// MachineTypeRootVolume represents the values of the 'machine_type_root_volume' type.
//
// Machine type root volume.
type MachineTypeRootVolume struct {
	bitmap_ uint32
	aws     *AWSVolume
	gcp     *GCPVolume
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *MachineTypeRootVolume) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AWS returns the value of the 'AWS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS volume specification
func (o *MachineTypeRootVolume) AWS() *AWSVolume {
	if o != nil && o.bitmap_&1 != 0 {
		return o.aws
	}
	return nil
}

// GetAWS returns the value of the 'AWS' attribute and
// a flag indicating if the attribute has a value.
//
// AWS volume specification
func (o *MachineTypeRootVolume) GetAWS() (value *AWSVolume, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.aws
	}
	return
}

// GCP returns the value of the 'GCP' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP Volume specification
func (o *MachineTypeRootVolume) GCP() *GCPVolume {
	if o != nil && o.bitmap_&2 != 0 {
		return o.gcp
	}
	return nil
}

// GetGCP returns the value of the 'GCP' attribute and
// a flag indicating if the attribute has a value.
//
// GCP Volume specification
func (o *MachineTypeRootVolume) GetGCP() (value *GCPVolume, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.gcp
	}
	return
}

// MachineTypeRootVolumeListKind is the name of the type used to represent list of objects of
// type 'machine_type_root_volume'.
const MachineTypeRootVolumeListKind = "MachineTypeRootVolumeList"

// MachineTypeRootVolumeListLinkKind is the name of the type used to represent links to list
// of objects of type 'machine_type_root_volume'.
const MachineTypeRootVolumeListLinkKind = "MachineTypeRootVolumeListLink"

// MachineTypeRootVolumeNilKind is the name of the type used to nil lists of objects of
// type 'machine_type_root_volume'.
const MachineTypeRootVolumeListNilKind = "MachineTypeRootVolumeListNil"

// MachineTypeRootVolumeList is a list of values of the 'machine_type_root_volume' type.
type MachineTypeRootVolumeList struct {
	href  string
	link  bool
	items []*MachineTypeRootVolume
}

// Len returns the length of the list.
func (l *MachineTypeRootVolumeList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *MachineTypeRootVolumeList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *MachineTypeRootVolumeList) Get(i int) *MachineTypeRootVolume {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *MachineTypeRootVolumeList) Slice() []*MachineTypeRootVolume {
	var slice []*MachineTypeRootVolume
	if l == nil {
		slice = make([]*MachineTypeRootVolume, 0)
	} else {
		slice = make([]*MachineTypeRootVolume, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *MachineTypeRootVolumeList) Each(f func(item *MachineTypeRootVolume) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *MachineTypeRootVolumeList) Range(f func(index int, item *MachineTypeRootVolume) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
