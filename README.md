[![Build Status](https://travis-ci.org/alesanfra/ground-control.svg?branch=master)](https://travis-ci.org/alesanfra/ground-control)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/63957071aa024e59accb6e9a628b4987)](https://www.codacy.com/project/alesanfra/ground-control/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=alesanfra/ground-control&amp;utm_campaign=Badge_Grade_Dashboard)
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


## Get data
Simple REST API
```
curl http://localhost:3000/devices 
```
