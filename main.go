package main

import (
  "fmt"
  "io"
  "os"
  "golang.org/x/crypto/ssh"
)

func main() {
  // Configure base client
  sshConfig := &ssh.ClientConfig{
    User: os.Getenv("PI_USER"),
    Auth: []ssh.AuthMethod{
      ssh.Password(os.Getenv("PI_PASSWORD")),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
  }

  // Establish connection
  connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", os.Getenv("PI_ADDRESS")), sshConfig)
  if err != nil {
    fmt.Println(fmt.Errorf("Failed to dial: %s", err))
  }

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

  err = session.Run("uname -a")

  // Ensure session gets closed
  session.Close()
}
