#!/bin/bash

# Add a file, to use global variables 
source /zap/helpers/.env

# Ensure report directory exists
mkdir /zap/report/

echo "Zap .env file:"
cat /zap/helpers/.env

echo "Zap headers file:"
cat /zap/helpers/auth_post_request.json

# Make a request to the API, to get authorization (The token)
# The "$@" parameter adds the ability to add headers from within this script, 
# So there is no need for script modification.
/zap/helpers/genJSON.sh Authorization "$ACCESS_TOKEN" "$@" > /zap/helpers/headers.json

echo "Zap configs:"
cat /zap/helpers/zapConfigs

# First execution of zap, to install addons and configure ZAP to run.
# add-headers.py script on every request
zap.sh -cmd -addoninstall jython -configfile /zap/helpers/zapConfigs
# Second execution of zap, to run actual tests against specified endpoints.
zap.sh -cmd -autorun /zap/"$ZAP_CONFIG_FILE"
