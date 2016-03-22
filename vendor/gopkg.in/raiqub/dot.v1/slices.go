/*
 * Copyright 2015 Fabrício Godoy
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dot

import "strings"

// A StringSlice represents an array of string.
type StringSlice []string

// IndexOf looks for specified string into current slice, and optionally ignores
// letter casing.
func (s StringSlice) IndexOf(str string, ignoreCase bool) int {
	if ignoreCase {
		str = strings.ToLower(str)
		for i, v := range s {
			if str == strings.ToLower(v) {
				return i
			}
		}

		return -1
	}

	for i, v := range s {
		if str == v {
			return i
		}
	}

	return -1
}

// Exists determines whether specified string exists into current slice.
func (s StringSlice) Exists(str string, ignoreCase bool) bool {
	return s.IndexOf(str, ignoreCase) != -1
}

// ExistsAll determine whether all specified strings exists into
// current slice.
func (s StringSlice) ExistsAll(str []string, ignoreCase bool) bool {
	for _, v := range str {
		if s.IndexOf(v, ignoreCase) == -1 {
			return false
		}
	}

	return true
}

// ExistsAny determine whether any of specified strings exists into current
// slice
func (s StringSlice) ExistsAny(str []string, ignoreCase bool) bool {
	for _, v := range str {
		if s.IndexOf(v, ignoreCase) != -1 {
			return true
		}
	}

	return false
}

// TrueForAll tests whether every element of current slice matches the
// conditions specified by predicate.
func (s StringSlice) TrueForAll(pred PredicateStringFunc) bool {
	for _, v := range s {
		if !pred(v) {
			return false
		}
	}

	return true
}

// TrueForAny tests whether any element of current slice matches the conditions
// specified by predicate.
func (s StringSlice) TrueForAny(pred PredicateStringFunc) bool {
	for _, v := range s {
		if pred(v) {
			return true
		}
	}

	return false
}
