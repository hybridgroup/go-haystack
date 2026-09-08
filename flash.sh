# This script will flash the target device with the keys for the keys file
#
# For a device on a battery, add the low power settings, for example:
#   EXTRA_FLAGS="-serial=none" EXTRA_LDFLAGS="-X main.TxPower=-8" ./flash.sh TARGET KEYSFILE
TARGET=$1
KEYSFILE=$2
echo "Flashing $TARGET device with keys $KEYSFILE"

ADVKEY=$(awk -F: 'NR==2 {gsub(/^ +/, "", $2); print $2}' ${KEYSFILE})
cd ./firmware
tinygo flash -target $TARGET $EXTRA_FLAGS -ldflags="-X main.AdvertisingKey='$ADVKEY' $EXTRA_LDFLAGS" .
