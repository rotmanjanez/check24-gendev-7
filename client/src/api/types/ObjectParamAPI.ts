import { HttpInfo } from '../http/http';
import type { Configuration, ConfigurationOptions } from '../configuration'

import { Health } from '../models/Health';

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
