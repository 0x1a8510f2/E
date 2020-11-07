#!/bin/sh

if [[ -z "$GID" ]]; then
	GID="$UID"
fi

# Define functions.
function fixperms {
	chown -R $UID:$GID /data /opt/E
}

# Temporarily commented out for development
#if [[ ! -f /data/config.yaml ]]; then
#	cp /opt/E/example-config.yaml /data/config.yaml
#	echo "Could not find config.yaml"
#	echo "Copied default config file to /data/config.yaml"
#	echo "Please edit /data/config.yaml and re-start the container. This will generate a registration file."
#	exit
#fi

#if [[ ! -f /data/registration.yaml ]]; then
#	/usr/bin/E -g -c /data/config.yaml -r /data/registration.yaml
#	echo "Could not find a regastration file so it has been generated."
#	echo "Copy /data/registration.yaml to your homeserver."
#	exit
#fi

cd /data
fixperms
exec su-exec $UID:$GID /usr/bin/E
