package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sensu/sensu-plugin-sdk"
	"github.com/sensu/sensu-go/types"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Username string
	Password string
	Database string
	Host string
	Port string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "mysql-health",
			Short:    "Simple mysql health check",
			Keyspace: "sensu.io/plugins/mysql-health/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "username",
			Argument:  "username",
			Shorthand: "u",
			Default:   "root",
			Usage:     "username.",
			Value:     &plugin.Username,
		},
		{
			Path:      "password",
			Argument:  "password",
			Shorthand: "p",
			Default:   "",
			Usage:     "password.",
			Value:     &plugin.Password,
		},
		{
			Path:      "database",
			Argument:  "database",
			Shorthand: "d",
			Default:   "mysql",
			Usage:     "database.",
			Value:     &plugin.Database,
		},
		{
			Path:      "host",
			Argument:  "host",
			Shorthand: "z",
			Default:   "127.0.0.1",
			Usage:     "host.",
			Value:     &plugin.Host,
		},
		{
			Path:      "port",
			Argument:  "port",
			Shorthand: "i",
			Default:   "3306",
			Usage:     "port.",
			Value:     &plugin.Port,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {

	dsl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", plugin.Username, plugin.Password, plugin.Host, plugin.Port, plugin.Database)
	db, err := sql.Open("mysql", dsl)
	if err != nil {
		fmt.Printf("%s CRITICAL: %s.\n", plugin.PluginConfig.Name, err.Error())
		return sensu.CheckStateCritical, nil
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		fmt.Printf("%s CRITICAL: %s.\n", plugin.PluginConfig.Name, err.Error())
		return sensu.CheckStateCritical, nil
	}

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		fmt.Printf("%s CRITICAL: %s.\n", plugin.PluginConfig.Name, err.Error())
		return sensu.CheckStateCritical, nil
	}

	fmt.Printf("%s OK: Server version: %s.\n", plugin.PluginConfig.Name, version)
	return sensu.CheckStateOK, nil
}
