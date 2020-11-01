### go-ssh

The intention of this project is to automate the setup and maintenance of a raspberry-pi k3s cluster with golang. This project could be used as a seperate terraform provider in the future.

Get list of node hostnames and/or addresses `nmap -sn 192.168.0.*` 

#### Build

To build the project, simple run the following command.
```bash
go build -o syncwave ./cmd/syncwave
```

#### Run

To run the project, execute the compiled binary and pass in a configuration files with a list of your hosts.
```bash
./syncwave execute --config-file config.yaml
```
