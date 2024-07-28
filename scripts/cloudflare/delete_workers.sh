#!/bin/bash

# Cloudflare API details
CLOUDFLARE_ACCOUNT_ID="efb0ce8ef67369c42f4056264dac6f8c"
CLOUDFLARE_API_TOKEN="ngA1ttp8_LoIz-YGF20262sM9HGK6Nd--d47bsIU"

# List and count Cloudflare Workers
response=$(curl -X GET "https://api.cloudflare.com/client/v4/accounts/efb0ce8ef67369c42f4056264dac6f8c/workers/scripts" \
     -H "Authorization: Bearer ngA1ttp8_LoIz-YGF20262sM9HGK6Nd--d47bsIU" \
     -H "Content-Type: application/json")    

worker=$($response | jq -r '.result[].id')
worker_count=$(echo $response | jq -r '.result | length')

for worker in $workers; do
  echo "Deleting Worker: $worker"
  del_response=$(curl -s -X DELETE "https://api.cloudflare.com/client/v4/accounts/efb0ce8ef67369c42f4056264dac6f8c/workers/scripts/$worker" \
   -H "Authorization: Bearer ngA1ttp8_LoIz-YGF20262sM9HGK6Nd--d47bsIU" \
   -H "Content-Type: application/json")
   echo "$response"
done

echo "Total number of Workers deleted: $worker_count"


