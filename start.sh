#!/bin/bash
export ADDRESS=localhost
export PORT=8080
export client_secret=oGDXMy5jsD4qEEjfWr1AOW1fuzKKZvQ9y1RcmhbYPNiesfi0haNYO7e9sL6FtlGr
export client_id=IXhXcAAfhY4WUIck8Sp1CgavF1zMVjVoihMoIvSuePQN6dpr
export grant_type=client_credentials
export scope=openapi
export AUTH_URL=https://test.optii.io/oauth/authorize
export OPTII_URL=https://test.optii.io/api/v1
go run main.go
