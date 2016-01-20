#!/bin/bash

mkdir -p log
nohup ./httpdns -addr 0.0.0.0:80 ./resolv.conf > ./log/httpdns.log 2>&1&
echo "httpdns server is running..."
