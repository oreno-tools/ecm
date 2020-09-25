# ecm

## これは

* ECS のコンテナインスタンスを管理するツールです

## Install


```sh
# Get latest version
v=$(curl -s 'https://api.github.com/repos/oreno-tools/ecm/releases' | jq -r '.[0].tag_name')
# For macOS
$ wget https://github.com/oreno-tools/ecm/releases/download/${v}/ecm_darwin_amd64 -O ~/bin/ecm && chmod +x ~/bin/ecm
# For Linux
$ wget https://github.com/oreno-tools/ecm/releases/download/${v}/ecm_linux_amd64 -O ~/bin/ecm && chmod +x ~/bin/ecm
```

## Help

```sh
$ ecm --help
Usage of ecm:
  -cluster string
        Set a AutoScaling Group Name.
  -drain
        Execute draining.
  -instance string
        Specify the instances.
  -version
        Print version number.
```

## Usage

### List up AutoScaling Group

```sh
$ ecm
+--------------------------------------------------------------------------+
|                               CLUSTER NAME                               |
+--------------------------------------------------------------------------+
| default                                                                  |
+--------------------------------------------------------------------------+
```

### List up Container Instances

```sh
$ ecm -cluster=default
+--------------------------------------+---------------------+---------------+----------------+---------------+--------+
|          CONTAINER INSTNACE          |     INSTANCE ID     | AGENT VERSION | DOCKER VERSION | RUNNING TASKS | STATUS |
+--------------------------------------+---------------------+---------------+----------------+---------------+--------+
| xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx | i-1234567890abcdefg | 1.44.4        | 19.03.6-ce     |             2 | ACTIVE |
+--------------------------------------+---------------------+---------------+----------------+---------------+--------+
```

### Update Container Instance Status (Active to Draining)

```sh
$ ecm -cluster=default -drain -instance=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
