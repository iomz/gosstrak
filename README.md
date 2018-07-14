gosstrak-fc
==

[![Build Status](https://travis-ci.org/iomz/gosstrak.svg?branch=master)](https://travis-ci.org/iomz/gosstrak)
[![Coverage Status](https://coveralls.io/repos/iomz/gosstrak/badge.svg?branch=master)](https://coveralls.io/github/iomz/gosstrak?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/iomz/gosstrak)](https://goreportcard.com/report/github.com/iomz/gosstrak)
[![GoDoc](https://godoc.org/github.com/iomz/gosstrak?status.svg)](http://godoc.org/github.com/iomz/gosstrak)

Stat Monitoring
--
gosstrak collects statistical metrics and write them to InfluxDB for visualization in Grafana.
For quick start, use the docker-compose file and initialize the InfluxDB.

```bash
% sudo docker-compose up -d
% influx
>CREATE USER gosstrak WITH PASSWORD '<password>' WITH ALL PRIVILEGES
>CREATE DATABASE gosstrak
>exit
```

Then, run `gosstrak-fc` with `--enableStat` flag.

Author
--

Iori Mizutani (iomz)

License
--
See `LICENSE` file.
