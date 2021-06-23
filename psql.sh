#!/bin/bash
# Connects to database for development purposes
. vars.sh
psql $DB_URL
