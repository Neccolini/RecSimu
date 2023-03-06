# RecSimu
RecSimu is a CLI tool for simulating the network of shape-changeable computer system.


## Requirements
- [Golang](https://go.dev/)


## Installation


### Via go install
```bash
$ go install github.com/Neccolini/RecSimu@latest
```


## Usage

### Generate Network Simulation Configuration file(json)
```bash
$ RecSimu gen -n <network_size> -c <simulation_cycles> -r <packet_injection_rate> -f <output_file_path>
```

### Run Simulation
```bash
$ RecSimu run -i <input_file_path>
```