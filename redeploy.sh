#!/bin/bash

# To get process ID run
# ps -aux | grep "go run ."
# Or, if run by nohup, then
# ps -ef | grep "go"
# then kill the processes with "kill -9 {PID}"

nohup go run . &
