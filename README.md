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
  -agent-version string
        Specify the agent version.
  -cluster string
        Set a AutoScaling Group Name.
  -drain
        Execute draining.
  -drain-all
        Execute all instance draining.
  -instance string
        Specify the instances.
  -type string
        Specify the launch type (EC2 or Fargate)
  -version
        Print version number.

```

## Usage

### List up ECS cluster

```sh
$ ecm
```

You can use the `type` option to filter the Lauch Type.

```sh
$ ecm -type=[EC2|FARGATE]
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

The `drain-all` option and the latest Agent Version (`agent-version`) can be used to drain all instances of an older Agent Version in a cluster.

```sh
$ ecm -cluster=default -drain-all -agent-version=[New Agent Version]
```
