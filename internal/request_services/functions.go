package request_services

import "encoding/json"


func convertInterfaceToStruct(m interface{}, s interface{}) error {
	// convert map to json
	jsonString, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// convert json to struct
	err = json.Unmarshal(jsonString, &s)
	if err != nil {
		return err
	}

	return nil
}
