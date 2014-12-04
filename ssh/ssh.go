package ssh

import (
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// Config for the SSH connection
type Config struct {
	Username string
	Password string
	Host     string
	Port     int
}

// Copy the output from a command to the specified io.Writer
func (config *Config) Copy(dest io.Writer, src io.Reader) error {
	// TODO: error if port <= 0
	clientconfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), clientconfig)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	session.Stdout = dest
	command, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	if err := session.Run(string(command[:])); err != nil {
		return err
	}

	return nil
}
