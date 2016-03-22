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
Package crypt provides some cryptographic operations.

RandomAggr

A RandomAggr allows to aggregate random data sources to fill a buffer. Each
source has weight to control the percentage from total to be read.

It implements the io.ReadCloser interface to allow to close sources if needed.

Salter

A Salter provides a random data generator to password salt and unique session
IDs. Every token generated is used to salt next token to increase
unpredictability of generated data.

SSTDEG

A SSTDEG provides a pseudo-random generator based on syscall time deltas of
Sleep calls. It implements io.Reader interface.
*/
package crypt
