# \DataApi

All URIs are relative to *https://[node-id].nodes.mcc-f.red/v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**KeygroupGroupIdDataItemIdDelete**](DataApi.md#KeygroupGroupIdDataItemIdDelete) | **Delete** /keygroup/{group_id}/data/{item_id} | Deletes an item value from a Keygroup
[**KeygroupGroupIdDataItemIdGet**](DataApi.md#KeygroupGroupIdDataItemIdGet) | **Get** /keygroup/{group_id}/data/{item_id} | Gets an item value from a Keygroup
[**KeygroupGroupIdDataItemIdPut**](DataApi.md#KeygroupGroupIdDataItemIdPut) | **Put** /keygroup/{group_id}/data/{item_id} | Sets an item value in a Keygroup


# **KeygroupGroupIdDataItemIdDelete**
> KeygroupGroupIdDataItemIdDelete(ctx, groupId, itemId)
Deletes an item value from a Keygroup

Deletes the value of item `item_id` in Keygroup with the name `group_id`.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 
  **itemId** | **string**| ID of item to update | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **KeygroupGroupIdDataItemIdGet**
> Item KeygroupGroupIdDataItemIdGet(ctx, groupId, itemId)
Gets an item value from a Keygroup

Gets the value of item `item_id` in Keygroup with the name `group_id`.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 
  **itemId** | **string**| ID of item to update | 

### Return type

[**Item**](Item.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **KeygroupGroupIdDataItemIdPut**
> KeygroupGroupIdDataItemIdPut(ctx, groupId, itemId, body)
Sets an item value in a Keygroup

Sets the value of item `item_id` in Keygroup with the name `group_id` to the provided value.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 
  **itemId** | **string**| ID of item to update | 
  **body** | [**Item**](Item.md)| Data that should be saved in this item. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: text/plain, application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

