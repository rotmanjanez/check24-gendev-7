import { HttpInfo } from '../http/http';
import type { Configuration, ConfigurationOptions } from '../configuration'

import { Address } from '../models/Address';
import { Health } from '../models/Health';
import { InternetProductsCursor } from '../models/InternetProductsCursor';
import { InternetProductsResponse } from '../models/InternetProductsResponse';
import { SharedInternetProductsResponse } from '../models/SharedInternetProductsResponse';
import { Version } from '../models/Version';

import { ObservableHealthApi } from "./ObservableAPI";
import { HealthApiRequestFactory, HealthApiResponseProcessor } from "../apis/HealthApi";

export class ObjectHealthApi {
    private api: ObservableHealthApi

    public constructor(configuration: Configuration, requestFactory?: HealthApiRequestFactory, responseProcessor?: HealthApiResponseProcessor) {
        this.api = new ObservableHealthApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Returns the status of the API
     * Health check endpoint
     * @param param the request object
     */
    public healthCheckWithHttpInfo(options?: ConfigurationOptions): Promise<HttpInfo<Health>> {
        return this.api.healthCheckWithHttpInfo(options).toPromise();
    }

    /**
     * Returns the status of the API
     * Health check endpoint
     * @param param the request object
     */
    public healthCheck(options?: ConfigurationOptions): Promise<Health> {
        return this.api.healthCheck(options).toPromise();
    }

}

import { ObservableInternetProductsApi } from "./ObservableAPI";
import { InternetProductsApiRequestFactory, InternetProductsApiResponseProcessor } from "../apis/InternetProductsApi";

export interface InternetProductsApiContinueInternetProductsQueryRequest {
    /**
     * Cursor to continue fetching products
     * Defaults to: undefined
     * @type string
     * @memberof InternetProductsApicontinueInternetProductsQuery
     */
    cursor: string
}

export interface InternetProductsApiGetSharedInternetProductsRequest {
    /**
     * Cursor to retrieve the shared products
     * Defaults to: undefined
     * @type string
     * @memberof InternetProductsApigetSharedInternetProducts
     */
    cursor: string
}

export interface InternetProductsApiInitiateInternetProductsQueryRequest {
    /**
     * 
     * @type Address
     * @memberof InternetProductsApiinitiateInternetProductsQuery
     */
    address: Address
    /**
     * Providers to filter the products
     * Defaults to: undefined
     * @type Array&lt;string&gt;
     * @memberof InternetProductsApiinitiateInternetProductsQuery
     */
    providers?: Array<string>
}

export interface InternetProductsApiShareInternetProductsRequest {
    /**
     * Cursor to share the products
     * Defaults to: undefined
     * @type string
     * @memberof InternetProductsApishareInternetProducts
     */
    cursor: string
}

export class ObjectInternetProductsApi {
    private api: ObservableInternetProductsApi

    public constructor(configuration: Configuration, requestFactory?: InternetProductsApiRequestFactory, responseProcessor?: InternetProductsApiResponseProcessor) {
        this.api = new ObservableInternetProductsApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Fetches the next batch of internet products using a continuation cursor
     * @param param the request object
     */
    public continueInternetProductsQueryWithHttpInfo(param: InternetProductsApiContinueInternetProductsQueryRequest, options?: ConfigurationOptions): Promise<HttpInfo<InternetProductsResponse>> {
        return this.api.continueInternetProductsQueryWithHttpInfo(param.cursor, options).toPromise();
    }

    /**
     * Fetches the next batch of internet products using a continuation cursor
     * @param param the request object
     */
    public continueInternetProductsQuery(param: InternetProductsApiContinueInternetProductsQueryRequest, options?: ConfigurationOptions): Promise<InternetProductsResponse> {
        return this.api.continueInternetProductsQuery(param.cursor, options).toPromise();
    }

    /**
     * Retrieves the shared internet products using a given cursor
     * @param param the request object
     */
    public getSharedInternetProductsWithHttpInfo(param: InternetProductsApiGetSharedInternetProductsRequest, options?: ConfigurationOptions): Promise<HttpInfo<SharedInternetProductsResponse>> {
        return this.api.getSharedInternetProductsWithHttpInfo(param.cursor, options).toPromise();
    }

    /**
     * Retrieves the shared internet products using a given cursor
     * @param param the request object
     */
    public getSharedInternetProducts(param: InternetProductsApiGetSharedInternetProductsRequest, options?: ConfigurationOptions): Promise<SharedInternetProductsResponse> {
        return this.api.getSharedInternetProducts(param.cursor, options).toPromise();
    }

    /**
     * Initiates retrieval of internet products and returns a product version and a cursor to retrieve the first batch of products
     * @param param the request object
     */
    public initiateInternetProductsQueryWithHttpInfo(param: InternetProductsApiInitiateInternetProductsQueryRequest, options?: ConfigurationOptions): Promise<HttpInfo<InternetProductsCursor>> {
        return this.api.initiateInternetProductsQueryWithHttpInfo(param.address, param.providers, options).toPromise();
    }

    /**
     * Initiates retrieval of internet products and returns a product version and a cursor to retrieve the first batch of products
     * @param param the request object
     */
    public initiateInternetProductsQuery(param: InternetProductsApiInitiateInternetProductsQueryRequest, options?: ConfigurationOptions): Promise<InternetProductsCursor> {
        return this.api.initiateInternetProductsQuery(param.address, param.providers, options).toPromise();
    }

    /**
     * Shares the internet products with a given cursor. This cursor must be the same as the one returned by the initial query and the query must have completed.
     * @param param the request object
     */
    public shareInternetProductsWithHttpInfo(param: InternetProductsApiShareInternetProductsRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.shareInternetProductsWithHttpInfo(param.cursor, options).toPromise();
    }

    /**
     * Shares the internet products with a given cursor. This cursor must be the same as the one returned by the initial query and the query must have completed.
     * @param param the request object
     */
    public shareInternetProducts(param: InternetProductsApiShareInternetProductsRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.shareInternetProducts(param.cursor, options).toPromise();
    }

}

import { ObservableSystemApi } from "./ObservableAPI";
import { SystemApiRequestFactory, SystemApiResponseProcessor } from "../apis/SystemApi";

export class ObjectSystemApi {
    private api: ObservableSystemApi

    public constructor(configuration: Configuration, requestFactory?: SystemApiRequestFactory, responseProcessor?: SystemApiResponseProcessor) {
        this.api = new ObservableSystemApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Returns version information about the API
     * Version information endpoint
     * @param param the request object
     */
    public getVersionWithHttpInfo(options?: ConfigurationOptions): Promise<HttpInfo<Version>> {
        return this.api.getVersionWithHttpInfo(options).toPromise();
    }

    /**
     * Returns version information about the API
     * Version information endpoint
     * @param param the request object
     */
    public getVersion(options?: ConfigurationOptions): Promise<Version> {
        return this.api.getVersion(options).toPromise();
    }

}
