# Tollgate Module - whoami (go)

This Tollgate module will give back a user's MAC address based on the IP address that's been assigned to them. This is necesarry as a workaround for restrictions by the Android operating sytstem, which does not allow for non-system apps to get the device's MAC address.

# Compile for ATH79 (GL-AR300 NOR)

```bash
cd ./src
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o whoami -trimpath -ldflags="-s -w"

# Hint: copy to connected router 
scp -O whoami root@192.168.8.1:/tmp/whoami
```

# Compile for GL-MT3000

## Build

```bash
cd ./src
env GOOS=linux GOARCH=arm64 go build -o whoami -trimpath -ldflags="-s -w"

# Hint: copy to connected router 
scp -O whoami root@192.168.1.1:/root/whoami # X.X == Router IP
```

## Required Firewall rules 

First, test if the whoami is up by going to your router's ip on port `2122`. You should get a JSON response with your IP and mac address.

Add to `/etc/config/firewall`:
```uci
config rule
	option name 'Allow-whoami-In'
	option src 'lan'
	option proto 'tcp'
	option dest_port '2122' # whoami port
	option target 'ACCEPT'

config redirect
	option name 'TollGate - Nostr whoami DNAT'
	option src 'lan'
	option dest 'lan'
	option proto 'tcp'
	option src_dip '192.168.21.21'
	option src_dport '2121'
	option dest_ip '192.168.X.X' # Router IP
	option dest_port '2122' # whoami port
	option target 'DNAT'

config redirect
        option name 'TollGate - Nostr whoami DNAT port'
        option src 'lan'
        option dest 'lan'
        option proto 'tcp'
        option src_dip '192.168.X.X' # Router IP
        option src_dport '2121'
        option dest_ip '192.168.X.X' # Router IP
        option dest_port '2122' # whoami port
        option target 'DNAT'
```

Run `service firewall restart` to make changes go into effect.

To test the firewall rule, go to `192.168.21.21:2122`. You should be greeted with the same JSON.


## OpenNDS rules
**Prerequisite: OpenNDS is installed**

To allow unauthenticated clients to reach the whoami, we need to explicitly allow access.

Add to `/etc/config/opennds` under `config opennds`:
```uci
config opennds
    list users_to_router 'allow tcp port 2122' # whoami port
    list preauthenticated_users 'allow tcp port 2122 to 192.168.21.21'
```

Run `service opennds restart` to make changes go into effect.

## License
This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

