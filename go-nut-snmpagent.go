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

//	"github.com/posteo/go-agentx"
	"github.com/posteo/go-agentx/pdu"
//	"github.com/posteo/go-agentx/value"
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
	NUTserver string	`validate:"ip|hostname"`	// NUT server to contact.
//	*GNSAuth  `ini:"auth"`							// I have trouble parsing sections...
	Username string									// Username to connect to the NUT server.
	Password string									// Password for the connection with NUT server.
	SNMPserver string	`validate:"ip|hostname"`	// SNMP server to contact.
	SNMPport int		`validate:"integer"`		// AgentX port for SNMP server (default 705).
	SubagentOID string	`validate:""`				// OID for subagent to attach to (default PowerNet-MIB { iso org(3) dod(6) internet(1) private(4) enterprises(1) apc(318) }).
}

type oneOID struct {
	NUTvar string					// NUT UPS variable name.
	NUTvalue interface{}			// value expressed in different possible formats (string, integer. etc.)
	NUTtype pdu.VariableType		// essentially to tell what type that value comes from, as per SNMP specs.
	SNMPoid string					// OID for this mapping (NUT var name <-> OID).
	SNMPname string					// Name for this mapping (e.g. PowerNet-MIB::upsHighPrecExtdBatteryTemperature)
	SNMPdesc string					// SNMP description of this mapping.
}


