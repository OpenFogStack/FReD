# \KeygroupApi

All URIs are relative to *https://[node-id].nodes.mcc-f.red/v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**KeygroupGroupIdDelete**](KeygroupApi.md#KeygroupGroupIdDelete) | **Delete** /keygroup/{group_id} | Delete an existing Keygroup
[**KeygroupGroupIdPost**](KeygroupApi.md#KeygroupGroupIdPost) | **Post** /keygroup/{group_id} | Create a new Keygroup
[**KeygroupGroupIdReplicaGet**](KeygroupApi.md#KeygroupGroupIdReplicaGet) | **Get** /keygroup/{group_id}/replica | Gets all Replica Nodes for a Keygroup
[**KeygroupGroupIdReplicaNodeIdDelete**](KeygroupApi.md#KeygroupGroupIdReplicaNodeIdDelete) | **Delete** /keygroup/{group_id}/replica/{node_id} | Remove an existing replica node for a Keygroup
[**KeygroupGroupIdReplicaNodeIdPost**](KeygroupApi.md#KeygroupGroupIdReplicaNodeIdPost) | **Post** /keygroup/{group_id}/replica/{node_id} | Create a new Replica node for a Keygroup


# **KeygroupGroupIdDelete**
> KeygroupGroupIdDelete(ctx, groupId)
Delete an existing Keygroup

Deletes the Keygroup with the name `group_id` if it exists.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **KeygroupGroupIdPost**
> KeygroupGroupIdPost(ctx, groupId)
Create a new Keygroup

Creates a new Keygroup with the name `group_id` if it does not exist already.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **KeygroupGroupIdReplicaGet**
> ReplicationList KeygroupGroupIdReplicaGet(ctx, groupId)
Gets all Replica Nodes for a Keygroup

Returns replica nodes for a Keygroup with the name `group_id` if it does not exist already.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 

### Return type

[**ReplicationList**](ReplicationList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **KeygroupGroupIdReplicaNodeIdDelete**
> KeygroupGroupIdReplicaNodeIdDelete(ctx, groupId, nodeId)
Remove an existing replica node for a Keygroup

De-registers the node with the name {node_id} as a replica node for the Keygroup with the name `group_id` if it exists.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 
  **nodeId** | **string**| Name of Replica Node | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **KeygroupGroupIdReplicaNodeIdPost**
> KeygroupGroupIdReplicaNodeIdPost(ctx, groupId, nodeId)
Create a new Replica node for a Keygroup

Registers the node with the name `node_id` as a replica node for a Keygroup with the name `group_id` if it does not exist already.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 
  **nodeId** | **string**| Name of Replica Node | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

