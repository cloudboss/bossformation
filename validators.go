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
	"fmt"
	r "regexp"

	v "github.com/asaskevich/govalidator"
)

func init() {
	v.TagMap["region"] = v.Validator(IsRegion)
	v.TagMap["subnet"] = v.Validator(IsSubnet)
	v.TagMap["effect"] = v.Validator(IsEffect)
	v.TagMap["vpc"] = v.Validator(IsVpc)
	v.TagMap["scheme"] = v.Validator(IsScheme)
}

func IsRegion(s string) (b bool) {
	regions := map[string]bool{
		"ap-northeast-1": true,
		"ap-northeast-2": true,
		"ap-southeast-1": true,
		"ap-southeast-2": true,
		"eu-central-1":   true,
		"eu-west-1":      true,
		"sa-east-1":      true,
		"us-east-1":      true,
		"us-west-1":      true,
		"us-west-2":      true,
	}
	if !regions[s] {
		return
	}
	return true
}

func isAwsId(s, rsrc string) bool {
	pat := fmt.Sprintf("^%s-[0-9a-fA-F]+$", rsrc)
	b, _ := r.MatchString(pat, s)
	return b
}

func IsSubnet(s string) (b bool) {
	return isAwsId(s, "subnet")
}

func IsVpc(s string) (b bool) {
	return isAwsId(s, "vpc")
}

func IsScheme(s string) (b bool) {
	return s == "internal" || s == "internet-facing"
}

func IsEffect(s string) (b bool) {
	return
}
