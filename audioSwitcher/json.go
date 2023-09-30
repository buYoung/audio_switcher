package audioSwitcher

import (
	"encoding/json"
	"os"
)

type SelectedDevice struct {
	Name  string
	Count int
}

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
	data, err := json.Marshal(SelectedDeviceData)
	if err != nil {
		return err
	}

	err = os.WriteFile("selectedDevice.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
