# .InternetProductsApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**continueInternetProductsQuery**](InternetProductsApi.md#continueInternetProductsQuery) | **GET** /internet-products/continue | 
[**getSharedInternetProducts**](InternetProductsApi.md#getSharedInternetProducts) | **GET** /internet-products/share/{cursor} | 
[**initiateInternetProductsQuery**](InternetProductsApi.md#initiateInternetProductsQuery) | **POST** /internet-products | 
[**shareInternetProducts**](InternetProductsApi.md#shareInternetProducts) | **POST** /internet-products/share/{cursor} | 


# **continueInternetProductsQuery**
> void | InternetProductsResponse continueInternetProductsQuery()

Fetches the next batch of internet products using a continuation cursor

### Example


```typescript
import { createConfiguration, InternetProductsApi } from '';
import type { InternetProductsApiContinueInternetProductsQueryRequest } from '';

const configuration = createConfiguration();
const apiInstance = new InternetProductsApi(configuration);

const request: InternetProductsApiContinueInternetProductsQueryRequest = {
    // Cursor to continue fetching products
  cursor: "cursor_example",
};

const data = await apiInstance.continueInternetProductsQuery(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cursor** | [**string**] | Cursor to continue fetching products | defaults to undefined


### Return type

**void | InternetProductsResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Next batch of internet products |  -  |
**202** | Query is still in progress, no products available yet |  -  |
**400** | Bad request, invalid cursor |  -  |
**404** | Not found, cursor not found |  -  |
**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **getSharedInternetProducts**
> SharedInternetProductsResponse getSharedInternetProducts()

Retrieves the shared internet products using a given cursor

### Example


```typescript
import { createConfiguration, InternetProductsApi } from '';
import type { InternetProductsApiGetSharedInternetProductsRequest } from '';

const configuration = createConfiguration();
const apiInstance = new InternetProductsApi(configuration);

const request: InternetProductsApiGetSharedInternetProductsRequest = {
    // Cursor to retrieve the shared products
  cursor: "cursor_example",
};

const data = await apiInstance.getSharedInternetProducts(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cursor** | [**string**] | Cursor to retrieve the shared products | defaults to undefined


### Return type

**SharedInternetProductsResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successful retrieval of shared internet products |  -  |
**400** | Bad request, invalid cursor |  -  |
**404** | Not found, cursor not found |  -  |
**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **initiateInternetProductsQuery**
> InternetProductsCursor initiateInternetProductsQuery(address)

Initiates retrieval of internet products and returns a product version and a cursor to retrieve the first batch of products

### Example


```typescript
import { createConfiguration, InternetProductsApi } from '';
import type { InternetProductsApiInitiateInternetProductsQueryRequest } from '';

const configuration = createConfiguration();
const apiInstance = new InternetProductsApi(configuration);

const request: InternetProductsApiInitiateInternetProductsQueryRequest = {
  
  address: {
    street: "street_example",
    houseNumber: "houseNumber_example",
    city: "city_example",
    postalCode: "postalCode_example",
    countryCode: "DE",
  },
    // Providers to filter the products (optional)
  providers: [
    "providers_example",
  ],
};

const data = await apiInstance.initiateInternetProductsQuery(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **address** | **Address**|  |
 **providers** | **Array&lt;string&gt;** | Providers to filter the products | (optional) defaults to undefined


### Return type

**InternetProductsCursor**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Initial batch of internet products with a continuation cursor |  -  |
**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **shareInternetProducts**
> void shareInternetProducts()

Shares the internet products with a given cursor. This cursor must be the same as the one returned by the initial query and the query must have completed.

### Example


```typescript
import { createConfiguration, InternetProductsApi } from '';
import type { InternetProductsApiShareInternetProductsRequest } from '';

const configuration = createConfiguration();
const apiInstance = new InternetProductsApi(configuration);

const request: InternetProductsApiShareInternetProductsRequest = {
    // Cursor to share the products
  cursor: "cursor_example",
};

const data = await apiInstance.shareInternetProducts(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cursor** | [**string**] | Cursor to share the products | defaults to undefined


### Return type

**void**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successful sharing of internet products |  -  |
**400** | Bad request, invalid cursor or query not completed |  -  |
**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


