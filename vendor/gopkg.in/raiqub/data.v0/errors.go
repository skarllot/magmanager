/*
 * Copyright 2016 Fabrício Godoy
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

package data

import (
	"fmt"
)

// A InvalidTypeError represents an error when value type is different than
// expected.
type InvalidTypeError struct {
	Value interface{}
}

// NewInvalidTypeError returns a new instance of InvalidTypeError.
func NewInvalidTypeError(value interface{}) InvalidTypeError {
	return InvalidTypeError{value}
}

// Error returns string representation of current instance error.
func (e InvalidTypeError) Error() string {
	return fmt.Sprintf("Unexpected type: %T", e.Value)
}
