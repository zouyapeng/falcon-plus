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

package db

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
	"log"
)

func AddGroup(grpName string){
	sql := ""

	sql = fmt.Sprintf("insert into grp(grp_name) values ('%s')", grpName)

	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}
}

func AddTemplate(tplName string, commonTplID int){
	sql := ""

	sql = fmt.Sprintf("insert into tpl(tpl_name, parent_id) values ('%s', %d)", tplName, commonTplID)

	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}
}


func UpdateAgentToGroup(hid int, gid int){
	sql := ""

	sql = fmt.Sprintf("insert into grp_host(grp_id, host_id) values ('%d', '%d')", gid, hid)

	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}
}

func UpdateTemplateToGroup(tid int, gid int){
	sql := ""

	sql = fmt.Sprintf("insert into grp_tpl(tpl_id, grp_id) values ('%d', '%d')", tid, gid)

	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}
}

func UpdateAgent(agentInfo *model.AgentUpdateInfo) {
	sql := ""
	if g.Config().Hosts == "" {
		sql = fmt.Sprintf(
			"insert into host(hostname, instance_id, region, role, product_version, environment, ip, agent_version, plugin_version) values ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') on duplicate key update ip='%s', agent_version='%s', plugin_version='%s'",
			agentInfo.ReportRequest.Hostname,
			agentInfo.ReportRequest.InstanceID,
			agentInfo.ReportRequest.Region,
			agentInfo.ReportRequest.Role,
			agentInfo.ReportRequest.ProductVersion,
			agentInfo.ReportRequest.Environment,
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
		)
	} else {
		// sync, just update
		sql = fmt.Sprintf(
			"update host set ip='%s', instance_id=%s, region=%s, role=%s, product_version=%s, env=%s, agent_version='%s', plugin_version='%s' where hostname='%s'",
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.InstanceID,
			agentInfo.ReportRequest.Region,
			agentInfo.ReportRequest.Role,
			agentInfo.ReportRequest.ProductVersion,
			agentInfo.ReportRequest.Environment,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.Hostname,
		)
	}

	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}

}
