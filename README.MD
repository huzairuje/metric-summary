# CLI Summary Metric
## Introduction
This is Shell CLI Apps Summary Metrics.

## Getting Started
1. Usage with binary 
2. Usage with golang
3. Usage with parameters

### 1. Usage with binary
you can run the program itself by type this command, based on your operating system

#### a. windows binary
```shell
./bin/metric_summary_windows_amd64.exe <PARAMETER>
```

#### b. linux binary
```shell
./bin/metric_summary_linux_amd64 <PARAMETER>
```

#### c. MacOS (intel arch) binary 
```shell
./bin/metric_summary_macos_amd64 <PARAMETER>
```

#### c. MacOS (Apple Silicon arch) binary
```shell
./bin/metric_summary_macos_arm64 <PARAMETER>
```

### 2. Usage with golang
```shell
go run main.go <PARAMETER>
```

### 3. Parameter
example parameter
```shell
./bin/metric_summary_windows_amd64.exe -d metrics/json/ -t json -startTime 2021-12-12T00:00:00.00Z -endTime 2022-11-30T00:23:59.21Z
```

or like this 
```shell
./bin/metric_summary_windows_amd64.exe --duration metrics/json/ --type json -startTime 2021-12-12T00:00:00.00Z -endTime 2022-11-30T00:23:59.21Z
```

or this
```shell
go run main.go --duration metrics/json/ --type json -startTime 2021-12-12T00:00:00.00Z -endTime 2022-11-30T00:23:59.21Z
```

the default folder for output file is in `summary` folder, this folder is automatically created.