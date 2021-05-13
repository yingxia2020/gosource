/* Copyright (C) Intel Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia <ying.xia@intel.com>, 2019
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//const HEADER = `#!/bin/bash`

const HEADER = `#===============================================================================
# Copyright 2020 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#===============================================================================

#!/bin/bash`

const NEWHEADER = `#!/bin/bash
#===============================================================================
# Copyright 2020 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#===============================================================================
`

func main() {
	var files []string

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sh") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		myText := strings.Replace(string(content), HEADER, NEWHEADER, -1) 
		err = ioutil.WriteFile(file, []byte(myText), 0664)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Finished add header for:", file)
	}
}
