#!/bin/bash

# inject the env name into datadog config
sed -i 's@__APM_ENV__@'"${GO_ENV}"'@g' /etc/datadog-agent/datadog.yaml