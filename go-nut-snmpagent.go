// go-nut-snmpagent
// Reads MIB.nut, connects to a NUT server, and exposes a SNMP API to
//
//	snmpd
//
// (c) 2023 by Gwyneth Llewelyn. All rights reserved.
// Distributed under a MIT license: https://gwyneth-llewelyn.mit-license.org/
package main

import (
	"fmt"
	"os"

	"github.com/robbiet480/go.nut"
	"gopkg.in/ini.v1"
)

// // Authentication object, which will be embedded later to reflect the INI structure.
// type GNSAuth struct {
// 	Username string
// 	Password string
// }

// GNSConfig ia a configuration object, to be filled in later.
type GNSConfig struct {
	NUTServer string `validate:"ip|hostname" ini:"hostname"`
//	*GNSAuth  `ini:"auth"`
	Username string
	Password string
}

// getFirstUPS connects to NUT, authenticates and returns the first UPS listed.
func getFirstUPS(config *GNSConfig) (*nut.UPS, error) {
	client, err := nut.Connect(config.NUTServer)
	if err != nil {
		return nil, fmt.Errorf("NUT connection error: %v", err)
	}
	// No authentication is perfectly valid, we just skip a step.
	if config.Username != "" {
		_, err = client.Authenticate(config.Username, config.Password)
		if err != nil {
			return nil, fmt.Errorf("NUT authentication error: %v", err)
		}
	} else {
		fmt.Println("[DEBUG] No authentication used to connect with NUT")
	}

	upsList, err := client.GetUPSList()
	if err != nil {
		return nil, fmt.Errorf("NUT getting UPS list error: %v", err)
	}
	fmt.Printf("\n[DEBUG]---------\nFirst UPS: %#v\n---------\n[/DEBUG]\n", upsList[0])
	return &upsList[0], nil
}

// Everything starts here.
func main() {
	// using LooseLoad so that we can simply ignore non-existing files.
	cfg, err := ini.LooseLoad("config.main.ini", "config.ini")
	// ,
	// 	ini.LoadOptions{
	// 		SkipUnrecognizableLines: true,
	// 	})
	if err != nil {
		fmt.Printf("error reading configuration file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[DEBUG] Value of configuration file is: %v\n", cfg)

	// place reasonable defaults into configuration struct.
	// config := &GNSConfig{
	// 	string("127.0.0.1"), // hostname
	// 	&GNSAuth{
	// 		Username: "",
	// 		Password: "",
	// 	},
	// }
	config := &GNSConfig{
		NUTServer: "127.0.0.1",
		Username: "",
		Password: "",
	}

	if err = cfg.MapTo(config); err != nil {
		fmt.Printf("invalid configuration, error was: %#v\n", err)
	}

	fmt.Printf("[DEBUG] Value of config object is: %#v\n", config)

	myUPS, err := getFirstUPS(config)
	if err != nil {
		fmt.Printf("could not read information of first UPS: %v\n", err)
		os.Exit(2)
	}
	fmt.Printf("Contact with UPS %q (%s) successful.\n", myUPS.Name, myUPS.Description)
	commands, err := myUPS.GetCommands()
	if err != nil {
		fmt.Printf("could not read list of commands for UPS %q: %v\n", myUPS.Name, err)
	} else {
		fmt.Printf("Commands available for UPS %q:\n", myUPS.Name)
		for _, cmd := range commands {
			fmt.Printf("%s (%s)\n", cmd.Name, cmd.Description)
		}
	}
}
