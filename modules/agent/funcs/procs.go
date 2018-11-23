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

package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/nux"
	"log"
	"strings"
	"io/ioutil"
	"github.com/toolkits/file"
	"fmt"
	"github.com/struCoder/pidusage"
)

type Resource struct {
	Fd int
	MemUsage float64
	CpuUsage float64
}

func ProcMetrics() (L []*model.MetricValue) {
	reportProcs := g.ReportProcs()
	sz := len(reportProcs)
	if sz == 0 {
		return
	}

	ps, err := nux.AllProcs()
	if err != nil {
		log.Println(err)
		return
	}

	pslen := len(ps)

	for tags, m := range reportProcs {
		cnt := 0
		for i := 0; i < pslen; i++ {
			if is_a(ps[i], m) {
				cnt++
			}
		}

		L = append(L, GaugeValue(g.PROC_NUM, cnt, tags))
	}

	return
}


func ProcResourceMetrics() (L []*model.MetricValue) {
	reportProcs := g.ReportProcsResource()
	sz := len(reportProcs)
	if sz == 0 {
		return
	}

	ps, err := nux.AllProcs()
	if err != nil {
		log.Println(err)
		return
	}

	pslen := len(ps)
	var findPss []*nux.Proc
	memInfo, err := nux.MemInfo()
	cpuNum := nux.NumCpu()

	for tags, m := range reportProcs {
		for i := 0; i < pslen; i++ {
			if is_a(ps[i], m) {
				findPss = append(findPss, ps[i])
			}
		}

		var resourceAllOfOne *Resource
		for index, ps := range findPss {
			resource, err := collectProcessResources(ps.Pid)
			if err != nil{
				continue
			}

			if index == 0{
				resourceAllOfOne = resource
			} else {
				resourceAllOfOne.Fd += resource.Fd
				resourceAllOfOne.CpuUsage += resource.CpuUsage
				resourceAllOfOne.MemUsage += resource.MemUsage
			}
		}

		L = append(L, GaugeValue(g.PROC_CPU, resourceAllOfOne.CpuUsage / float64(cpuNum), tags))
		L = append(L, GaugeValue(g.PROC_MEM, resourceAllOfOne.MemUsage/ float64(memInfo.MemTotal), tags))
		L = append(L, GaugeValue(g.PROC_FD, resourceAllOfOne.Fd, tags))
	}

	return
}


func collectProcessResources(pid int) (resource *Resource,err error){
	FdFile := fmt.Sprintf("/proc/%d/fd", pid)
	if !file.IsExist(FdFile) {
		return
	}
	files, err := ioutil.ReadDir(FdFile)

	if err != nil {
		log.Fatal(err)
	}

	resource.Fd = len(files) - 3
	if resource.Fd < 0{
		resource.Fd = 0
	}

	sysInfo, err := pidusage.GetStat(pid)
	resource.MemUsage = sysInfo.Memory
	resource.CpuUsage = sysInfo.CPU

	return
}


func is_a(p *nux.Proc, m map[int]string) bool {
	// only one kv pair
	for key, val := range m {
		if key == 1 {
			// name
			if val != p.Name {
				return false
			}
		} else if key == 2 {
			// cmdline
			if !strings.Contains(p.Cmdline, val) {
				return false
			}
		}
	}
	return true
}
