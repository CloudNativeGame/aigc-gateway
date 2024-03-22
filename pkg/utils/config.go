package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type AigcGatewayConfig struct {
	Namespaces []string          `json:"namespaces"`
	GssLabels  map[string]string `json:"gss_labels"`
}

func ParseConfig(name string) (*AigcGatewayConfig, error) {
	// 打开 JSON 文件
	file, err := os.Open(name)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return nil, err
	}
	defer file.Close()

	// 读取文件内容
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return nil, err
	}

	// 将 JSON 数据解析到 Config 结构体中
	config := &AigcGatewayConfig{}
	if err := json.Unmarshal(bytes, config); err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return nil, err
	}

	// 输出解析结果
	fmt.Printf("Namespaces: %s\n", config.Namespaces)
	fmt.Printf("GSS Labels: %+v\n", config.GssLabels)
	return config, nil
}
