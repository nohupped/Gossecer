package modules

import (
	"os"
	"io/ioutil"
	"encoding/xml"
	"gopkg.in/ini.v1"
)

// GetOSSecConfig returns the IP and Port at which Ossec is sending the syslog output
// after parsing the ossec's main configuration file.
func GetOSSecConfig(configfileparam string) (string, int) {
	xmlconfig, err := os.Open(configfileparam)
	CheckError(err)
	defer xmlconfig.Close()
	config := new(Xmlconfig)
	configbytes, err := ioutil.ReadAll(xmlconfig)
	CheckError(err)
	xml.Unmarshal(configbytes, &config.Ossec_config)
	return config.Ossec_config.Syslog_output.Server, config.Ossec_config.Syslog_output.Port

}

// GetConfig returns the ini file object of self's configuration file.
func GetConfig(configFile string) (*ini.File, error) {
	var Cfg *ini.File
	Cfg, err := ini.Load(configFile)
	if err != nil {
		return Cfg, err
	}
	return Cfg, err
}