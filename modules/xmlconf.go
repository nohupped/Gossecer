package modules

// Xmlconfig type is used to store the Server and Port attributes stored in the Ossec main configuration file.
type Xmlconfig struct {
	Ossec_config struct{
		Syslog_output struct{
			Server string `xml:"server"`
			Port int `xml:"port"`
			      }`xml:"syslog_output"`
		     } `xml:"ossec_config"`
}
