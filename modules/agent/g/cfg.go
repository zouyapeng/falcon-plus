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

package g

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/toolkits/file"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type PluginConfig struct {
	Enabled bool   `json:"enabled"`
	Dir     string `json:"dir"`
	Git     string `json:"git"`
	LogDir  string `json:"logs"`
}

type HeartbeatConfig struct {
	Enabled  bool   `json:"enabled"`
	Addr     string `json:"addr"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type TransferConfig struct {
	Enabled  bool     `json:"enabled"`
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

type HttpConfig struct {
	Enabled  bool   `json:"enabled"`
	Listen   string `json:"listen"`
	Backdoor bool   `json:"backdoor"`
}

type CollectorConfig struct {
	IfacePrefix []string `json:"ifacePrefix"`
	MountPoint  []string `json:"mountPoint"`
}

type GlobalConfig struct {
	Debug          bool              `json:"debug"`
	Hostname       string            `json:"hostname"`
	IP             string            `json:"ip"`
	InstanceID     string            `json:"instance_id"`
	Region         string            `json:"region"`
	Role           string            `json:"role"`
	ProductVersion string            `json:"product_version"`
	Environment    string            `json:"environment"`
	Plugin         *PluginConfig     `json:"plugin"`
	Heartbeat      *HeartbeatConfig  `json:"heartbeat"`
	Transfer       *TransferConfig   `json:"transfer"`
	Http           *HttpConfig       `json:"http"`
	Collector      *CollectorConfig  `json:"collector"`
	DefaultTags    map[string]string `json:"default_tags"`
	IgnoreMetrics  map[string]bool   `json:"ignore"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func Hostname() (string, error) {
	hostname := Config().Hostname
	if hostname != "" {
		return hostname, nil
	}

	if os.Getenv("FALCON_ENDPOINT") != "" {
		hostname = os.Getenv("FALCON_ENDPOINT")
		return hostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	hostname = strings.Replace(hostname, ".", "-", -1)
	return hostname, err
}

func InstanceID() string {
	instanceID := Config().InstanceID
	if instanceID != "" {
		return instanceID
	}

	httpClient := &http.Client{}
	httpClient.Timeout = 3 * time.Second

	resp, err := httpClient.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		log.Println("ERROR: Get instance_id from AWS fail", err)
		return ""
	}
	defer resp.Body.Close()

	Response, _ := ioutil.ReadAll(resp.Body)
	instanceID = string(Response)
	return instanceID
}

func Region() string {
	region := Config().Region
	if region != "" {
		return region
	}

	httpClient := &http.Client{}
	httpClient.Timeout = 3 * time.Second

	resp, err := httpClient.Get("http://169.254.169.254/latest/meta-data/placement/availability-zone")
	if err != nil {
		log.Println("ERROR: Get region from AWS fail", err)
		return ""
	}

	defer resp.Body.Close()

	Response, _ := ioutil.ReadAll(resp.Body)
	region = string(Response)
	region = region[:len(region)-1]

	return region
}

func Role() string {
	role := Config().Role
	if role != "" {
		return role
	}

	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	hostnameSplit := strings.Split(hostname, "-")

	if len(hostnameSplit) == 7 &&
		(hostnameSplit[2] == "20" || hostnameSplit[2] == "30") &&
		(hostnameSplit[5] == "ops" || hostnameSplit[5] == "pro") {
		role = hostnameSplit[0]
		return role
	}

	return "default"
}

func ProductVersion() string {
	productVersion := Config().ProductVersion
	if productVersion != "" {
		return productVersion
	}

	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	hostnameSplit := strings.Split(hostname, "-")

	if len(hostnameSplit) == 7 &&
		(hostnameSplit[2] == "20" || hostnameSplit[2] == "30") &&
		(hostnameSplit[5] == "ops" || hostnameSplit[5] == "pro") {
		productVersion = hostnameSplit[2]
		return productVersion
	}

	return ""
}

func Environment() string {
	environment := Config().Environment
	if environment != "" {
		return environment
	}

	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	hostnameSplit := strings.Split(hostname, "-")

	if len(hostnameSplit) == 7 &&
		(hostnameSplit[2] == "20" || hostnameSplit[2] == "30") &&
		(hostnameSplit[5] == "ops" || hostnameSplit[5] == "pro") {
		environment = hostnameSplit[5]
		return environment
	}

	return ""
}

func IP() string {
	ip := Config().IP
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIp) > 0 {
		ip = LocalIp
	}

	return ip
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
