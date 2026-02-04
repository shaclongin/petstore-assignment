#!/bin/bash

set -euo pipefail

base=${1:-http://localhost:9900/api}

create_pet() {
    pet_name=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 13; echo '')
    pet_id=$(echo $(( $RANDOM % 200 + 1 )))
    echo "Creating pet ${pet_id}"
    curl 2>/dev/null --location --request POST "${base}/pet" \
    --header 'Accept: application/json' \
    --header 'Content-Type: application/json' \
    --data-raw "{
    \"name\": \"${pet_name}\",
    \"photoUrls\": [
        \"velit mollit dolore\",
        \"sed\"
    ],
    \"id\": ${pet_id},
    \"tags\": [
        {
            \"id\": -15522919,
            \"name\": \"velit ut in esse aliquip\"
        },
        {
            \"id\": 36974757,
            \"name\": \"ullamco mollit sed commodo\"
        }
    ],
    \"status\": \"available\"
}"
}

get_pet() {
    pet_id=$(echo $(( $RANDOM % 200 + 1 )))
    echo "Getting pet ${pet_id}"
    curl 2>/dev/null --location --request GET "${base}/pet/${pet_id}" \
--header 'Accept: application/json'
}

while true
do
    (( RANDOM % 2 )) || (get_pet &)
    (( RANDOM % 10 )) || (create_pet &)
    sleep 0.1
done