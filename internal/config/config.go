package config

import (
	"fmt"
)

type Config interface {
	DSN() string
}

type StandardConfig struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
}

func (c *StandardConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.User, c.Password, c.Password, c.Port, c.Database)
}

type FileConfig struct {
	filePath string
}

func (c *FileConfig) DSN() string {
	return c.filePath
}
