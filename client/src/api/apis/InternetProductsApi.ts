// TODO: better import syntax?
import { BaseAPIRequestFactory, RequiredError } from './baseapi';
import type { Configuration } from '../configuration';
import { RequestContext, HttpMethod, ResponseContext, HttpInfo } from '../http/http';
import { ObjectSerializer } from '../models/ObjectSerializer';
import { ApiException } from './exception';
import { isCodeInRange } from '../util';
import type { SecurityAuthentication } from '../auth/auth';


import { Address } from '../models/Address';
import { InternetProductsCursor } from '../models/InternetProductsCursor';
import { InternetProductsResponse } from '../models/InternetProductsResponse';
import { SharedInternetProductsResponse } from '../models/SharedInternetProductsResponse';

/**
 * no description
 */
export class InternetProductsApiRequestFactory extends BaseAPIRequestFactory {

    /**
     * Fetches the next batch of internet products using a continuation cursor
     * @param cursor Cursor to continue fetching products
     */
    public async continueInternetProductsQuery(cursor: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'cursor' is not null or undefined
        if (cursor === null || cursor === undefined) {
            throw new RequiredError("InternetProductsApi", "continueInternetProductsQuery", "cursor");
        }


        // Path Params
        const localVarPath = '/internet-products/continue';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (cursor !== undefined) {
            requestContext.setQueryParam("cursor", ObjectSerializer.serialize(cursor, "string", ""));
        }



        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Retrieves the shared internet products using a given cursor
     * @param cursor Cursor to retrieve the shared products
     */
    public async getSharedInternetProducts(cursor: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'cursor' is not null or undefined
        if (cursor === null || cursor === undefined) {
            throw new RequiredError("InternetProductsApi", "getSharedInternetProducts", "cursor");
        }


        // Path Params
        const localVarPath = '/internet-products/share/{cursor}'
            .replace('{' + 'cursor' + '}', encodeURIComponent(String(cursor)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")



        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Initiates retrieval of internet products and returns a product version and a cursor to retrieve the first batch of products
     * @param address 
     * @param providers Providers to filter the products
     */
    public async initiateInternetProductsQuery(address: Address, providers?: Array<string>, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'address' is not null or undefined
        if (address === null || address === undefined) {
            throw new RequiredError("InternetProductsApi", "initiateInternetProductsQuery", "address");
        }



        // Path Params
        const localVarPath = '/internet-products';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (providers !== undefined) {
            const serializedParams = ObjectSerializer.serialize(providers, "Array<string>", "");
            for (const serializedParam of serializedParams) {
                requestContext.appendQueryParam("providers", serializedParam);
            }
        }


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(address, "Address", ""),
            contentType
        );
        requestContext.setBody(serializedBody);


        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Shares the internet products with a given cursor. This cursor must be the same as the one returned by the initial query and the query must have completed.
     * @param cursor Cursor to share the products
     */
    public async shareInternetProducts(cursor: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'cursor' is not null or undefined
        if (cursor === null || cursor === undefined) {
            throw new RequiredError("InternetProductsApi", "shareInternetProducts", "cursor");
        }


        // Path Params
        const localVarPath = '/internet-products/share/{cursor}'
            .replace('{' + 'cursor' + '}', encodeURIComponent(String(cursor)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")



        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

}

export class InternetProductsApiResponseProcessor {

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to continueInternetProductsQuery
     * @throws ApiException if the response code was not in [200, 299]
     */
    public async continueInternetProductsQueryWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void | InternetProductsResponse>> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: InternetProductsResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "InternetProductsResponse", ""
            ) as InternetProductsResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("202", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Bad request, invalid cursor", undefined, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Not found, cursor not found", undefined, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Internal server error", undefined, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void | InternetProductsResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void | InternetProductsResponse", ""
            ) as void | InternetProductsResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to getSharedInternetProducts
     * @throws ApiException if the response code was not in [200, 299]
     */
    public async getSharedInternetProductsWithHttpInfo(response: ResponseContext): Promise<HttpInfo<SharedInternetProductsResponse>> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: SharedInternetProductsResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SharedInternetProductsResponse", ""
            ) as SharedInternetProductsResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Bad request, invalid cursor", undefined, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Not found, cursor not found", undefined, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Internal server error", undefined, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: SharedInternetProductsResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SharedInternetProductsResponse", ""
            ) as SharedInternetProductsResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to initiateInternetProductsQuery
     * @throws ApiException if the response code was not in [200, 299]
     */
    public async initiateInternetProductsQueryWithHttpInfo(response: ResponseContext): Promise<HttpInfo<InternetProductsCursor>> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: InternetProductsCursor = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "InternetProductsCursor", ""
            ) as InternetProductsCursor;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Internal server error", undefined, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: InternetProductsCursor = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "InternetProductsCursor", ""
            ) as InternetProductsCursor;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to shareInternetProducts
     * @throws ApiException if the response code was not in [200, 299]
     */
    public async shareInternetProductsWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void>> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Bad request, invalid cursor or query not completed", undefined, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            throw new ApiException<undefined>(response.httpStatusCode, "Internal server error", undefined, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

}
