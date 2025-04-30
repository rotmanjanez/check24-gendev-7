openapi-generator generate \
    -g typescript \
    --additional-properties=platform=browser,framework=fetch-api \
    -i ../openapi.yaml \
    -o src/api 