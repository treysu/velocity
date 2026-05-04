#!/bin/sh

PUID=${PUID:-1000}
PGID=${PGID:-1000}
VELOCITY_VERSION=${VELOCITY_VERSION:-"latest"}

echo "Starting Velocity with UID: $PUID, GID: $PGID"

VELOCITY_JAR="/opt/velocity/velocity.jar"
mkdir -p /opt/velocity

if [ "$VELOCITY_VERSION" = "latest" ]; then
    echo "Fetching latest Velocity version..."
    VELOCITY_VERSION=$(wget -qO- "https://api.papermc.io/v2/projects/velocity" | \
        grep -o '"versions":\["[^"]*"' | grep -o '[0-9][^"]*' | tail -1)
    echo "Latest version: $VELOCITY_VERSION"
fi

LATEST_BUILD=$(wget -qO- "https://api.papermc.io/v2/projects/velocity/versions/${VELOCITY_VERSION}" | \
    grep -o '"builds":\[[^]]*\]' | grep -o '[0-9]*' | tail -1)
echo "Latest build: $LATEST_BUILD"

DOWNLOAD_URL="https://api.papermc.io/v2/projects/velocity/versions/${VELOCITY_VERSION}/builds/${LATEST_BUILD}/downloads/velocity-${VELOCITY_VERSION}-${LATEST_BUILD}.jar"

echo "Downloading Velocity from $DOWNLOAD_URL..."
wget -O "$VELOCITY_JAR" "$DOWNLOAD_URL"

if getent group velocity > /dev/null 2>&1; then
    groupmod -g $PGID velocity
else
    addgroup -g $PGID velocity
fi

if getent passwd velocity > /dev/null 2>&1; then
    usermod -u $PUID -g $PGID velocity
else
    adduser -u $PUID -S velocity -G velocity
fi

chown -R velocity:velocity /data
chown -R velocity:velocity /opt/velocity

exec gosu velocity java \
    -Xms$JAVA_MEMORY \
    -Xmx$JAVA_MEMORY \
    $JAVA_FLAGS \
    -jar "$VELOCITY_JAR"