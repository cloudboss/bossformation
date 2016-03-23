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

	v "github.com/asaskevich/govalidator"
	"github.com/aws/aws-sdk-go/service/ec2"
	cf "github.com/crewjam/go-cloudformation"
)

var ec2Client *ec2.EC2

type Cluster struct {
	Name             string `valid:"alphanum" json:"name"`
	Region           string `valid:"region" json:"region"`
	VpcId            string `valid:"required,vpc" json:"vpcId"`
	AutoscalingGroup `valid:"required" json:"autoscalingGroup"`
	LoadBalancer     `valid:"optional" json:"loadBalancer"`
}

type AutoscalingGroup struct {
	Image              string `valid:"required" json:"image"`
	InstanceType       string `valid:"required" json:"instanceType"`
	Subnets            `valid:"required" json:"subnets"`
	IamInstanceProfile `valid:"optional" json:"iamInstanceprofile"`
}

type IamInstanceProfile struct {
	Roles  []string `valid:"optional" json:"roles"`
	Policy `valid:"optional" json:"roles"`
}

type Policy struct {
	Effect string   `valid:"required,effect" json:"effect"`
	Action []string `valid:"required,action" json:"action"`
}

type LoadBalancer struct {
	Scheme      string `valid:"scheme" json:"scheme"`
	Public      bool   `valid:"optional" json:"public"`
	HealthCheck `valid:"required" json:"healthCheck"`
	Subnets     `valid:"required" json:"subnets"`
}

type Subnets struct {
	Tag     string   `valid:"optional,ascii" json:"tag"`
	TagName string   `valid:"optional,ascii" json:"tagName"`
	Ids     []string `valid:"optional" json:"ids,omitempty"`
}

type HealthCheck struct {
	Target             string `valid:"alphanum" json:"target"`
	HealthyThreshold   string `valid:"numeric" json:"healthyThreshold"`
	UnhealthyThreshold string `valid:"numeric" json:"unhealthyThreshold"`
	Interval           string `valid:"numeric" json:"interval"`
	Timeout            string `valid:"numeric" json:"timeout"`
}

func NewCluster() *Cluster {
	return &Cluster{
		AutoscalingGroup: AutoscalingGroup{
			Subnets: Subnets{
				TagName: "Name",
			},
		},
	}
}

func (c *Cluster) Validate() (bool, error) {
	if b, err := v.ValidateStruct(c); err != nil {
		return b, err
	}
	if c.LoadBalancer.Subnets.Tag == "" && len(c.LoadBalancer.Subnets.Ids) == 0 {
		return false, fmt.Errorf("Subnet tag or ids required")
	}
	return true, nil
}

func (c *Cluster) BeforeRender() (Context, error) {
	return Context{}, nil
}

func (c *Cluster) Render(ctx Context) (string, error) {
	t := cf.NewTemplate()
	t.AddResource("VPC", cf.EC2VPC{CidrBlock: cf.String("10.0.0.0/16")})

	if buf, err := json.MarshalIndent(t, "", "  "); err != nil {
		return "", err
	} else {
		return string(buf), nil
	}
}

func (c *Cluster) subnets() ([]string, error) {
	if len(c.AutoscalingGroup.Subnets.Ids) > 0 {
		return c.AutoscalingGroup.Subnets.Ids, nil
	} else {
		subnets, err := SubnetIdsByTag(
			c.ec2Client(),
			c.AutoscalingGroup.Subnets.Tag,
			c.VpcId,
		)
		if err != nil {
			return nil, err
		} else {
			return subnets, nil
		}
	}
}

func (c *Cluster) ec2Client() *ec2.EC2 {
	if ec2Client == nil {
		ec2Client = Client(c.Region)
	}
	return ec2Client
}
