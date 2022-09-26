package converter

import "encoding/json"

func Unmarshal(data interface{}, target interface{}) error {
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawData, &target)
	if err != nil {
		return err
	}

	return err
}
