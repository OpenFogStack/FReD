# \TriggersApi

All URIs are relative to *https://[node-id].nodes.mcc-f.red/v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**KeygroupGroupIdTriggersGet**](TriggersApi.md#KeygroupGroupIdTriggersGet) | **Get** /keygroup/{group_id}/triggers | Gets all Trigger Nodes for a Keygroup
[**KeygroupGroupIdTriggersTriggerNodeIdDelete**](TriggersApi.md#KeygroupGroupIdTriggersTriggerNodeIdDelete) | **Delete** /keygroup/{group_id}/triggers/{trigger_node_id} | Remove an existing trigger node for a Keygroup
[**KeygroupGroupIdTriggersTriggerNodeIdPost**](TriggersApi.md#KeygroupGroupIdTriggersTriggerNodeIdPost) | **Post** /keygroup/{group_id}/triggers/{trigger_node_id} | Create a new Trigger node for a Keygroup


# **KeygroupGroupIdTriggersGet**
> TriggerList KeygroupGroupIdTriggersGet(ctx, groupId)
Gets all Trigger Nodes for a Keygroup

Returns trigger nodes for a Keygroup with the name `group_id`.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 

### Return type

[**TriggerList**](TriggerList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **KeygroupGroupIdTriggersTriggerNodeIdDelete**
> KeygroupGroupIdTriggersTriggerNodeIdDelete(ctx, groupId, triggerNodeId)
Remove an existing trigger node for a Keygroup

De-registers the node with the name {trigger_node_id} as a trigger node for the Keygroup with the name `group_id` if it exists.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 
  **triggerNodeId** | **string**| Name of Trigger Node | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **KeygroupGroupIdTriggersTriggerNodeIdPost**
> KeygroupGroupIdTriggersTriggerNodeIdPost(ctx, groupId, triggerNodeId, host)
Create a new Trigger node for a Keygroup

Registers the trigger node with the name `trigger_node_id` and host `host` as a trigger node for a Keygroup with the name `group_id` if it does not exist already.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**| Name of Keygroup | 
  **triggerNodeId** | **string**| Name of Trigger Node | 
  **host** | [**TriggerNode**](TriggerNode.md)| Host of the trigger node | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

