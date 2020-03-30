/*
 * Copyright (c) 2020 wellwell.work, LLC by Zoe
 *
 * Licensed under the Apache License 2.0 (the "License");
 * You may not use this file except in compliance with the License.
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

package unipaw

import (
	"strings"
)

// PathJoin ...
func PathJoin(a ...string) string {
	switch len(a) {
	case 0:
		return "/"
	case 1:
		if strings.HasPrefix(a[0], "/") {
			return "/" + a[0]
		}
		return a[0]
	case 2:
		var x, y = a[0], a[1]
		if !strings.HasSuffix(x, "/") && !strings.HasPrefix(y, "/") {
			return x + "/" + y
		}
		return x + y
	default:
		var x = PathJoin(a[1:]...)
		return PathJoin(a[0], x)
	}
}
