# .HealthApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**healthCheck**](HealthApi.md#healthCheck) | **GET** /health | Health check endpoint


# **healthCheck**
> Health healthCheck()

Returns the status of the API

### Example


```typescript
import { createConfiguration, HealthApi } from '';

const configuration = createConfiguration();
const apiInstance = new HealthApi(configuration);

const request = {};

const data = await apiInstance.healthCheck(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


### Return type

**Health**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successful health check |  -  |
**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


