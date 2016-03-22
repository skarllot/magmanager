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

/*
Package dot provides implementation of basic common tasks.

Disposable

A Disposable defines a predictable way of releasing resources that cannot be
handled by garbage collector as open files and streams.

Call Dispose function when the object is no longer needed and its allocated
resources can be released.

Can be used in conjunction of Using function to define a execution block and
automatically release allocated resources at the end of the block.

LockStatus

A LockStatus defines values for expected lock status for a lockable object.

MulticastDispose

A MulticastDispose allows to register multiple Dispose functions on a single
Disposable object (itself).

Very useful as returning object of a function that allocates several Disposable
objects.

StringSlice

A StringSlice is a named type of slice of strings and provides some useful
methods as:

	- IndexOf to determine index of specified string;
	- Exists to determine whether specified string exists on slice;
	- ExistsAll to determine whether all specified strings exists on slice;
	- TrueForAll to test whether every element of slice matches the conditions
	specified by the predicate function.

Using

A Using function allows to define a execution block for a Disposable object to
automatically release allocated resources at the end of the function/block.

WaitFunc

A WaitFunc function allows to synchronously waits for a function until it
returns true or until a specified timeout is reached.

WaitPeerListening

A WaitPeerListening function allows to synchronously waits for a peer until it
is ready for new connections or until a specified timeout is reached.
*/
package dot
