#!/bin/sh

# getting env settings from .env
set -a
. /mnt/UserData/ic-pcp-wifi-plus-server/.env
set +a

echo "------------------------------------------------------------------------"
echo "--------------------------- wifi-plus ----------------------------------"
echo "------------------------------------------------------------------------"
echo "Starting wifi-plus-startup script..."
echo "Copying go binary and script files to web folders..."
sudo chmod 777 /mnt/UserData/ic-pcp-wifi-plus-server/pcp-scripts/wifi-plus.sh
if sudo cp /mnt/UserData/ic-pcp-wifi-plus-server/pcp-scripts/wifi-plus.sh /var/www/cgi-bin/wifi-plus.sh; then
  sudo chown tc:staff /var/www/cgi-bin/wifi-plus.sh
  echo "Successfully copied wifi-plus shell script to cgi-bin"
else
  echo "Unable to copy shell file to cgi-bin"
  echo "Exiting..."
  exit 1
fi

echo "Checking for any running 'wifiplus' processes..."
wifiplus_pid=0
wifiplus_pid=$(pidof wifiplus)
if [ "$wifiplus_pid" ] && [ "$wifiplus_pid" -ne 0 ]; then
    echo "'wifiplus' process [$wifiplus_pid] found!"
  if sudo kill -9 "$wifiplus_pid"; then
    echo "Process [$wifiplus_pid] terminated"
  else
    echo "'wifiplus' process could not be terminated!"
    echo "Exiting..."
    exit 1
  fi
else
  echo "No 'wifiplus' process found"
fi

echo "Copying wifiplus hostapd.conf to /usr/local/etc/pcp/"
cp /mnt/UserData/ic-pcp-wifi-plus-server/confs/hostapd.conf /usr/local/etc/pcp/hostapd.conf
if [ -f /usr/local/etc/pcp/dnsmasq.conf ]; then
  sudo chown tc:staff /usr/local/etc/pcp/dnsmasq.conf
fi
if [ -f /usr/local/etc/pcp/hostapd.conf ]; then
  echo "Copying wifiplus pcp_hosts to /usr/local/etc/pcp/pcp_hosts"
  cp /mnt/UserData/ic-pcp-wifi-plus-server/confs/pcp_hosts /usr/local/etc/pcp/pcp_hosts
  sudo chown root:root /usr/local/etc/pcp/pcp_hosts
  sudo chmod 644 /usr/local/etc/pcp/pcp_hosts
  if [ $(pidof dnsmasq) ]; then
    sudo kill -9 $(pidof dnsmasq)
    sleep 5
  fi

  sudo dnsmasq -C /usr/local/etc/pcp/dnsmasq.conf
  sudo chown tc:staff /usr/local/etc/pcp/hostapd.conf
  sudo chmod 644 /usr/local/etc/pcp/hostapd.conf
fi

if sudo cp /mnt/UserData/ic-pcp-wifi-plus-server/wifiplus /var/www/wifiplus &&
   sudo cp /mnt/UserData/ic-pcp-wifi-plus-server/.env /var/www/.env; then

  echo "Successfully copied wifi-plus binary to webroot"
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
  wifiplus_pid=0
  wifiplus_pid=$(pidof wifiplus)
  if [ "$wifiplus_pid" ] && [ "$wifiplus_pid" -ne 0 ]; then
    printf "Binary started successfully.\nProcess [$wifiplus_pid] listening in background on port [$PORT]...\n"
    echo "Testing connection..."
    sleep 1
    rc=$(curl -s -o /dev/null -w "%{http_code}" http://"$ICHOST""$PORT"/system/status)

    if [ ! "$rc" = "200" ]; then
      rc=$(curl -s -o /dev/null -w "%{http_code}" http://10.10.10.1/system/status)
    fi

    if [ "$rc" = "200" ]; then
      echo "[$rc OK] API up and running :)"
      exit 0
    else
      echo "Unable to connect to API successfully"
      echo "Status code is [$rc]"
      echo "Exiting..."
    fi
  else
    echo "No PID found for binary. Looks like it failed to start :("
    echo "Exiting..."
  fi

else
  echo "Unable to copy binary to webroot"
  echo "Exiting..."
fi

exit 1