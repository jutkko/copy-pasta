#!/bin/bash -e

# set env var DEBUG to 1 to enable debugging
[ -z "$DEBUG" ] || set -x

# print out usage when something goes wrong
usage() {
  printf "usage: %s <target>\n" "$0"
  exit 1
}

main() {
  if [ "$#" -lt 1 ]
  then
    usage "$@"
  fi

  if ! lpass status > /dev/null
  then
    echo "login to lpass"
    exit 1
  fi

  local accesskey=<get your key>
  local secretaccesskey=<get your secret key>

  local target="$1"
  printf "Setting a new target: %s\n" "$target"
  printf "%s\n%s\n" "$accesskey" "$secretaccesskey" | copy-pasta login --target "$target"
}

main $@
