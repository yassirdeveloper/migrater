package cmd

import (
	"fmt"
	"testing"

	"github.com/yassirdeveloper/cli/errors"
)

type MockOperator struct{}

func (m *MockOperator) Write(message string) errors.Error {
	fmt.Println(message)
	return nil
}
func (m *MockOperator) Read() (string, errors.Error) {
	return "", nil
}

func TestDescribeCommand(t *testing.T) {
	c := DescribeCommand()
	if c.String() != "describe" {
		t.Errorf("Expected command name 'describe', got '%s'", c.String())
	}
	input, err := c.Parse([]string{"--database", "test_db"})
	if err != nil {
		t.Errorf("Failed to parse command: %s", err)
	}
	err = c.Handle(input, &MockOperator{})
	if err != nil {
		t.Errorf("Failed to handle command: %s", err)
	}
}
