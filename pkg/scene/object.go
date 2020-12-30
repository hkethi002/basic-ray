package scene

import (
	geometry "basic-ray/pkg/geometry"
	"encoding/json"
	"io/ioutil"
)

func ReadObject(filename string) (*geometry.Object, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// json data
	var obj geometry.Object

	// unmarshall it
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func WriteObject(filename string, object *geometry.Object) error {
	fileData, err := json.MarshalIndent(*object, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, fileData, 0644)
}
