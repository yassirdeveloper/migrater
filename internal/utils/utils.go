package utils

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/yassirdeveloper/cli/errors"
)

type DSNFormat string

const (
	DSNFormatMySQL     DSNFormat = "mysql"
	DSNFormatPostgres  DSNFormat = "postgres"
	DSNFormatSQLServer DSNFormat = "sqlserver"
	DSNFormatSQLite    DSNFormat = "sqlite"
)

type DSN struct {
	Protocol string
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
	format   DSNFormat
}

func (c *DSN) String() string {
	switch c.format {
	case DSNFormatMySQL:
		return fmt.Sprintf("%s:%s@%s(%s:%d)/%s", c.User, c.Password, c.Protocol, c.Host, c.Port, c.Database)
	case DSNFormatPostgres:
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.Database)
	case DSNFormatSQLServer:
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", c.User, c.Password, c.Host, c.Port, c.Database)
	case DSNFormatSQLite:
		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.Database)
	default:
		return ""
	}
}

func ToDSN(dsn string, format DSNFormat) (*DSN, errors.Error) {
	conf := &DSN{}
	switch format {
	case DSNFormatMySQL:
		// user:pass@tcp(host:3306)/dbname
		atIdx := strings.Index(dsn, "@")
		if atIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing '@'")
		}
		userPass := dsn[:atIdx]
		rest := dsn[atIdx+1:]

		colonIdx := strings.Index(userPass, ":")
		if colonIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing ':' in user:pass")
		}
		conf.User = userPass[:colonIdx]
		conf.Password = userPass[colonIdx+1:]

		openIdx := strings.Index(rest, "(")
		closeIdx := strings.Index(rest, ")")
		if openIdx == -1 || closeIdx == -1 || closeIdx < openIdx {
			return nil, errors.New("Could not parse DSN: invalid protocol/address format")
		}
		conf.Protocol = rest[:openIdx]
		address := rest[openIdx+1 : closeIdx]
		after := rest[closeIdx+1:]

		colonIdx = strings.LastIndex(address, ":")
		if colonIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing ':' in host:port")
		}
		conf.Host = address[:colonIdx]
		var port uint16
		_, err := fmt.Sscanf(address[colonIdx+1:], "%d", &port)
		if err != nil {
			return nil, errors.New("Could not parse DSN: invalid port")
		}
		conf.Port = port

		if !strings.HasPrefix(after, "/") {
			return nil, errors.New("Could not parse DSN: missing '/' before database")
		}
		conf.Database = after[1:]
		conf.format = DSNFormatMySQL

	case DSNFormatPostgres:
		// user:pass@host:port/dbname
		atIdx := strings.Index(dsn, "@")
		if atIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing '@'")
		}
		userPass := dsn[:atIdx]
		rest := dsn[atIdx+1:]

		colonIdx := strings.Index(userPass, ":")
		if colonIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing ':' in user:pass")
		}
		conf.User = userPass[:colonIdx]
		conf.Password = userPass[colonIdx+1:]

		slashIdx := strings.Index(rest, "/")
		if slashIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing '/' before database")
		}
		hostPort := rest[:slashIdx]
		dbName := rest[slashIdx+1:]

		colonIdx = strings.LastIndex(hostPort, ":")
		if colonIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing ':' in host:port")
		}
		conf.Host = hostPort[:colonIdx]
		var port uint16
		_, err := fmt.Sscanf(hostPort[colonIdx+1:], "%d", &port)
		if err != nil {
			return nil, errors.New("Could not parse DSN: invalid port")
		}
		conf.Port = port
		conf.Database = dbName
		conf.format = DSNFormatPostgres

	case DSNFormatSQLServer:
		// user:pass@host:port?database=dbname
		colonIdx := strings.Index(dsn, ":")
		if colonIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing ':' in user:pass")
		}
		conf.User = dsn[:colonIdx]
		passwordAtIdx := strings.Index(dsn, "@")
		if passwordAtIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing '@'")
		}
		conf.Password = dsn[colonIdx+1 : passwordAtIdx]
		after := dsn[passwordAtIdx+1:]

		slashIdx := strings.Index(after, "/")
		if slashIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing '/' before database")
		}
		questionIdx := strings.Index(after, "?")
		if questionIdx == -1 {
			return nil, errors.New("Could not parse DSN: missing '?' before database")
		}
		conf.Host = after[:slashIdx]
		portStr := after[slashIdx+1 : questionIdx]
		var port uint16
		_, err := fmt.Sscanf(portStr, "%d", &port)
		if err != nil {
			return nil, errors.New("Could not parse DSN: invalid port")
		}
		conf.Port = port
		conf.Database = strings.TrimPrefix(after[questionIdx+1:], "database=")
		conf.format = DSNFormatSQLServer

	case DSNFormatSQLite:
		// file:dbname?cache=shared&mode=rwc
		if !strings.HasPrefix(dsn, "file:") {
			return nil, errors.New("Could not parse DSN: missing 'file:' prefix")
		}
		conf.format = DSNFormatSQLite
		dbName := strings.TrimPrefix(dsn, "file:")
		if strings.Contains(dbName, "?") {
			parts := strings.SplitN(dbName, "?", 2)
			conf.Database = parts[0]
			// Handle additional parameters if needed
			// For SQLite, we can ignore the parameters as they are not critical for connection
		} else {
			conf.Database = dbName
		}

	default:
		return nil, errors.New("Unsupported DSN format")
	}

	if conf.Database == "" {
		return nil, errors.New("Could not parse DSN: missing database name")
	}

	return conf, nil
}

func ValidateSQLName(s string) errors.Error {
	s = strings.TrimSpace(s)

	if s == "" {
		return errors.New("cannot be empty")
	}

	if unicode.IsDigit(rune(s[0])) {
		return errors.New("cannot start with a digit")
	}

	for _, r := range s {
		if unicode.IsSpace(r) {
			return errors.New("cannot include spaces")
		}
		// Allow letters, digits, and underscores
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return errors.New("cannot include a special character")
		}
	}

	reservedKeywords := []string{"SELECT", "FROM", "WHERE", "AND", "OR", "NOT", "IN", "IS", "NULL", "TRUE", "FALSE"}
	for _, keyword := range reservedKeywords {
		if strings.ToUpper(s) == keyword {
			return errors.New("cannot be a reserved keyword")
		}
	}

	return nil
}
