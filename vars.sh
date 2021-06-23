#!/bin/bash
# Defines env vars required by other scripts (this script must be sourced using "source" or "." in
# order for exported variables to be accessible outside itself).
export DB_URL="postgres://postgres:password123@localhost:5432/readcommend?sslmode=disable"
