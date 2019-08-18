package util

import (
	"errors"
	"fmt"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

type SshClient struct {
	clients []*ssh.Client
}

func NewSshClient(host string) (*SshClient, error) {
	out := new(SshClient)
	for i := 10; i > 0; i -= 5 {
		client, err := sshConnect(host)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		out.clients = append(out.clients, client)
	}
	return out, nil
}
func (this SshClient) GetSession() (*ssh.Session, error) {
	for _, client := range this.clients {
		session, err := client.NewSession()
		if err != nil {
			continue
		}
		return session, nil
	}
	return nil, errors.New("Unable to get a session")
}

/**
 * Easy shorthand for multiple calls to sshExec
 * @param  ...string    commands    The commands to execute
 * @return []string                 The results of the execution of each command
 */
func (this SshClient) MultiRun(commands ...string) ([]string, error) {
	out := []string{}
	for _, command := range commands {
		res, err := this.Run(command)
		if err != nil {
			return nil, err
		}
		out = append(out, string(res))
	}
	return out, nil
}

/**
 * Speeds up remote execution by chaining commands together
 * @param  ...string    commands    The commands to execute
 * @return string                   The result of the execution
 */
func (this SshClient) FastMultiRun(commands ...string) (string, error) {
	cmd := ""
	for i, command := range commands {
		if i != 0 {
			cmd += "&&"
		}
		cmd += command
	}
	return this.Run(cmd)
}
func (this SshClient) Run(command string) (string, error) {
	session, err := this.GetSession()
	defer session.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	out, err := session.CombinedOutput(command)
	return string(out), err
}

func (this SshClient) DockerExec(node int, command string) (string, error) {
	return this.Run(fmt.Sprintf("docker exec whiteblock-node%d %s", node, command))
}
func (this SshClient) DockerExecd(node int, command string) (string, error) {
	return this.Run(fmt.Sprintf("docker exec -d whiteblock-node%d %s", node, command))
}
func (this SshClient) DockerRead(node int, file string) (string, error) {
	return this.Run(fmt.Sprintf("docker exec whiteblock-node%d cat %s", node, file))
}
func (this SshClient) DockerMultiExec(node int, commands []string) (string, error) {
	merged_command := ""
	for _, command := range commands {
		if len(merged_command) != 0 {
			merged_command += "&&"
		}
		merged_command += fmt.Sprintf("docker exec -d whiteblock-node%d %s", node, command)
	}
	return this.Run(merged_command)
}

/**
 * Copy a file over to a remote machine
 * @param  string   host    The IP address of the remote host
 * @param  string   src     The source path of the file
 * @param  string   dest    The destination path of the file
 */
func (this SshClient) Scp(src string, dest string) error {
	session, err := this.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()
	err = scp.CopyPath(src, dest, session)
	if err != nil {
		return err
	}
	return nil
}

func (this SshClient) Close() {
	for _, client := range this.clients {
		if client == nil {
			continue
		}
		client.Close()
	}
}

func sshConnect(host string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(conf.SSHPrivateKey)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), sshConfig)
	if err != nil {
		fmt.Println("First ssh attempt failed: " + err.Error())
	}

	return client, err
}
