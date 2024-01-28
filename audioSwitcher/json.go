package audioSwitcher

import (
	"encoding/json"
	"os"
	"strings"
)

func loadData() error {
	_, err := os.Stat("selectedDevice.json")
	if err != nil {
		err := os.WriteFile("selectedDevice.json", []byte("[]"), 0644)
		if err != nil {
			return err
		}
		return nil
	}

	data, err := os.ReadFile("selectedDevice.json")
	if err != nil {
		return err
	}

	SelectedDeviceData = []SelectedDevice{}

	err = json.Unmarshal(data, &SelectedDeviceData)
	if err != nil {
		return err
	}
	return nil
}

func saveData() error {
	for i := range SelectedDeviceData {
		SelectedDeviceData[i].Name = strings.ReplaceAll(SelectedDeviceData[i].Name, " (현재 사용중)", "")
	}

	data, err := json.MarshalIndent(SelectedDeviceData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("selectedDevice.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
