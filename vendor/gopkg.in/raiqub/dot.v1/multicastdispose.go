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

package dot

// A MulticastDispose allows to register multiple Dispose functions in one
// object.
type MulticastDispose struct {
	list []func()
}

// Add adds one or more Dispose function.
func (md *MulticastDispose) Add(f ...func()) {
	md.list = append(md.list, f...)
}

// AddDisposable adds one or more Disposable instances.
func (md *MulticastDispose) AddDisposable(d ...Disposable) {
	for _, v := range d {
		md.list = append(md.list, v.Dispose)
	}
}

// Dispose disposes all Dispose functions registered by this instance.
func (md *MulticastDispose) Dispose() {
	for i := len(md.list) - 1; i >= 0; i-- {
		md.list[i]()
	}

	md.list = nil
}

// NewMulticastDispose creates a new instance of MulticastDispose.
func NewMulticastDispose() *MulticastDispose {
	return &MulticastDispose{
		make([]func(), 0),
	}
}
