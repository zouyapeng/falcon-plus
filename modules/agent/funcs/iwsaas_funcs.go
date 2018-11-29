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
)

func IWSaaSBootstrapErrorMetrics() []*model.MetricValue {
	stage1Error := 0
	stage2Error := 0

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

func IWSaaSPatchHistoryMetrics() []*model.MetricValue {
	hasError, err := readPatchHistory()
	if err != nil {
		return nil
	}

	return []*model.MetricValue{
		GaugeValue("iwsaas.patch.history.error", hasError),
	}
}
