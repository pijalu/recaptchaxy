#!/bin/bash

export RC_LISTEN=":8080"
export RC_TARGET="https://myendpoint"
export RC_PROJECT_ID="<projectid>"
export RC_APIKEY="<apikey>"

go run main.go