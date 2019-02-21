# /bin/sh

export FLOTILLA_MODE=dev
export HTTP_SERVER_CORS_ALLOWED_ORIGINS=http://localhost:5000

# these 2 envvars need to be present
# AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}
# AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY}

go install github.com/datagovsg/flotilla-os && ~/go/bin/flotilla-os conf
