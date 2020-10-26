package ssh

import (
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

var hostKey ssh.PublicKey

func Execute(address string) {
	// Configure base client
	sshConfig := createSession("cian")
	fmt.Println("Created.")
	// Establish connection
	connection, err := establishConnection("tcp", address, sshConfig)
	fmt.Println("Connected.")
	session, err := prepareSession(connection)
	if err != nil {
		log.Fatalf("Failed to establish pty session: %v", err)
	}
	fmt.Println("Prepared.")
	// Execute Command
	//if err = session.Run("uname -a"); err != nil {
	//log.Fatalf("LOG: %v", err)
	//}

	// Close Session
	session.Close()
}

func createSession(user string) *ssh.ClientConfig {
	key, err := ioutil.ReadFile("/home/cian/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
}

func establishConnection(protocol, address string, config *ssh.ClientConfig) (*ssh.Client, error) {
	return ssh.Dial(protocol, fmt.Sprintf("%s:22"), config)
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
