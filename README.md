### syncwave

The intention of this project is to automate the setup and maintenance of a raspberry-pi k3s cluster with golang. This project could be used as a seperate terraform provider in the future.

Get list of node hostnames and/or addresses `nmap -sn 192.168.0.*` 

#### Build

To build the project, simple run the following command.
```bash
make build
```

Or you can build and run using the following command as well.
```bash
make build-run
```

#### Usage

You can interact with the syncwave tool using the following commands and flags.
```bash
Pass a hosts configuration file and scenario file in order to execute tasks on remote hosts.

Usage:
  syncwave execute [flags]

Flags:
  -h, --help   help for execute

Global Flags:
  -c, --config string     Pass configuration path/file.
  -s, --scenario string   Pass scenario path/file.
```

Below denotes the structure of the `config,yaml` and `scenario.yaml` files which  you can use as a base for building out your hosts and scenarios.

config.yaml
```yaml

all-nodes:
  master-nodes:
  worker-nodes:

master-nodes:
  hosts:
    - hostname: master-1
      ip-address: 192.168.0.1

worker-nodes:
  hosts:
    - hostname: worker-1
      ip-address: 192.168.0.3
    - hostname:  worker-2
      ip-address: 192.168.0.4

configuration:
  username: acid-burn
```

scenario.yaml
```yaml

scenario:
  name: Name of the scenario
  description: Description of the scenario
  tasks:
    - name: list os version
      exec: |
        uname -a
    - name: list ip address of remote machine
      exec: |
        hostname
    - name: test
      exec: |
        apt update -y
```

To execute a set of tasks on your hosts, pass in a scenario file as a flag, like so.
```bash
./syncwave execute -c config.yaml -s scenario.yaml
```

Examples of the structures of both files can be found in the repo above.

![Syncwave Sample Output](https://i.imgur.com/qY0KUKG.png)
