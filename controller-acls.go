// Copyright 2019 HAProxy Technologies LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/haproxytech/models"
)

func (c *HAProxyController) addACL(acl models.ACL, frontends ...string) {
	if len(frontends) == 0 {
		frontends = []string{FrontendHTTP, FrontendHTTPS}
	}
	for _, frontend := range frontends {
		_, acls, err := c.NativeAPI.Configuration.GetACLs("frontend", frontend, c.ActiveTransaction)
		found := false
		if err == nil {
			for _, d := range acls {
				if acl.ACLName == d.ACLName {
					found = true
					break
				}
			}
		}
		if !found {
			err = c.NativeAPI.Configuration.CreateACL("frontend", frontend, &acl, c.ActiveTransaction, 0)
			LogErr(err)
		}
	}
}

func (c *HAProxyController) removeACL(acl models.ACL, frontends ...string) {
	nativeAPI := c.NativeAPI
	for _, frontend := range frontends {
		_, acls, err := nativeAPI.Configuration.GetACLs("frontend", frontend, c.ActiveTransaction)
		if err == nil {
			indexShift := int64(0)
			for _, d := range acls {
				if acl.ACLName == d.ACLName {
					err = nativeAPI.Configuration.DeleteACL(*d.ID-indexShift, "frontend", frontend, c.ActiveTransaction, 0)
					LogErr(err)
					if err == nil {
						indexShift++
					}
				}
			}
		}
	}
}
