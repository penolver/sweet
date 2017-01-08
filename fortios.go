package sweet

import (
	"fmt"
)

type FortiOS struct {
}

func newFortiOSCollector() Collector {
	return FortiOS{}
}

func (collector FortiOS) Collect(device DeviceConfig) (map[string]string, error) {
	result := make(map[string]string)

	c, err := newSSHCollector(device)
	if err != nil {
		return result, fmt.Errorf("Error connecting to %s: %s", device.Hostname, err.Error())
	}

	if err := expect("assword:", c.Receive); err != nil {
		return result, fmt.Errorf("Missing password prompt: %s", err.Error())
	}

	c.Send <- device.Config["pass"] + "\n"
	multi := []string{"#", "assword:"}
	m, err := expectMulti(multi, c.Receive)
	if err != nil {
		return result, fmt.Errorf("Invalid response to password: %s", err.Error())
	}
	if m == "assword:" {
		return result, fmt.Errorf("Bad username or password.")
	}
	//c.Send <- "get system status\n"
	//if err := expect("#", c.Receive); err != nil {
	//	return result, fmt.Errorf("Command 'get system status' failed: %s", err.Error())
	//}
	c.Send <- "get system status\n"
	result["version"], err = expectSaveTimeout("#", c.Receive, device.CommandTimeout)
	if err != nil {
		return result, fmt.Errorf("Command 'get system status' failed: %s", err.Error())
	}
	// we probably should check re vdoms.. but for now.. this will run a backup aganst global..
	// config global
	// or config vdom xxx
	// config system console
	// set output standard
	// end
	// show full-configuration
	c.Send <- "config global \n"
	_ = expect("#", c.Receive)
  c.Send <- "config system console \n"
	_ = expect("#", c.Receive)
	c.Send <- "set output standard \n"
	_ = expect("#", c.Receive)
	c.Send <- "end \n"
	_ = expect("(global) #", c.Receive)
	c.Send <- "show full-configuration \n"
	result["config"], err = expectSaveTimeout("(global) #", c.Receive, device.CommandTimeout)
	if err != nil {
		return result, fmt.Errorf("Command 'show full-configuration' failed: %s", err.Error())
	}
	c.Send <- "exit\n"

	return result, nil
}
