#!/bin/bash

go tool oapi-codegen -config client-config.yaml ../contracts/analysis_service/api.yaml

go tool oapi-codegen -config server-config.yaml ../contracts/analysis_service/api.yaml
