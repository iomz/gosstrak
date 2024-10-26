# gosstrak-fc

[![Test](https://github.com/iomz/gosstrak/actions/workflows/test.yml/badge.svg)](https://github.com/iomz/gosstrak/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iomz/gosstrak)](https://goreportcard.com/report/github.com/iomz/gosstrak)
[![codecov](https://codecov.io/gh/iomz/gosstrak/branch/main/graph/badge.svg?token=fN1tyc6ssX)](https://codecov.io/gh/iomz/gosstrak)
[![GoDoc](https://godoc.org/github.com/iomz/gosstrak?status.svg)](http://godoc.org/github.com/iomz/gosstrak)
[![License](https://img.shields.io/github/license/iomz/gosstrak.svg)](https://github.com/iomz/gosstrak/blob/main/LICENSE)

## Stat Monitoring

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

## TDT Benchmark

```bash
BenchmarkTranslate100Tags-32               10000            283151 ns/op           34096 B/op       1321 allocs/op
BenchmarkTranslate200Tags-32                5000            548079 ns/op           68752 B/op       2662 allocs/op
BenchmarkTranslate300Tags-32                3000            832507 ns/op          102801 B/op       3978 allocs/op
BenchmarkTranslate400Tags-32                2000           1151762 ns/op          137600 B/op       5341 allocs/op
BenchmarkTranslate500Tags-32                2000           1435293 ns/op          171312 B/op       6589 allocs/op
BenchmarkTranslate600Tags-32                2000           1769288 ns/op          204112 B/op       7906 allocs/op
BenchmarkTranslate700Tags-32                2000           2039621 ns/op          240305 B/op       9276 allocs/op
BenchmarkTranslate800Tags-32                1000           2256183 ns/op          274352 B/op      10614 allocs/op
BenchmarkTranslate900Tags-32                1000           2599413 ns/op          307761 B/op      11898 allocs/op
BenchmarkTranslate1000Tags-32               1000           2938569 ns/op          342385 B/op      13247 allocs/op
```

## Author

Iori Mizutani (iomz)

## License

See `LICENSE` file.
