package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const configFile = ".ggrok"

type Config struct {
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) SaveToConfig(proxyServer string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("get home dir error when save config")
	}
	filePath := home + "/" + configFile
	if !isExist(filePath) {
		os.Create(filePath)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	write.WriteString(strings.TrimSpace(proxyServer))

	write.Flush()
}

func (c *Config) ReadConfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("get home dir error when read config")
	}
	filePath := home + "/" + configFile
	if !isExist(filePath) {
		os.Create(filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(content))
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}
