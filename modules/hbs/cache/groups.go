// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
	"sync"
)

// 一个机器可能在多个group下，做一个map缓存hostid与groupid的对应关系
type SafeHostGroupsMap struct {
	sync.RWMutex
	M map[int][]int
}

type SafeGroupsMap struct {
	sync.RWMutex
	M map[string]int
}

var HostGroupsMap = &SafeHostGroupsMap{M: make(map[int][]int)}
var GroupsMap = &SafeGroupsMap{M: make(map[string]int)}

func (this *SafeHostGroupsMap) GetGroupIds(hid int) ([]int, bool) {
	this.RLock()
	defer this.RUnlock()
	gids, exists := this.M[hid]
	return gids, exists
}


func (this *SafeHostGroupsMap) Init() {
	m, err := db.QueryHostGroups()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeGroupsMap) GetID(grpName string) (int, bool) {
	this.RLock()
	defer this.RUnlock()
	id, exists := this.M[grpName]
	return id, exists
}

func (this *SafeGroupsMap) Init() {
	m, err := db.QueryGroups()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}
