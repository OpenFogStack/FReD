# \ReplicationApi

All URIs are relative to *https://[node-id].nodes.mcc-f.red/v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**KeygroupGroupIdReplicaGet**](ReplicationApi.md#KeygroupGroupIdReplicaGet) | **Get** /keygroup/{group_id}/replica | Gets all Replica Nodes for a Keygroup
[**KeygroupGroupIdReplicaNodeIdDelete**](ReplicationApi.md#KeygroupGroupIdReplicaNodeIdDelete) | **Delete** /keygroup/{group_id}/replica/{node_id} | Remove an existing replica node for a Keygroup
[**KeygroupGroupIdReplicaNodeIdPost**](ReplicationApi.md#KeygroupGroupIdReplicaNodeIdPost) | **Post** /keygroup/{group_id}/replica/{node_id} | Create a new Replica node for a Keygroup
[**ReplicaGet**](ReplicationApi.md#ReplicaGet) | **Get** /replica | Gets all Replica Nodes
[**ReplicaNodeIdDelete**](ReplicationApi.md#ReplicaNodeIdDelete) | **Delete** /replica/{node_id} | Remove an existing replica node
[**ReplicaNodeIdGet**](ReplicationApi.md#ReplicaNodeIdGet) | **Get** /replica/{node_id} | Gets a Replica Node


# **KeygroupGroupIdReplicaGet**
> ReplicationList KeygroupGroupIdReplicaGet(ctx, groupId)
Gets all Replica Nodes for a Keygroup

Returns replica nodes for a Keygroup with the name `group_id`.

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

# **ReplicaGet**
> ReplicationList ReplicaGet(ctx, )
Gets all Replica Nodes

Returns all known replica nodes.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ReplicationList**](ReplicationList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReplicaNodeIdDelete**
> ReplicaNodeIdDelete(ctx, nodeId)
Remove an existing replica node

Removes the replica node `node_id` if it exists.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **nodeId** | **string**| ID of node to delete | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReplicaNodeIdGet**
> Node ReplicaNodeIdGet(ctx, nodeId)
Gets a Replica Node

Returns the replica node `node_id` if it exists.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **nodeId** | **string**| ID of node to return | 

### Return type

[**Node**](Node.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

