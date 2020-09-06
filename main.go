package main

import (
	"fmt"
	"io"
	"os"
  "sync"
  "golang.org/x/crypto/ssh"
)

func main() {
  var wg sync.WaitGroup

  // Create array of node ips
  var clusterNodes = []string{"192.168.0.103", "pi-worker-1", "pi-worker-2", "pi-worker-4-nfs"}
  wg.Add(1)

  go func() {
    for i := 0; i < len(clusterNodes); i++ {
      execute(clusterNodes[i])
    }
    wg.Done()
  }()

  wg.Wait()
}

func createSSHSession(user, password string) *ssh.ClientConfig {
  return &ssh.ClientConfig{
    User: user,
    Auth: []ssh.AuthMethod{
      ssh.Password(password),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
  }
}

func establishConnection(protocol, address string, config *ssh.ClientConfig) (*ssh.Client, error) {
  return ssh.Dial(protocol, fmt.Sprintf("%s:22", address), config)
}


func prepareSession(connection *ssh.Client) (*ssh.Session, error) {
  // Create a session to the remote terminal
  session, err := connection.NewSession()
  if err != nil {
    fmt.Println(fmt.Errorf("Failed to create session: %s", err))
  }

  // Create pseudo terminal (pty) in order to execute commands on remote machine
  modes := ssh.TerminalModes{
    ssh.ECHO: 0,              // Disable echoing
    ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
    ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
  }

  // Create xterm terminal that has 80 columns & 40 rows
  if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
    session.Close()
    fmt.Println(fmt.Errorf("Request for pseudo terminal failed: %s", err))
  }

  // Open stdout, stdin, stderr pipes
  stdin, err := session.StdinPipe()
  if err != nil {
    fmt.Println(fmt.Errorf("Unable to setup stdin session: %v", err))
  }
  // Copy to local stdin
  go io.Copy(stdin, os.Stdin)

  stdout, err := session.StdoutPipe()
  if err != nil {
    fmt.Println(fmt.Errorf("Unable to setup stdout session: %v", err))
  }
  go io.Copy(os.Stdout, stdout)

  stderr, err := session.StderrPipe()
  if err != nil {
    fmt.Println(fmt.Errorf("Unable to setup stderr session: %v", err))
  }
  go io.Copy(os.Stderr, stderr)

  return session, err
}

func execute(address string) {
  // Configure base client
  sshConfig := createSSHSession(os.Getenv("PI_USER"), os.Getenv("PI_PASSWORD"))

  // Establish connection
  connection, err := establishConnection("tcp", address, sshConfig)
  if err != nil {
    fmt.Println(fmt.Errorf("Failed to dial: %s", err))
  }

  session, err := prepareSession(connection)
  if err != nil {
    fmt.Errorf("Failed to establish pty session.", err)
  }

  err = session.Run("uname -a")
  // Ensure session gets closed
  session.Close()
}
