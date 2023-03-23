#!/usr/bin/env bash

set -Eeuo pipefail
set -o xtrace

API_URL=http://bridge.zalizniak.duckdns.org
# API_URL=http://localhost:8080

user1_email=test1@bridge.test
user1_password=Password123

response=$(curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"email\":\"$user1_email\", \"password\":\"$user1_password\"}" \
  $API_URL/auth/login)
user_id1=$( jq -r  '.user.id' <<< "${response}" ) 
token1=$( jq -r  '.user.token' <<< "${response}" ) 
success=$( jq -r  '.success' <<< "${response}" ) 
if [[ "$success" == "false" ]]; then
  exit 1
fi

user2_email=zalizniak@zalizniak
user2_password=coco

response=$(curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"email\":\"$user2_email\", \"password\":\"$user2_password\"}" \
  $API_URL/auth/login)
user_id2=$( jq -r  '.user.id' <<< "${response}" ) 
token2=$( jq -r  '.user.token' <<< "${response}" ) 
success=$( jq -r  '.success' <<< "${response}" ) 
if [[ "$success" == "false" ]]; then
  exit 1
fi

response=$(curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"host_id\":\"$user_id1\", \"user_id\":\"$user_id1\", \"token\":\"$token1\"}" \
  $API_URL/room/create)
room_id=$( jq -r  '.room_id' <<< "${response}" ) 
success=$( jq -r  '.success' <<< "${response}" ) 
if [[ "$success" == "false" ]]; then
  exit 1
fi

response=$(curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"room_id\":\"$room_id\", \"user_id\":\"$user_id2\", \"token\":\"$token2\"}" \
  $API_URL/room/join )
success=$( jq -r  '.success' <<< "${response}" ) 
if [[ "$success" == "false" ]]; then
  exit 1
fi
echo $response |  jq | jq

response=$(curl --request GET \
  "$API_URL/room/$room_id?token=$token1&user_id=$user_id1"  | jq)
success=$( jq -r  '.success' <<< "${response}" ) 
if [[ "$success" == "false" ]]; then
  exit 1
fi
echo $response |  jq

response=$(curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"room_id\":\"$room_id\", \"user_id\":\"$user_id1\", \"token\":\"$token1\"}" \
  $API_URL/session/create)
session_id=$( jq -r  '.session_id' <<< "${response}" ) 
success=$( jq -r  '.success' <<< "${response}" ) 
if [[ "$success" == "false" ]]; then
  exit 1
fi

response=$(curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"user_id\":\"$user_id2\", \"token\":\"$token2\"}" \
  $API_URL/session/getByUser)
user2_session_id=$( jq -r  '.session_id' <<< "${response}" ) 
success=$( jq -r  '.success' <<< "${response}" ) 
if [[ "$success" == "false" ]]; then
  exit 1
fi

if [[ "$user2_session_id" != "$session_id" ]]; then
  exit 1
fi


response=$(curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"session_id\":\"$session_id\", \"user_id\":\"$user_id1\", \"token\":\"$token1\"}" \
  $API_URL/session/close)
if [[ "$success" == "false" ]]; then
  exit 1
fi
echo $response |  jq

response=$(curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"room_id\":\"$room_id\", \"user_id\":\"$user_id1\", \"token\":\"$token1\"}" \
  $API_URL/room/delete)
if [[ "$success" == "false" ]]; then
  exit 1
fi
echo $response |  jq
