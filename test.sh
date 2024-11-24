#!/bin/bash

POST_URL="http://localhost:4000/api/receipts/process"
GET_URL="http://localhost:4000/api/receipts"

TEST_PAYLOADS="payload.json"

# Read and parse each JSON objects
jq -c '.[]' "$TEST_PAYLOADS" | while read -r PAYLOAD; do
  echo "Post Process Receipt"
  echo "$PAYLOAD"

  # Send POST request
  POST_RESPONSE=$(curl -s -X POST "$POST_URL" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD")

  echo "POST Response: $POST_RESPONSE"

  # Get receipt ID from response
  ID=$(echo "$POST_RESPONSE" | sed -n 's/.*"id"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')

  echo "Receipt ID: $ID"

  if [ -z "$ID" ] || [ "$ID" == "null" ]; then
    echo "Failed to extract 'id' from the POST response."
    continue
  fi

  echo "Receipt ID: $ID"

  # Send GET request using the Id
  GET_RESPONSE=$(curl -s -X GET "$GET_URL/$ID/points")

  echo "$GET_RESPONSE"
  echo "--------------------------"
done
