[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=alesanfra_ground-control&metric=alert_status)](https://sonarcloud.io/dashboard?id=alesanfra_ground-control)
[![codecov](https://codecov.io/gh/alesanfra/ground-control/branch/master/graph/badge.svg)](https://codecov.io/gh/alesanfra/ground-control)
# Ground Control - network monitor
Ground control is an automated network monitor written in go from scratch.

It features an arp scanner to monitor the presence of hosts in your local network 
and a speed test based on [Mlab](https://www.measurementlab.net).

## Build

You can build **ground-control** from sources with the standard `go build` command. 

In order to do so you have to install the golang toolchain, version 1.12 or newer. 

```
go build
```

## Run

In order to perform ARP scan **ground-control** needs to access to the ARP table of your OS, 
you have to start it with administrative privileges.

Example on macOS:
```
sudo ./ground-control
```
sudo is required in order to access ARP table


## Get data
Simple REST API
```
curl http://localhost:3000/api/v1/devices 
```
