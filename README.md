# snmp-simulator

## Commands

```sh
snmpget -v2c -c public 127.0.0.1:1610 1.3.6.1.2.1.1.3.0

snmpget -v3 -l authPriv -u snmpuser -a SHA -A "authpass" -x AES -X "privpass" 127.0.0.1:1610 1.3.6.1.2.1.1.3.0

snmpwalk -v2c -c public 127.0.0.1:1610 1.3.6.1.2.1.2.2
```

## Docker

```sh
docker build . -t snmpsimulator:local
docker run -e CONFIG_FILEPATH=config.yaml -v $(pwd)/config.yaml:/config.yaml snmpsimulator:local
```

## TODO

- [ ] golangci lint
- [ ] create/publish helm chart
- [ ] unit tests
- [ ] release process
- [ ] Contributing guide

