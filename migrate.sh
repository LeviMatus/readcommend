#!/bin/bash
# Creates the database tables and seeds them with sample data
. vars.sh
psql -a -q -f migrate.sql $DB_URL