// getFirstUPS connects to NUT, authenticates and returns the first UPS listed.
func getFirstUPS(config *GNSConfig) (*nut.UPS, error) {
	client, err := nut.Connect(config.NUTserver)
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
	// fmt.Printf("\n[DEBUG]---------\nFirst UPS: %#v\n---------\n[/DEBUG]\n", upsList[0])
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
	// fmt.Printf("[DEBUG] Value of configuration file is: %v\n", cfg)

	// place reasonable defaults into configuration struct.
	// config := &GNSConfig{
	// 	string("127.0.0.1"), // hostname
	// 	&GNSAuth{
	// 		Username: "",
	// 		Password: "",
	// 	},
	// }
	config := &GNSConfig{
		NUTserver: "127.0.0.1",				// localhost
		SNMPserver: "127.0.0.1",
		SNMPport: 705,						// Default for AgentX connection
		SubagentOID: ".1.3.6.1.4.1.318.1.1",	// PowerNet-MIB { iso org(3) dod(6) internet(1) private(4) enterprises(1) apc(318) products(1) hardware(1) }
		Username: "",
		Password: "",
	}

	if err = cfg.MapTo(config); err != nil {
		fmt.Printf("invalid configuration, error was: %#v\n", err)
	}

	// fmt.Printf("[DEBUG] Value of config object is: %#v\n", config)

	myUPS, err := getFirstUPS(config)
	if err != nil {
		fmt.Printf("could not read information of first UPS: %v\n", err)
		os.Exit(2)
	}
	fmt.Printf("Contact with UPS %q (%s) successful.\n", myUPS.Name, myUPS.Description)

/*	// debug stuff
	commands, err := myUPS.GetCommands()
	if err != nil {
		fmt.Printf("could not read list of commands for UPS %q: %v\n", myUPS.Name, err)
	} else {
		fmt.Printf("Commands available for UPS %q:\n", myUPS.Name)
		for _, cmd := range commands {
			fmt.Printf("%s (%s)\n", cmd.Name, cmd.Description)
		}
	}
	variables, err := myUPS.GetVariables()
	if err != nil {
		fmt.Printf("could not read list of variables for UPS %q: %v\n", myUPS.Name, err)
	} else {
		fmt.Printf("\nVariables available for UPS %q:\n", myUPS.Name)
		for _, vars := range variables {
			fmt.Printf("%s: %v (%s)\n", vars.Name, vars.Value, vars.Description)
		}
	}
*/
	// We ought to do some magic mapping between the MIB and the NUT variables... they do not exactly overlap

	var oneMap = &[]oneOID{
		{"battery.charge", 100, pdu.VariableTypeInteger, "", "", ""},
		{"battery.charge.low", 10, pdu.VariableTypeInteger, "", "", ""},
		{"battery.charge.warning", 50, pdu.VariableTypeInteger, "", "", ""},
		{"battery.mfr.date", "2022/02/18", pdu.VariableTypeOctetString, "", "", ""},	// date YYYY/MM/DD
		{"battery.runtime", 3720, pdu.VariableTypeInteger, ".1.1.1.2.2.3.0", "", "The UPS battery run time remaining before battery exhaustion."},
		{"battery.runtime.low", 120, pdu.VariableTypeInteger, "", "", ""},
		{"battery.temperature", 270, pdu.VariableTypeInteger, ".1.1.1.2.3.13.0", "PowerNet-MIB::upsHighPrecExtdBatteryTemperature", "The current internal UPS temperature expressed in tenths of degrees Celsius. Can be negative."},	// divide by 10
		{"battery.type", "PbAc", pdu.VariableTypeOctetString, "", "", ""},
		{"battery.voltage", 275, pdu.VariableTypeInteger, "", "", ""},		// divide by 10
		{"battery.voltage.nominal", 240, pdu.VariableTypeInteger, "", "", ""}, // divide by 10
		{"device.mfr", "American Power Conversion", pdu.VariableTypeOctetString, "", "", ""},
		{"device.model", "Smart-UPS 1000", pdu.VariableTypeOctetString, "", "", ""},
		{"device.serial", "AS0632330748", pdu.VariableTypeOctetString, "", "", ""},
		{"device.type", "ups", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.name", "usbhid-ups", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.parameter.bus", "001", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.parameter.pollfreq", 30, pdu.VariableTypeInteger, "", "", ""},
		{"driver.parameter.pollinterval", 2, pdu.VariableTypeInteger, "", "", ""},
		{"driver.parameter.port", "auto", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.parameter.product", "Smart-UPS 1000 FW:652.13.I USB FW:7.3", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.parameter.productid", "0002", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.parameter.serial", "AS0632330748", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.parameter.synchronous", "no", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.parameter.vendor", "American Power Conversion", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.parameter.vendorid", "051D", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.version", "2.7.4", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.version.data", "APC HID 0.96", pdu.VariableTypeOctetString, "", "", ""},
		{"driver.version.internal", "0.41", pdu.VariableTypeOctetString, "", "", ""},
		{"input.sensitivity", "medium", pdu.VariableTypeOctetString, "", "", ""},
		{"input.transfer.high", 2530, pdu.VariableTypeInteger, "", "", ""},	//divide by 10
		{"input.transfer.low", 2080, pdu.VariableTypeInteger, "", "", ""},// divide by 10
		{"input.voltage", 2260, pdu.VariableTypeInteger, ".1.1.1.3.2.1.0", "", "The current utility line voltage in VAC."},	// divide by 10
		{"output.current", 46, pdu.VariableTypeInteger, "", "", ""},	// divide by 100
		{"output.frequency", 500, pdu.VariableTypeInteger, "", "", ""},	// divide by 10
		{"output.voltage", 2260, pdu.VariableTypeInteger, ".1.1.1.4.2.1.0", "", "The output voltage of the UPS system in VAC."},	// divide by 10
		{"output.voltage.nominal", 2300, pdu.VariableTypeInteger, "", "", ""},	// divide by 10
		{"ups.beeper.status", "disabled", pdu.VariableTypeOctetString, "", "", ""},
		{"ups.delay.shutdown", 20, pdu.VariableTypeInteger, "", "", ""},
		{"ups.delay.start", 30, pdu.VariableTypeInteger, "", "", ""},
		{"ups.firmware", "652.13.I", pdu.VariableTypeOctetString, "", "", ""},
		{"ups.firmware.aux", "7.3", pdu.VariableTypeOctetString, "", "", ""},
		{"ups.load", 169, pdu.VariableTypeGauge32, ".1.1.1.4.2.3.0", "", "The current UPS load expressed in percent of rated capacity."},		// divide by 10
		{"ups.mfr", "American Power Conversion", pdu.VariableTypeOctetString, "", "", ""},
		{"ups.mf22r.date", "2006/08/03", pdu.VariableTypeOctetString, "", "", ""},	// date YYYY/MM/DD
		{"ups.model", "Smart-UPS 1000", pdu.VariableTypeOctetString, "", "", ""},
		{"ups.productid", "0002", pdu.VariableTypeOctetString, "", "", ""},
		{"ups.serial", "AS0632330748", pdu.VariableTypeOctetString, "", "", ""},
		{"ups.status", "OL", pdu.VariableTypeOctetString, "", "", ""},
		{"ups.test.result", "No test initiated", pdu.VariableTypeOctetString, "", "", ""},
		{"ups.timer.reboot", -1, pdu.VariableTypeInteger, "", "", ""},
		{"ups.timer.shutdown", -1, pdu.VariableTypeInteger, "", "", ""},
		{"ups.timer.start", -1, pdu.VariableTypeInteger, "", "", ""},
		{"ups.vendorid", "051d", pdu.VariableTypeOctetString, "", "", ""},
	}
	fmt.Println("NUT var name\tValue\tType\tOID Name Description")
	for _, apcRawData := range *oneMap {
		fmt.Printf("%s:\t%v\t(%s)\t=>%s %q %q\n", apcRawData.NUTvar, apcRawData.NUTvalue, apcRawData.NUTtype.String(), apcRawData.SNMPoid, apcRawData.SNMPname, apcRawData.SNMPdesc)
	}
}
