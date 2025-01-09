# This script will flash the target device with the keys for the keys file
TARGET=$1
KEYSFILE=$2
echo "Flashing $TARGET device with keys $KEYSFILE"

ADVKEY=$(awk -F: 'NR==2 {gsub(/^ +/, "", $2); print $2}' ${KEYSFILE})
cd ./firmware
tinygo flash -target $TARGET -ldflags="-X main.AdvertisingKey='$ADVKEY'" .
