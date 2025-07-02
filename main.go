package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	URL         string            `yaml:"url"`
	IntervalSec int               `yaml:"interval_sec"`
	DataFormat  map[string]string `yaml:"data_format"`
}

func loadConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func generateRandomData(template map[string]string) map[string]interface{} {
	data := make(map[string]interface{})
	for key, valType := range template {
		switch valType {
		case "int":
			data[key] = rand.Intn(100)
		case "float":
			data[key] = rand.Float64() * 100
		case "string":
			data[key] = fmt.Sprintf("str_%d", rand.Intn(10000))
		case "bool":
			data[key] = rand.Intn(2) == 1
		default:
			data[key] = nil
		}
	}
	return data
}

func sendJSON(url string, jsonData map[string]interface{}) error {
	body, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 2 {
		log.Fatal("Usage: ./client config.yaml")
	}

	configPath := os.Args[1]
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ticker := time.NewTicker(time.Duration(config.IntervalSec) * time.Second)
	defer ticker.Stop()

	log.Println("Starting REST client...")
	for {
		select {
		case <-ticker.C:
			data := generateRandomData(config.DataFormat)
			err := sendJSON(config.URL, data)
			if err != nil {
				log.Printf("Failed to send data: %v", err)
			} else {
				log.Printf("Sent data: %+v", data)
			}
		}
	}
}
