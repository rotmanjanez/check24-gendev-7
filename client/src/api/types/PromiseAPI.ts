import { HttpInfo } from '../http/http';
import type { Configuration, ConfigurationOptions, PromiseConfigurationOptions } from '../configuration'
import { PromiseMiddlewareWrapper } from '../middleware';

import { Address } from '../models/Address';
import { Health } from '../models/Health';
import { InternetProductsCursor } from '../models/InternetProductsCursor';
import { InternetProductsResponse } from '../models/InternetProductsResponse';
import { SharedInternetProductsResponse } from '../models/SharedInternetProductsResponse';
import { Version } from '../models/Version';
import { ObservableHealthApi } from './ObservableAPI';

import { HealthApiRequestFactory, HealthApiResponseProcessor } from "../apis/HealthApi";
export class PromiseHealthApi {
    private api: ObservableHealthApi

    public constructor(
        configuration: Configuration,
        requestFactory?: HealthApiRequestFactory,
        responseProcessor?: HealthApiResponseProcessor
    ) {
        this.api = new ObservableHealthApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Returns the status of the API
     * Health check endpoint
     */
    public healthCheckWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<Health>> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.healthCheckWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     * Returns the status of the API
     * Health check endpoint
     */
    public healthCheck(_options?: PromiseConfigurationOptions): Promise<Health> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.healthCheck(observableOptions);
        return result.toPromise();
    }


}



import { ObservableInternetProductsApi } from './ObservableAPI';

import { InternetProductsApiRequestFactory, InternetProductsApiResponseProcessor } from "../apis/InternetProductsApi";
export class PromiseInternetProductsApi {
    private api: ObservableInternetProductsApi

    public constructor(
        configuration: Configuration,
        requestFactory?: InternetProductsApiRequestFactory,
        responseProcessor?: InternetProductsApiResponseProcessor
    ) {
        this.api = new ObservableInternetProductsApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Fetches the next batch of internet products using a continuation cursor
     * @param cursor Cursor to continue fetching products
     */
    public continueInternetProductsQueryWithHttpInfo(cursor: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<InternetProductsResponse>> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.continueInternetProductsQueryWithHttpInfo(cursor, observableOptions);
        return result.toPromise();
    }

    /**
     * Fetches the next batch of internet products using a continuation cursor
     * @param cursor Cursor to continue fetching products
     */
    public continueInternetProductsQuery(cursor: string, _options?: PromiseConfigurationOptions): Promise<InternetProductsResponse> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.continueInternetProductsQuery(cursor, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieves the shared internet products using a given cursor
     * @param cursor Cursor to retrieve the shared products
     */
    public getSharedInternetProductsWithHttpInfo(cursor: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<SharedInternetProductsResponse>> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.getSharedInternetProductsWithHttpInfo(cursor, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieves the shared internet products using a given cursor
     * @param cursor Cursor to retrieve the shared products
     */
    public getSharedInternetProducts(cursor: string, _options?: PromiseConfigurationOptions): Promise<SharedInternetProductsResponse> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.getSharedInternetProducts(cursor, observableOptions);
        return result.toPromise();
    }

    /**
     * Initiates retrieval of internet products and returns a product version and a cursor to retrieve the first batch of products
     * @param address
     * @param [providers] Providers to filter the products
     */
    public initiateInternetProductsQueryWithHttpInfo(address: Address, providers?: Array<string>, _options?: PromiseConfigurationOptions): Promise<HttpInfo<InternetProductsCursor>> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.initiateInternetProductsQueryWithHttpInfo(address, providers, observableOptions);
        return result.toPromise();
    }

    /**
     * Initiates retrieval of internet products and returns a product version and a cursor to retrieve the first batch of products
     * @param address
     * @param [providers] Providers to filter the products
     */
    public initiateInternetProductsQuery(address: Address, providers?: Array<string>, _options?: PromiseConfigurationOptions): Promise<InternetProductsCursor> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.initiateInternetProductsQuery(address, providers, observableOptions);
        return result.toPromise();
    }

    /**
     * Shares the internet products with a given cursor. This cursor must be the same as the one returned by the initial query and the query must have completed.
     * @param cursor Cursor to share the products
     */
    public shareInternetProductsWithHttpInfo(cursor: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.shareInternetProductsWithHttpInfo(cursor, observableOptions);
        return result.toPromise();
    }

    /**
     * Shares the internet products with a given cursor. This cursor must be the same as the one returned by the initial query and the query must have completed.
     * @param cursor Cursor to share the products
     */
    public shareInternetProducts(cursor: string, _options?: PromiseConfigurationOptions): Promise<void> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.shareInternetProducts(cursor, observableOptions);
        return result.toPromise();
    }


}



import { ObservableSystemApi } from './ObservableAPI';

import { SystemApiRequestFactory, SystemApiResponseProcessor } from "../apis/SystemApi";
export class PromiseSystemApi {
    private api: ObservableSystemApi

    public constructor(
        configuration: Configuration,
        requestFactory?: SystemApiRequestFactory,
        responseProcessor?: SystemApiResponseProcessor
    ) {
        this.api = new ObservableSystemApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Returns version information about the API
     * Version information endpoint
     */
    public getVersionWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<Version>> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.getVersionWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     * Returns version information about the API
     * Version information endpoint
     */
    public getVersion(_options?: PromiseConfigurationOptions): Promise<Version> {
        let observableOptions: undefined | ConfigurationOptions
        if (_options) {
            observableOptions = {
                baseServer: _options.baseServer,
                httpApi: _options.httpApi,
                middleware: _options.middleware?.map(
                    m => new PromiseMiddlewareWrapper(m)
                ),
                middlewareMergeStrategy: _options.middlewareMergeStrategy,
                authMethods: _options.authMethods
            }
        }
        const result = this.api.getVersion(observableOptions);
        return result.toPromise();
    }


}



