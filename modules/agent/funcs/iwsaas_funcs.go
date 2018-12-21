package funcs

import (
	"bufio"
	"bytes"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/file"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"net/http"
	"time"
	"log"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func IWSaaSBootstrapErrorMetrics() (L []*model.MetricValue) {
	stage1Error := 0
	stage2Error := 0

	if _, err := os.Stat("/opt/icsops/bootstrap/result"); os.IsNotExist(err) {
		return
	}

	if _, err := os.Stat("/opt/icsops/bootstrap/result/stage1.error"); !os.IsNotExist(err) {
		stage1Error = 1
	}

	if _, err := os.Stat("/opt/icsops/bootstrap/result/stage2.error"); !os.IsNotExist(err) {
		stage2Error = 1
	}

	return []*model.MetricValue{
		GaugeValue("iwsaas.bootstrap.stage1", stage1Error),
		GaugeValue("iwsaas.bootstrap.stage2", stage2Error),
	}
}

func readPatchHistory() (int, error) {
	contents, err := ioutil.ReadFile("/opt/icsops/patch/history.ini")
	if err != nil {
		return 0, err
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	for {
		line, err := file.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return 0, err
		}

		lineString := strings.Split(strings.Replace(string(line), " ", "", -1), "=")
		for index, value := range lineString{
			if index == 0 && (value == "version" || value == "time" || value == "reason"){
				continue
			} else if index != 0{
				if value != ""{
					return 1, nil
				}
			} else {
				break
			}
		}
	}

	return 0, nil
}

func IWSaaSPatchHistoryMetrics() (L []*model.MetricValue) {
	if _, err := os.Stat("/opt/icsops/patch"); os.IsNotExist(err) {
		return
	}

	hasError, err := readPatchHistory()
	if err != nil {
		return
	}

	return []*model.MetricValue{
		GaugeValue("iwsaas.patch.history.error", hasError),
	}
}

func getCurrentIP() (publicIP string) {
	httpClient := &http.Client{}
	httpClient.Timeout = 3 * time.Second

	resp, err := httpClient.Get("http://169.254.169.254/latest/meta-data/public-ipv4")
	if err != nil {
		log.Println("ERROR: Get public-ipv4 from AWS fail", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200{

		Response, _ := ioutil.ReadAll(resp.Body)
		return string(Response)
	}
	log.Println("WARN: Instance do not have public ip")

	return
}

func IWSaaSEIPIsExistMetrics() (L []*model.MetricValue) {
	reportEipIsExist := g.ReportEipIsExist()
	reportCurrentRegionEipIsExist := reportEipIsExist[g.Config().Environment][g.Config().Region]

	if len(reportCurrentRegionEipIsExist) == 0{
		return
	}

	publicIP := getCurrentIP()
	if publicIP == ""{
		return
	}
	isExist := 0
	for _,ip := range reportCurrentRegionEipIsExist{
		if ip == publicIP{
			isExist = 1
			break
		}
	}

	return []*model.MetricValue{
		GaugeValue(g.EIP_IS_EXIST, isExist),
	}
}
