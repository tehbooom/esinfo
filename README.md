# esinfo

[![Go Report Card](https://goreportcard.com/badge/github.com/tehbooom/project_name)](https://goreportcard.com/report/github.com/tehbooom/esinfo)

![Build Status](https://github.com/tehbooom/esinfo/actions/workflows/build.yml/badge.svg)

When running large elasticsearch clusters it can be difficult to know what indexes you have in the cluster without manually searching through index management or using dev tools and scrolling. `esinfo` queries elasticsearch for all indexes in the cluster and outputs them in a nice format (csv, json, yaml).

## Installation

```bash
git clone https://github.com/tehbooom/esinfo.git
cd esinfo
make build
```

Or download the latest binary [here](https://github.com/tehbooom/esinfo/releases) for your OS

## Options

| Flag     | Description                        | Default        |
|----------|------------------------------------|----------------|
| endpoint | url to elasticsearch API           | localhost:9200 |
| username | username to authenticate           | elastic        |
| password | password for user                  | changeme       |
| cacert   | path to certificate authority file |                |
| unsafe   | option to verify ssl               | true           |
| format   | output format for esinfo           | csv            |

## Usage

`esinfo` looks for a config file at `esinfo.yaml` in the directory you are currently at or your `$HOME` directory

[Here](esinfo.yaml) is an example file

Optionally you can set any flag at the commandline which overrules anything in your config file

```bash
esinfo [command] [flags]
```

## Examples

```bash
$ esinfo test
{
  "name" : "instance-0000000001",
  "cluster_name" : "ee0073bc7ffc44a8b62454c1a73f508e",
  "cluster_uuid" : "F37aL59HQZm41W4BBDWrbg",
  "version" : {
    "number" : "8.8.0",
    "build_flavor" : "default",
    "build_type" : "docker",
    "build_hash" : "c01029875a091076ed42cdb3a41c10b1a9a5a20f",
    "build_date" : "2023-05-23T17:16:07.179039820Z",
    "build_snapshot" : false,
    "lucene_version" : "9.6.0",
    "minimum_wire_compatibility_version" : "7.17.0",
    "minimum_index_compatibility_version" : "7.0.0"
  },
  "tagline" : "You Know, for Search"
}

Connection successful!
```

```bash
$ esinfo run
CSV file located at /home/alec/go/src/github.com/tehbooom/esinfo/indices.csv
```

```bash
$ esinfo run -f json -p "supersecretpassword" -u elastic -e "https://es-1:9200"
JSON file located at /home/alec/go/src/github.com/tehbooom/esinfo/indices.json
```

## Man Page

```text
When running large elasticsearch clusters it can be difficult to know what indexes you have in the cluster without manually 
        searching through index management or using dev tools and scrolling. Esinfo queries elasticsearch for all indexes in the cluster and 
        outputs them in a nice format(csv, json, yaml).

Usage:
  esinfo [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Queries the elasticsearch cluster for all indices and datastreams and outputs them.
  test        A brief description of your command

Flags:
      --cacert string     Certificate Authority for cluster
      --config string     config file (default is esinfo.yaml)
  -e, --endpoint string   Address to elasticsearch (default "localhost:9200")
  -f, --format string     Output type for file (default "csv")
  -h, --help              help for esinfo
  -p, --password string   Password for elasticsearch (default "changeme")
  -U, --unsafe            Ignore certificate errors
  -u, --username string   Username for elasticsearch (default "elastic")

Use "esinfo [command] --help" for more information about a command.
```
