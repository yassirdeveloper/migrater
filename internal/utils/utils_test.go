package utils

import (
	"reflect"
	"testing"

	"github.com/yassirdeveloper/cli/errors"
)

func TestValidateSQLName(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want errors.Error
	}{
		{
			name: "valid name",
			s:    "valid_name",
			want: nil,
		},
		{
			name: "empty name",
			s:    "",
			want: errors.New("cannot be empty"),
		},
		{
			name: "name with spaces",
			s:    "name with spaces",
			want: errors.New("cannot include spaces"),
		},
		{
			name: "name with special characters",
			s:    "name@with#special$characters",
			want: errors.New("cannot include a special character"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateSQLName(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateSQLName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDSN(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		format  DSNFormat
		want    *DSN
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "valid dsn",
			dsn:    "user:password@tcp(localhost:5432)/database",
			format: DSNFormatMySQL,
			want: &DSN{
				Protocol: "tcp",
				Host:     "localhost",
				Port:     5432,
				Database: "database",
				User:     "user",
				Password: "password",
				format:   DSNFormatMySQL,
			},
			wantErr: false,
		},
		{
			name:   "valid dsn",
			dsn:    "user:password@localhost:5432/database",
			format: DSNFormatPostgres,
			want: &DSN{
				Host:     "localhost",
				Port:     5432,
				Database: "database",
				User:     "user",
				Password: "password",
				format:   DSNFormatPostgres,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDSN(tt.dsn, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDSN() got = %v, want %v", got, tt.want)
			}
		})
	}
}
