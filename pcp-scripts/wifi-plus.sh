#!/bin/sh

# --------------------------------------------------------------- #
# wifi-plus.sh - called by the wifi-plus binary to interface with #
#                picoreplayer subroutines.                        #
# --------------------------------------------------------------- #

. pcp-functions
. pcp-wifi-functions

subroutine=$1
arg1=$2
arg2=$3
arg3=$4
arg4=$5

wp_status() {
  echo $arg1
}

wp_test() {
  pcp_set_coloured_text
  echo "Able to call pcp functions"
}

# ---------------------- main program ---------------------- #

case $subroutine in
  wp_status)
    wp_status
  ;;
  wp_test)
    wp_test
  ;;
  *)
    echo "$subroutine"
  ;;
esac