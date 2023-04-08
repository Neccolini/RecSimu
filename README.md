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
$ RecSimu gen -t "<topology> <node_num or {rows and colunms}>" -c <simulation_cycles> -r <packet_injection_rate> -f <output_file_path>
```
for example...
```bash
$ RecSimu gen -f ./examples/example20.json -c 10000 -t "random 20"  -r 0.01
```
```bash
$ RecSimu gen -f ./examples/example3_5.json -c 10000 -t "mesh 3 5"  -r 0.01
```
### Run Simulation
```bash
$ RecSimu run -i <input_file_path>
```