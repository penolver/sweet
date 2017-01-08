package sweet

import (
	"fmt"
)

type Unix struct {
}

func newUnixCollector() Collector {
	return Unix{}
}

func (collector Unix) Collect(device DeviceConfig) (map[string]string, error) {
	result := make(map[string]string)

	c, err := newSSHCollector(device)
	if err != nil {
		return result, fmt.Errorf("Error connecting to %s: %s", device.Hostname, err.Error())
	}

	if err := expect("assword:", c.Receive); err != nil {
		return result, fmt.Errorf("Missing password prompt: %s", err.Error())
	}
	// assuming we have a root priv.. (# prompt).. we'll need to fix this and provide options..
	c.Send <- device.Config["pass"] + "\n"
	multi := []string{"#", "assword:"}
	m, err := expectMulti(multi, c.Receive)
	if err != nil {
		return result, fmt.Errorf("Invalid response to password: %s", err.Error())
	}
	if m == "assword:" {
		return result, fmt.Errorf("Bad username or password.")
	}
	//c.Send <- "uname -a\n"
	//if err := expect("#", c.Receive); err != nil {
	//	return result, fmt.Errorf("Command 'uname -a' failed: %s", err.Error())
	//}
	c.Send <- "uname -a\n"
	result["version"], err = expectSaveTimeout("#", c.Receive, device.CommandTimeout)
	if err != nil {
		return result, fmt.Errorf("Command 'uname -a' failed: %s", err.Error())
	}
	// example of an important file to backup
	c.Send <- "cat /etc/passwd\n"
	result["config"], err = expectSaveTimeout("#", c.Receive, device.CommandTimeout)
	if err != nil {
		return result, fmt.Errorf("Command 'cat /etc/passwd' failed: %s", err.Error())
	}
	// example of an important file to backup
	c.Send <- "cat /etc/passwd\n"
	result["passwd"], err = expectSaveTimeout("#", c.Receive, device.CommandTimeout)
	if err != nil {
		return result, fmt.Errorf("Command 'cat /etc/passwd' failed: %s", err.Error())
	}
	c.Send <- "exit\n"

	return result, nil
}
