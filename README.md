[![Build Status](https://travis-ci.org/alesanfra/ground-control.svg?branch=master)](https://travis-ci.org/alesanfra/ground-control)
[![codecov](https://codecov.io/gh/alesanfra/ground-control/branch/master/graph/badge.svg)](https://codecov.io/gh/alesanfra/ground-control)
# Ground Control
Automated network scanner written in go

## Build
```
go build
```

## Run
Example on macOS
```
sudo ./ground-control -n 192.168.1.0/24
```
sudo is required in order to access ARP table


## Get data
Simple REST API
```
curl http://localhost:3000/devices 
```
