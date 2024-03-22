#!/bin/sh

# Function to check if a service is healthy
check_service_healthy() {
  name="$1"
  port="$2"
  max_attempts="$3"
  interval="$4"
  code=1

  for i in $(seq 1 "$max_attempts"); do
    nc -z "$name" "$port" >/dev/null 2>&1
    code=$?
    if [ $code -eq 0 ]; then
      echo "$name is healthy"
      return 0
    fi
    sleep "$interval"
  done

  echo "$name is unhealthy"
  return 1
}

# Usage: wait_for_service <name> <port> [<max_attempts>] [<interval>]
wait_for_service() {
  name="$1"
  port="$2"
  max_attempts="${3:-100}"
  interval="${4:-2}"

  echo "Waiting for $name to be healthy..."

  if check_service_healthy "$name" "$port" "$max_attempts" "$interval"; then
    echo "$name is healthy"
    return 0
  else
    echo "$name is unhealthy"
    return 1
  fi
}

# Example usage: wait_for_service db 5432 100 2
wait_for_service db 5432 10 2