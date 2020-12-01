wgman is a simple tool to manage wireguard configurations with wireguard binary "wg". It saves peers in config file ("data.json" by default) and generate wg-quick compitable config files.

Run the binary to see available commands.

Use environmental variable `WGMAN_CONFIG` to specify a different config file.

# Examples

```sh
# creates a new network in data.json
wgman new 10.1.1.0
# add your nb
wgman add me
# add your gateway server, nb can connect to it directly
wgman add gw gw.my.com:1234
# add db server, relayed through gateway
wgman add-inner gw db 192.168.1.5
# list all peers
wgman list
# creates a new network in test.json
WGMAN_CONFIG=test.json wgman new 10.1.1.0
```

# CAUTION

`wgman new` *WILL NOT* prompt you if `data.json` exists.

# License

MPL 2.0
