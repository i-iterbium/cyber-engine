package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/i-iterbium/cyber-engine/internal/pkg/database"
)

// DBConfig описывает структуру конфигурации БД
type DBConfig struct {
	*database.Settings
	DefaultConnection string `json:"defaultConnection"`
}

// Config описывает структуру настроек подключения к БД
type Config struct {
	Database DBConfig `json:"db"`
}

var config interface{}
var filename string
var cfg Config

// Get возвращает сcылку на экземпляр cfg
func Get() *Config {
	return &cfg
}

// Init иницализирует настройку конфигурации
func Init() {
	cfg.Database.Settings = database.GetSettings()

	setConfig(&cfg)
	setFilename(filepath.Dir(os.Args[0]) + string(os.PathSeparator) + "config.json")

	if err := readFromFile(filename); err != nil {
		log.Println(err)
	}
}

// setConfig устанавливает ссылку на структуру конфигурации
func setConfig(c interface{}) {
	config = c
}

// SetFilename устанавливает путь к файлу с конфигурацией
func setFilename(fn string) {
	filename = fn
}

// readFromFile получает данные из файла настроек (config.json) и обновляет конфигурацию
func readFromFile(fn string) error {
	if fn == "" {
		return errors.New("Не передано название файла")
	}

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, &config)
}
