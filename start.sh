#!/bin/bash

# ./miner &
echo "Miner Start Script"
cd /home/berry/mine/claymore && nohup ./start.bash > logs.txt &
