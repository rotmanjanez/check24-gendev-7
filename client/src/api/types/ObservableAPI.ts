import { ResponseContext, RequestContext, HttpInfo } from '../http/http';
import type { Configuration, ConfigurationOptions } from '../configuration'
import type { Middleware } from '../middleware';
import { Observable, of, from } from '../rxjsStub';
import { mergeMap, map } from '../rxjsStub';
import { Address } from '../models/Address';
import { Health } from '../models/Health';
import { InternetProductsCursor } from '../models/InternetProductsCursor';
import { InternetProductsResponse } from '../models/InternetProductsResponse';
import { SharedInternetProductsResponse } from '../models/SharedInternetProductsResponse';
import { Version } from '../models/Version';

import { HealthApiRequestFactory, HealthApiResponseProcessor } from "../apis/HealthApi";
export class ObservableHealthApi {
  private requestFactory: HealthApiRequestFactory;
  private responseProcessor: HealthApiResponseProcessor;
  private configuration: Configuration;

  public constructor(
    configuration: Configuration,
    requestFactory?: HealthApiRequestFactory,
    responseProcessor?: HealthApiResponseProcessor
  ) {
    this.configuration = configuration;
    this.requestFactory = requestFactory || new HealthApiRequestFactory(configuration);
    this.responseProcessor = responseProcessor || new HealthApiResponseProcessor();
  }

  /**
   * Returns the status of the API
   * Health check endpoint
   */
  public healthCheckWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<Health>> {
    let _config = this.configuration;
    let allMiddleware: Middleware[] = [];
    if (_options && _options.middleware) {
      const middlewareMergeStrategy = _options.middlewareMergeStrategy || 'replace' // default to replace behavior
      // call-time middleware provided
      const calltimeMiddleware: Middleware[] = _options.middleware;

      switch (middlewareMergeStrategy) {
        case 'append':
          allMiddleware = this.configuration.middleware.concat(calltimeMiddleware);
          break;
        case 'prepend':
          allMiddleware = calltimeMiddleware.concat(this.configuration.middleware)
          break;
        case 'replace':
          allMiddleware = calltimeMiddleware
          break;
        default:
          throw new Error(`unrecognized middleware merge strategy '${middlewareMergeStrategy}'`)
      }
    }
    if (_options) {
      _config = {
        baseServer: _options.baseServer || this.configuration.baseServer,
        httpApi: _options.httpApi || this.configuration.httpApi,
        authMethods: _options.authMethods || this.configuration.authMethods,
        middleware: allMiddleware || this.configuration.middleware
      };
    }

    const requestContextPromise = this.requestFactory.healthCheck(_config);
    // build promise chain
    let middlewarePreObservable = from<RequestContext>(requestContextPromise);
    for (const middleware of allMiddleware) {
      middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
    }

    return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
      pipe(mergeMap((response: ResponseContext) => {
        let middlewarePostObservable = of(response);
        for (const middleware of allMiddleware.reverse()) {
          middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
        }
        return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.healthCheckWithHttpInfo(rsp)));
      }));
  }

  /**
   * Returns the status of the API
   * Health check endpoint
   */
  public healthCheck(_options?: ConfigurationOptions): Observable<Health> {
    return this.healthCheckWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<Health>) => apiResponse.data));
  }

}

import { InternetProductsApiRequestFactory, InternetProductsApiResponseProcessor } from "../apis/InternetProductsApi";
export class ObservableInternetProductsApi {
  private requestFactory: InternetProductsApiRequestFactory;
  private responseProcessor: InternetProductsApiResponseProcessor;
  private configuration: Configuration;

  public constructor(
    configuration: Configuration,
    requestFactory?: InternetProductsApiRequestFactory,
    responseProcessor?: InternetProductsApiResponseProcessor
  ) {
    this.configuration = configuration;
    this.requestFactory = requestFactory || new InternetProductsApiRequestFactory(configuration);
    this.responseProcessor = responseProcessor || new InternetProductsApiResponseProcessor();
  }

