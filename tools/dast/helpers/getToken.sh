#!/bin/bash

# Get auth token, using auth_post_request.json POST data.
# We use this approach to prevent secrets being uploaded to git.

APIRESPONSE=$(curl --silent -X POST -H "Content-Type: application/json" \
	-d @/zap/helpers/auth_post_request.json \
	https://<insert your URL>)

# Remember to change this regex, so it matches your response
echo $APIRESPONSE | sed 's/^{.*"access_token":"\(.*\)",.*$/\1/' 
