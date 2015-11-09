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

// A PredicateStringFunc represents a function that defines a criteria and
// determines whether specified string meets that criteria.
type PredicateStringFunc func(string) bool

// TrueForAll always returns true, regardless specified string.
func TrueForAll(string) bool {
	return true
}
