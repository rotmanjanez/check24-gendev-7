# How to Contribute



## Client

## Server

### Extending the API
When extending the API, adjust the `server/api/openapi.yml` first with all required types.
Then use the `generate_api.sh` script in the server directory to generate the scaffolding.
There are currently only three manual adjustemts required:
1. Add your service to the router in `server/cmd/check24-gendev-7-server/main.go`
2. Change the package of any new model files in `pkg/models`from `api`to `models`
3. Add any new service files to the `.openapi-generator-ignore` file, so they dont get overritten in the future

> **Throubleshooting Tip**: Encountering missing definition?
> To have better separation, the generated `errors.go` and `helpers.go` are split into the `models` package and the `api` package. Try removing those files from the `.openapi-generator-ignore` file and check via git if any new definitions got added, that are in neither of the package files.


### Making a API request locally
You can use curl to test the api locally without spinning up the client dev environment.
```sh
curl -X GET "http://localhost:8080/internet-products" -H "Content-Type: application/json" -d '{"street": "Erika-Mann-Straße", "houseNumber": 62, "city": "München", "postalCode": "80636", "countryCode": "DE" }'
```