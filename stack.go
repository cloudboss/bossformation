// Copyright © 2016 Joseph Wright <rjosephwright@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package bf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type stackMap map[string]func() Stack

var types stackMap = stackMap{
	"Cluster": func() Stack { return new(Cluster) },
}

type Context map[string]interface{}

type Stack interface {
	Validate() (bool, error)
	BeforeRender() (Context, error)
	Render(Context) (string, error)
}

func validate() error {
	return nil
}

func LoadStackFromFile(path string) (Stack, error) {
	if bytes, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else {
		if stack, err := LoadStack(bytes); err != nil {
			return nil, err
		} else {
			if _, err := stack.Validate(); err != nil {
				return nil, err
			} else {
				return stack, nil
			}
		}
	}
}

func LoadStack(bytes []byte) (Stack, error) {
	var m map[string]interface{}

	// First unmarshal into map to get "type" as string
	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	stype := m["kind"].(string)
	if stype == "" {
		return nil, fmt.Errorf("Stack kind is required")
	}

	stack := types[stype]()

	// Now the real unmarshal into struct of actual type
	if err := json.Unmarshal(bytes, stack); err != nil {
		return nil, err
	}
	return stack, nil
}
