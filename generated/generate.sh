#!/bin/bash

oapi-codegen -config client-config.yaml ../contracts/analysis_service/api.yaml

oapi-codegen -config server-config.yaml ../contracts/analysis_service/api.yaml