  /**
   * Fetches the next batch of internet products using a continuation cursor
   * @param cursor Cursor to continue fetching products
   */
  public continueInternetProductsQueryWithHttpInfo(cursor: string, _options?: ConfigurationOptions): Observable<HttpInfo<InternetProductsResponse>> {
    let _config = this.configuration;
    let allMiddleware: Middleware[] = [];
    if (_options && _options.middleware) {
      const middlewareMergeStrategy = _options.middlewareMergeStrategy || 'replace' // default to replace behavior
      // call-time middleware provided
      const calltimeMiddleware: Middleware[] = _options.middleware;

      switch (middlewareMergeStrategy) {
        case 'append':
          allMiddleware = this.configuration.middleware.concat(calltimeMiddleware);
          break;
        case 'prepend':
          allMiddleware = calltimeMiddleware.concat(this.configuration.middleware)
          break;
        case 'replace':
          allMiddleware = calltimeMiddleware
          break;
        default:
          throw new Error(`unrecognized middleware merge strategy '${middlewareMergeStrategy}'`)
      }
    }
    if (_options) {
      _config = {
        baseServer: _options.baseServer || this.configuration.baseServer,
        httpApi: _options.httpApi || this.configuration.httpApi,
        authMethods: _options.authMethods || this.configuration.authMethods,
        middleware: allMiddleware || this.configuration.middleware
      };
    }

    const requestContextPromise = this.requestFactory.continueInternetProductsQuery(cursor, _config);
    // build promise chain
    let middlewarePreObservable = from<RequestContext>(requestContextPromise);
    for (const middleware of allMiddleware) {
      middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
    }

    return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
      pipe(mergeMap((response: ResponseContext) => {
        let middlewarePostObservable = of(response);
        for (const middleware of allMiddleware.reverse()) {
          middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
        }
        return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.continueInternetProductsQueryWithHttpInfo(rsp)));
      }));
  }

  /**
   * Fetches the next batch of internet products using a continuation cursor
   * @param cursor Cursor to continue fetching products
   */
  public continueInternetProductsQuery(cursor: string, _options?: ConfigurationOptions): Observable<InternetProductsResponse> {
    return this.continueInternetProductsQueryWithHttpInfo(cursor, _options).pipe(map((apiResponse: HttpInfo<InternetProductsResponse>) => apiResponse.data));
  }

  /**
   * Retrieves the shared internet products using a given cursor
   * @param cursor Cursor to retrieve the shared products
   */
  public getSharedInternetProductsWithHttpInfo(cursor: string, _options?: ConfigurationOptions): Observable<HttpInfo<SharedInternetProductsResponse>> {
    let _config = this.configuration;
    let allMiddleware: Middleware[] = [];
    if (_options && _options.middleware) {
      const middlewareMergeStrategy = _options.middlewareMergeStrategy || 'replace' // default to replace behavior
      // call-time middleware provided
      const calltimeMiddleware: Middleware[] = _options.middleware;

      switch (middlewareMergeStrategy) {
        case 'append':
          allMiddleware = this.configuration.middleware.concat(calltimeMiddleware);
          break;
        case 'prepend':
          allMiddleware = calltimeMiddleware.concat(this.configuration.middleware)
          break;
        case 'replace':
          allMiddleware = calltimeMiddleware
          break;
        default:
          throw new Error(`unrecognized middleware merge strategy '${middlewareMergeStrategy}'`)
      }
    }
    if (_options) {
      _config = {
        baseServer: _options.baseServer || this.configuration.baseServer,
        httpApi: _options.httpApi || this.configuration.httpApi,
        authMethods: _options.authMethods || this.configuration.authMethods,
        middleware: allMiddleware || this.configuration.middleware
      };
    }

    const requestContextPromise = this.requestFactory.getSharedInternetProducts(cursor, _config);
    // build promise chain
    let middlewarePreObservable = from<RequestContext>(requestContextPromise);
    for (const middleware of allMiddleware) {
      middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
    }

    return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
      pipe(mergeMap((response: ResponseContext) => {
        let middlewarePostObservable = of(response);
        for (const middleware of allMiddleware.reverse()) {
          middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
        }
        return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.getSharedInternetProductsWithHttpInfo(rsp)));
      }));
  }

  /**
   * Retrieves the shared internet products using a given cursor
   * @param cursor Cursor to retrieve the shared products
   */
  public getSharedInternetProducts(cursor: string, _options?: ConfigurationOptions): Observable<SharedInternetProductsResponse> {
    return this.getSharedInternetProductsWithHttpInfo(cursor, _options).pipe(map((apiResponse: HttpInfo<SharedInternetProductsResponse>) => apiResponse.data));
  }

  /**
   * Initiates retrieval of internet products and returns a product version and a cursor to retrieve the first batch of products
   * @param address
   * @param [providers] Providers to filter the products
   */
  public initiateInternetProductsQueryWithHttpInfo(address: Address, providers?: Array<string>, _options?: ConfigurationOptions): Observable<HttpInfo<InternetProductsCursor>> {
    let _config = this.configuration;
    let allMiddleware: Middleware[] = [];
    if (_options && _options.middleware) {
      const middlewareMergeStrategy = _options.middlewareMergeStrategy || 'replace' // default to replace behavior
      // call-time middleware provided
      const calltimeMiddleware: Middleware[] = _options.middleware;

      switch (middlewareMergeStrategy) {
        case 'append':
          allMiddleware = this.configuration.middleware.concat(calltimeMiddleware);
          break;
        case 'prepend':
          allMiddleware = calltimeMiddleware.concat(this.configuration.middleware)
          break;
        case 'replace':
          allMiddleware = calltimeMiddleware
          break;
        default:
          throw new Error(`unrecognized middleware merge strategy '${middlewareMergeStrategy}'`)
      }
    }
    if (_options) {
      _config = {
        baseServer: _options.baseServer || this.configuration.baseServer,
        httpApi: _options.httpApi || this.configuration.httpApi,
        authMethods: _options.authMethods || this.configuration.authMethods,
        middleware: allMiddleware || this.configuration.middleware
      };
    }

    const requestContextPromise = this.requestFactory.initiateInternetProductsQuery(address, providers, _config);
    // build promise chain
    let middlewarePreObservable = from<RequestContext>(requestContextPromise);
    for (const middleware of allMiddleware) {
      middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
    }

    return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
      pipe(mergeMap((response: ResponseContext) => {
        let middlewarePostObservable = of(response);
        for (const middleware of allMiddleware.reverse()) {
          middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
        }
        return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.initiateInternetProductsQueryWithHttpInfo(rsp)));
      }));
  }

  /**
   * Initiates retrieval of internet products and returns a product version and a cursor to retrieve the first batch of products
   * @param address
   * @param [providers] Providers to filter the products
   */
  public initiateInternetProductsQuery(address: Address, providers?: Array<string>, _options?: ConfigurationOptions): Observable<InternetProductsCursor> {
    return this.initiateInternetProductsQueryWithHttpInfo(address, providers, _options).pipe(map((apiResponse: HttpInfo<InternetProductsCursor>) => apiResponse.data));
  }

  /**
   * Shares the internet products with a given cursor. This cursor must be the same as the one returned by the initial query and the query must have completed.
   * @param cursor Cursor to share the products
   */
  public shareInternetProductsWithHttpInfo(cursor: string, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
    let _config = this.configuration;
    let allMiddleware: Middleware[] = [];
    if (_options && _options.middleware) {
      const middlewareMergeStrategy = _options.middlewareMergeStrategy || 'replace' // default to replace behavior
      // call-time middleware provided
      const calltimeMiddleware: Middleware[] = _options.middleware;

      switch (middlewareMergeStrategy) {
        case 'append':
          allMiddleware = this.configuration.middleware.concat(calltimeMiddleware);
          break;
        case 'prepend':
          allMiddleware = calltimeMiddleware.concat(this.configuration.middleware)
          break;
        case 'replace':
          allMiddleware = calltimeMiddleware
          break;
        default:
          throw new Error(`unrecognized middleware merge strategy '${middlewareMergeStrategy}'`)
      }
    }
    if (_options) {
      _config = {
        baseServer: _options.baseServer || this.configuration.baseServer,
        httpApi: _options.httpApi || this.configuration.httpApi,
        authMethods: _options.authMethods || this.configuration.authMethods,
        middleware: allMiddleware || this.configuration.middleware
      };
    }

    const requestContextPromise = this.requestFactory.shareInternetProducts(cursor, _config);
    // build promise chain
    let middlewarePreObservable = from<RequestContext>(requestContextPromise);
    for (const middleware of allMiddleware) {
      middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
    }

    return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
      pipe(mergeMap((response: ResponseContext) => {
        let middlewarePostObservable = of(response);
        for (const middleware of allMiddleware.reverse()) {
          middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
        }
        return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.shareInternetProductsWithHttpInfo(rsp)));
      }));
  }

  /**
   * Shares the internet products with a given cursor. This cursor must be the same as the one returned by the initial query and the query must have completed.
   * @param cursor Cursor to share the products
   */
  public shareInternetProducts(cursor: string, _options?: ConfigurationOptions): Observable<void> {
    return this.shareInternetProductsWithHttpInfo(cursor, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
  }

}

