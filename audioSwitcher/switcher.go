package audioSwitcher

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var (
	re                 = regexp.MustCompile(ansi)
	deviceList         []map[string]interface{}
	SelectedDeviceData []SelectedDevice
)

func InitAudioSwitcher() {
	err := loadData()
	if err != nil {
		log.Panicln("eqweqw", err)
	}
	deviceList = GetPlaybackDevices()

	deviceList = lo.Map(deviceList, func(device map[string]interface{}, _ int) map[string]interface{} {
		isDefault := device["Default1"].(bool)
		index := device["Index"].(int)
		name := device["Name"].(string)

		if isDefault {
			name = fmt.Sprintf("%s (현재 사용중)", name)
		}
		return map[string]interface{}{
			"Index": index,
			"Name":  name,
		}
	})

	sort.Slice(deviceList, func(i, j int) bool {
		deviceA, findA := lo.Find(SelectedDeviceData, func(device SelectedDevice) bool {
			return device.Name == strings.ReplaceAll(deviceList[i]["Name"].(string), " (현재 사용중)", "")
		})
		deviceB, findB := lo.Find(SelectedDeviceData, func(device SelectedDevice) bool {
			return device.Name == strings.ReplaceAll(deviceList[j]["Name"].(string), " (현재 사용중)", "")
		})

		if findA && findB {
			return deviceA.Count > deviceB.Count
		}
		if findA {
			return true
		}
		if findB {
			return false
		}

		return deviceList[i]["Index"].(int) < deviceList[j]["Index"].(int)
	})

	var items []list.Item
	for _, q := range deviceList {
		items = append(items, item(q["Name"].(string)))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "변경할 장치를 선택해주세요"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func GetPlaybackDevices() []map[string]interface{} {
	powershell := exec.Command("pwsh", "-Command", "Get-AudioDevice", "-List")
	stdout, err := powershell.Output()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	buffurs := new(bytes.Buffer)
	transformString := transform.NewWriter(buffurs, korean.EUCKR.NewDecoder())
	_, err = transformString.Write(stdout)
	if err != nil {
		log.Panicln(err)
		return nil
	}
	err = transformString.Close()
	if err != nil {
		log.Panicln(err)
		return nil

	}

	pureDeviceListString := re.ReplaceAllString(buffurs.String(), "")
	filteredDeviceListString := lo.Reduce(strings.Split(pureDeviceListString, "\r\n\r\n"), func(acc []map[string]interface{}, device string, _ int) []map[string]interface{} {
		deviceDataList := map[string]interface{}{}
		for _, deviceData := range strings.Split(device, "\r\n") {
			if deviceData == "" {
				continue
			}
			deviceDataSplit := strings.Split(deviceData, ":")
			key := strings.TrimSpace(deviceDataSplit[0])
			value := strings.TrimSpace(deviceDataSplit[1])
			switch key {
			case "Default":
				deviceDataList["Default1"] = value == "True"
				break
			case "DefaultCommunication":
				deviceDataList[key] = value == "True"
				break
			case "Index":
				deviceDataList[key], _ = strconv.Atoi(value)
			default:
				deviceDataList[key] = value
			}
		}
		if val, ok := deviceDataList["Type"]; ok && val == "Playback" {
			acc = append(acc, deviceDataList)
		}
		return acc
	}, []map[string]interface{}{})

	return filteredDeviceListString
}

func SetPlaybackDevices(index int) error {
	powershell := exec.Command("pwsh", "-Command", fmt.Sprintf("Set-AudioDevice -Index %d", index))
	_, err := powershell.Output()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return nil
}
