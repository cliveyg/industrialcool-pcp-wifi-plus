#!/bin/sh

# getting env files from .env
set -a
source .env
set +a

echo "-----------------------------------------------------------------------"
echo "Starting wifi-plus-startup..."
echo "Copying go binary and script files to web folders..."
sudo chmod 777 /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/wifi-plus.sh
if sudo cp /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/wifi-plus.sh /var/www/cgi-bin/wifi-plus.sh; then
  echo "Successfully copied wifi-plus shell script to cgi-bin"
else
  echo "Unable to copy shell file to cgi-bin."
  echo "Exiting..."
  exit 1
fi

if sudo cp /mnt/UserData/industrialcool-pcp-wifi-plus/wifiplus /var/www/wifiplus &&
   sudo cp /mnt/UserData/industrialcool-pcp-wifi-plus/.env /var/www/.env; then

  echo "Successfully copied wifi-plus binary to webroot"
  echo "Checking for any running 'wifiplus' processes..."
  wifiplus_pid=0
  wifiplus_pid=$(pidof wifiplus)
  if [ "$wifiplus_pid" ] && [ "$wifiplus_pid" -ne 0 ]; then
    if kill -9 "$wifiplus_pid"; then
      echo "Running 'wifiplus' process [\"$wifiplus_pid\"] terminated"
    else
      echo "'wifiplus' process could not be terminated!"
      echo "Exiting..."
      exit 1
    fi
  else
    echo "No 'wifiplus' process found."
  fi

  echo "Creating logfile in [$LOGFILE]..."
  if [ -f "$LOGFILE" ]; then
    echo "Logfile already exists"
  else
    if sudo touch "$LOGFILE" && sudo chmod 666 "$LOGFILE"; then
      echo "Logfile created"
    else
      echo "Unable to create logfile. Exiting..."
      exit 1
    fi
  fi

  echo "Attempting to start binary..."

  cd /var/www/ && nohup ./wifiplus > /dev/null 2>&1 &

  printf "Binary started successfully.\nListening in background on port 8020...\n"
  echo "Testing connection..."
  sleep 3
  rc=$(curl -s -o /dev/null -w "%{http_code}" http://"$ICHOST":8020/status)
  if [ "$rc" = "200" ]; then
    echo "[$rc OK] API up and running."
    exit 0
  else
    echo "Unable to connect to API successfully."
    echo "Status code is [$rc]"
    echo "Exiting..."
  fi

else
  echo "Unable to copy binary to webroot."
  echo "Exiting..."
fi

exit 1