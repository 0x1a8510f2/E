#!/bin/sh

if [[ -z "$GID" ]]; then
	GID="$UID"
fi

# Define functions.
function fixperms {
	chown -R $UID:$GID /data /opt/E
}

# If the config was not found
if [[ ! -f /data/config.yaml ]]; then
	echo "[!] Could not find config.yaml"
	cp /opt/E/example-config.yaml /data/config.yaml
	echo "[#] Copied default config file to /data/config.yaml"
	echo "[#] Please edit /data/config.yaml and re-start the container. This will generate a registration file to be placed on the homeserver."
	exit
fi

#if [[ ! -f /data/registration.yaml ]]; then
#	/usr/bin/E -g -c /data/config.yaml -r /data/registration.yaml
#	echo "[#] Could not find a registration file so it has been generated."
#	echo "[#] Copy /data/registration.yaml to your homeserver."
#	exit
#fi

cd /data
fixperms
exec su-exec $UID:$GID /usr/bin/E
