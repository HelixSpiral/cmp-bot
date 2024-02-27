package main

import (
	"encoding/json"
	"os"
)

func readCache(f string) (Cache, error) {
	var logTmp Cache
	rawdata, err := os.ReadFile(f)
	if err != nil {
		return Cache{}, err
	}

	err = json.Unmarshal(rawdata, &logTmp)
	if err != nil {
		return Cache{}, err
	}

	return logTmp, nil
}

func writeCache(f string, logTmp Cache) error {
	jsonData, err := json.Marshal(logTmp)
	if err != nil {
		return err
	}

	_ = jsonData
	err = os.WriteFile(f, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
