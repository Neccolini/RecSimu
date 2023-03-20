#!/bin/bash

for i in $(seq 0.01 0.01 0.1); do
    echo $i >> res
    go run main.go gen -f ./examples/example10.json -n 10 -r $i -c 100000
    go run main.go run -i ./examples/example10.json >> res
done

echo "next"
for i in $(seq 0.01 0.01 0.1); do
    echo $i >> res
    go run main.go gen -f ./examples/example10.json -n 20 -r $i -c 100000
    go run main.go run -i ./examples/example10.json >> res
done

echo "next"
for i in $(seq 0.01 0.01 0.1); do
    echo $i >> res
    go run main.go gen -f ./examples/example10.json -n 30 -r $i -c 100000
    go run main.go run -i ./examples/example10.json >> res
done

echo "next"
for i in $(seq 0.01 0.01 0.1); do
    echo $i >> res
    go run main.go gen -f ./examples/example10.json -n 40 -r $i -c 100000
    go run main.go run -i ./examples/example10.json >> res
done

echo "next"
for i in $(seq 0.01 0.01 0.1); do
    echo $i >> res
    go run main.go gen -f ./examples/example10.json -n 50 -r $i -c 100000
    go run main.go run -i ./examples/example10.json >> res
done