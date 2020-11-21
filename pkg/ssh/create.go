package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func Execute(address, task string) (string, error) {
	// Current user
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Could not get current user: %v", err)
	}

	// Configure base client
	sshConfig := createSession(user)

	// Establish connection
	connection, err := establishConnection("tcp", address, sshConfig)
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}

	session, err := prepareSession(connection)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Execute command
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	status := session.Run(task)

	// Close Session
	session.Close()

	return stdoutBuf.String(), status
}

func createHostkeyCallback(user *user.User) ssh.HostKeyCallback {
	hostkeyCallback, err := knownhosts.New(fmt.Sprintf("%s/.ssh/known_hosts", user.HomeDir))
	if err != nil {
		log.Fatalf("Could not read hosts file: %v", err)
	}

	return hostkeyCallback
}

func createSession(user *user.User) *ssh.ClientConfig {
	hostkeyCallback := createHostkeyCallback(user)

	key, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa", user.HomeDir))
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	return &ssh.ClientConfig{
		User: user.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostkeyCallback,
	}
}

func establishConnection(protocol, address string, config *ssh.ClientConfig) (*ssh.Client, error) {
	client, err := ssh.Dial(protocol, fmt.Sprintf("%s:22", address), config)
	if err != nil {
		log.Fatalf("Unable to dial connection: %v", err)
	}

	return client, err
}

func prepareSession(connection *ssh.Client) (*ssh.Session, error) {
	session, err := connection.NewSession()
	if err != nil {
		log.Fatalf("Session error: %v", err)
	}

	return session, err
}

func keySigner() (signer ssh.Signer) {
	key, err := ioutil.ReadFile("/home/cian/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("Unable to read ssh private key: %v", err)
	}

	signer, err = ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	return
}
