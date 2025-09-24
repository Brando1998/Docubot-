#!/bin/bash
# Simple health check for Rasa
curl -f http://localhost:5005 || exit 1