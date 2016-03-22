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

// A Disposable represents a type that needs to release resources before
// dereferencing.
type Disposable interface {
	// Dispose releases resources references.
	Dispose()
}

// Using function allows to define a execution block for a disposable object.
//
// 	type Foo struct {
// 		Number int
// 	}
//
// 	func (f *Foo) Dispose() {
// 		f.Number = 0
// 	}
//
// 	func main() {
// 		val1 := &Foo{50}
// 		defer val1.Dispose()
//
// 		Using(&Foo{40}, func(val2 Disposable) {
// 			fmt.Println("The number is", val2)
// 		})
// 	}
func Using(d Disposable, f func(Disposable)) {
	defer d.Dispose()
	f(d)
}