import { SystemApiRequestFactory, SystemApiResponseProcessor } from "../apis/SystemApi";
export class ObservableSystemApi {
  private requestFactory: SystemApiRequestFactory;
  private responseProcessor: SystemApiResponseProcessor;
  private configuration: Configuration;

  public constructor(
    configuration: Configuration,
    requestFactory?: SystemApiRequestFactory,
    responseProcessor?: SystemApiResponseProcessor
  ) {
    this.configuration = configuration;
    this.requestFactory = requestFactory || new SystemApiRequestFactory(configuration);
    this.responseProcessor = responseProcessor || new SystemApiResponseProcessor();
  }

  /**
   * Returns version information about the API
   * Version information endpoint
   */
  public getVersionWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<Version>> {
    let _config = this.configuration;
    let allMiddleware: Middleware[] = [];
    if (_options && _options.middleware) {
      const middlewareMergeStrategy = _options.middlewareMergeStrategy || 'replace' // default to replace behavior
      // call-time middleware provided
      const calltimeMiddleware: Middleware[] = _options.middleware;

      switch (middlewareMergeStrategy) {
        case 'append':
          allMiddleware = this.configuration.middleware.concat(calltimeMiddleware);
          break;
        case 'prepend':
          allMiddleware = calltimeMiddleware.concat(this.configuration.middleware)
          break;
        case 'replace':
          allMiddleware = calltimeMiddleware
          break;
        default:
          throw new Error(`unrecognized middleware merge strategy '${middlewareMergeStrategy}'`)
      }
    }
    if (_options) {
      _config = {
        baseServer: _options.baseServer || this.configuration.baseServer,
        httpApi: _options.httpApi || this.configuration.httpApi,
        authMethods: _options.authMethods || this.configuration.authMethods,
        middleware: allMiddleware || this.configuration.middleware
      };
    }

    const requestContextPromise = this.requestFactory.getVersion(_config);
    // build promise chain
    let middlewarePreObservable = from<RequestContext>(requestContextPromise);
    for (const middleware of allMiddleware) {
      middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
    }

    return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
      pipe(mergeMap((response: ResponseContext) => {
        let middlewarePostObservable = of(response);
        for (const middleware of allMiddleware.reverse()) {
          middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
        }
        return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.getVersionWithHttpInfo(rsp)));
      }));
  }

  /**
   * Returns version information about the API
   * Version information endpoint
   */
  public getVersion(_options?: ConfigurationOptions): Observable<Version> {
    return this.getVersionWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<Version>) => apiResponse.data));
  }

}
