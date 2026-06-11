#!/bin/bash
set -e

CLEAN=false
if [ "$1" = "clean" ]; then
    CLEAN=true
fi

# Always kill previous demo processes
pkill -9 demo 2>/dev/null || true
sleep 1

# Vendor refresh (only on clean, otherwise reuse existing)
if $CLEAN || [ ! -d vendor ]; then
    rm -rf go.sum go.mod vendor
    go mod init
    # Pin l8types to avoid breaking Register() change (l8ql not yet updated)
    GOPROXY=direct GOPRIVATE=github.com go get github.com/saichler/l8types@v0.0.0-20260502164503-192680a11be2 2>/dev/null
    GOPROXY=direct GOPRIVATE=github.com go mod tidy
    go mod vendor

    # Apply l8alarms nil-guard fix to vendor (until upstream is merged)
    L8ALM_SRC="../l8alarms/go/alm"
    L8ALM_VENDOR="vendor/github.com/saichler/l8alarms/go/alm"
    for f in maintenancewindows/checker.go notification/engine.go escalation/scheduler.go; do
        if [ -f "$L8ALM_SRC/$f" ] && [ -f "$L8ALM_VENDOR/$f" ]; then
            cp "$L8ALM_SRC/$f" "$L8ALM_VENDOR/$f"
        fi
    done
fi

if $CLEAN; then
    # Recreate database
    echo "Clean mode: recreating database..."
    docker rm -f unsecure-postgres 2>/dev/null || true
    docker run -d --name unsecure-postgres -p 5432:5432 \
        -v /data/:/data/ saichler/unsecure-postgres:latest \
        /bin/sh -c "/start-postgres.sh admin admin admin 5432 && tail -f /dev/null"
    sleep 5
else
    # Ensure database is running
    if ! docker ps --format '{{.Names}}' | grep -q unsecure-postgres; then
        echo "Database not running, starting..."
        docker rm -f unsecure-postgres 2>/dev/null || true
        docker run -d --name unsecure-postgres -p 5432:5432 \
            -v /data/:/data/ saichler/unsecure-postgres:latest \
            /bin/sh -c "/start-postgres.sh admin admin admin 5432 && tail -f /dev/null"
        sleep 5
    fi
fi

# Build binaries
rm -rf demo && mkdir -p demo
cd tests/mocks/cmd && go build -o ../../../demo/mocks_demo && cd ../../../
cd vend/vnet && go build -o ../../demo/vnet_demo && cd ../../
cd vend/main && go build -o ../../demo/vend_demo && cd ../../
cd vend/ui/main && go build -o ../../../demo/ui_demo && cd ../../../
cd vend/collector && go build -o ../../demo/collector_demo && cd ../../
cd vend/parser && go build -o ../../demo/parser_demo && cd ../../
cd vend/inv_vend && go build -o ../../demo/inv_demo && cd ../../
cp -r vend/ui/web demo/.

# Generate kill script
cd demo
cat > kill_demo.sh <<'EOF'
cd ..
rm -rf demo
rm -rf /data/postgres/admin
pkill -9 demo
EOF
chmod +x kill_demo.sh

LOGFILE="../demo/demo.log"
> $LOGFILE

# Start services — UI must start BEFORE backend so it receives web service broadcasts
echo "Starting VNet..."
./vnet_demo >> $LOGFILE 2>&1 &
sleep 2

echo "Starting UI..."
./ui_demo >> $LOGFILE 2>&1 &
sleep 2

echo "Starting collector..."
./collector_demo >> $LOGFILE 2>&1 &
sleep 2

echo "Starting parser..."
./parser_demo >> $LOGFILE 2>&1 &
sleep 2

echo "Starting inventory cache..."
./inv_demo >> $LOGFILE 2>&1 &
sleep 2

echo "Starting backend (local mode)..."
./vend_demo local >> $LOGFILE 2>&1 &
sleep 5

echo "All services started! Logs: demo/demo.log"

if $CLEAN; then
    # Upload mock data
    EXTERNAL_IP=$(ip route get 1.1.1.1 | grep -oP 'src \K[0-9.]+')
    read -p "Press Enter to upload mocks"
    ./mocks_demo --address https://${EXTERNAL_IP}:4443 --user admin --password admin --insecure --simulator 192.168.200.1 --simulator-port 8443 2>&1 | tee -a $LOGFILE
fi

read -p "Press Enter to kill the demo"
#./kill_demo.sh
