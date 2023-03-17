package osutils

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"time"

	osapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//	jsonDoc (byte): The request body in form of bytes
//
// Description:
//
//	Calls the osapi IndexRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func IndexMetrics(ctx context.Context, jsonDoc []byte) (*osapi.Response, error) {
	return doRetry(ctx, osapi.IndexRequest{
		Index:        IndexName,
		DocumentType: "_doc",
		Body:         bytes.NewReader(jsonDoc),
		Refresh:      "wait_for",
	})
}

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi ClusterStatsRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func GetClusterStats(ctx context.Context) (*osapi.Response, error) {
	return doRetry(ctx, osapi.ClusterStatsRequest{})
}

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi ClusterHealthRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func GetClusterHealth(ctx context.Context, WaitForShards *bool) (*osapi.Response, error) {
	return doRetry(ctx, osapi.ClusterHealthRequest{
		WaitForNoInitializingShards: WaitForShards,
		WaitForNoRelocatingShards:   WaitForShards,
		Timeout:                     time.Duration(90 * time.Second),
	})
}

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi ClusterStateRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func GetClusterState(ctx context.Context) (*osapi.Response, error) {
	return doRetry(ctx, osapi.ClusterStateRequest{})
}

// Input:
//
//	nodes ([]string): The list of nodes for which the stats needs to be fetched
//	metrics ([]string): The list of metrics that needs to be fetched for the specified node/s
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi NodeStatsRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func GetNodeStats(ctx context.Context, nodes []string, metrics []string) (*osapi.Response, error) {
	return doRetry(ctx, osapi.NodesStatsRequest{
		Pretty: true,
		NodeID: nodes,
		Metric: metrics,
	})
}

// Input:
//
//	nodes ([]string): List of nodes for which the allocation needs to be fetched
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi CatAllocationRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func CatAllocation(ctx context.Context, nodes []string) (*osapi.Response, error) {
	return doRetry(ctx, osapi.CatAllocationRequest{
		Pretty: true,
		NodeID: nodes,
	})
}

// Input:
//
//	jsonQuery ([]byte): The json query in bytes that needs to be queried from the index
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi SearchRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func SearchQuery(ctx context.Context, jsonQuery []byte) (*osapi.Response, error) {
	return doRetry(ctx, osapi.SearchRequest{
		Index: []string{IndexName},
		Body:  bytes.NewReader(jsonQuery),
	})
}

// Input:
//
//	docId (string): The _id of the document which needs to be searched
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi GetRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func SearchDoc(ctx context.Context, docId string) (*osapi.Response, error) {
	return doRetry(ctx, osapi.GetRequest{
		Index:      IndexName,
		DocumentID: docId,
	})
}

// Input:
//
//	docId (string): The _id of the document which needs to be updated
//	content (string): The body of the request that needs to be updated in the document
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi IndexRequest along with document ID to update and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func UpdateDoc(ctx context.Context, docId string, content string) (*osapi.Response, error) {
	return doRetry(ctx, osapi.IndexRequest{
		Index:      IndexName,
		DocumentID: docId,
		Body:       strings.NewReader(content),
		Refresh:    "wait_for",
	})
}

// Input:
//
//	jsonQuery ([]byte): Query by which the deletion of documents is carried
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi DeleteByQueryRequest and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func DeleteWithQuery(ctx context.Context, jsonQuery []byte) (*osapi.Response, error) {
	return doRetry(ctx, osapi.DeleteByQueryRequest{
		Index: []string{IndexName},
		Body:  bytes.NewReader(jsonQuery),
	})
}

// Input:
//
//	ctx (context.Context): Request-scoped data that transits processes and APIs.
//
// Description:
//
//	Calls the osapi ClusterRerouteRequest with true and returns the response
//
// Return:
//
//	(*osapi.Response, error): Returns the api response and error if any
func RerouteRetryFailed(ctx context.Context) (*osapi.Response, error) {
	retry := true
	return doRetry(ctx, osapi.ClusterRerouteRequest{
		RetryFailed: &retry,
	})
}

func doRetry(ctx context.Context, request interface{}) (*osapi.Response, error) {
	var resp *osapi.Response
	var err error
	for i := 0; i < 30; i++ {
		switch request.(type) {
		case osapi.IndexRequest:
			resp, err = request.(osapi.IndexRequest).Do(ctx, osClient)
		case osapi.ClusterStatsRequest:
			resp, err = request.(osapi.ClusterStatsRequest).Do(ctx, osClient)
		case osapi.ClusterHealthRequest:
			resp, err = request.(osapi.ClusterHealthRequest).Do(ctx, osClient)
		case osapi.ClusterStateRequest:
			resp, err = request.(osapi.ClusterStateRequest).Do(ctx, osClient)
		case osapi.NodesStatsRequest:
			resp, err = request.(osapi.NodesStatsRequest).Do(ctx, osClient)
		case osapi.CatAllocationRequest:
			resp, err = request.(osapi.CatAllocationRequest).Do(ctx, osClient)
		case osapi.SearchRequest:
			resp, err = request.(osapi.SearchRequest).Do(ctx, osClient)
		case osapi.GetRequest:
			resp, err = request.(osapi.GetRequest).Do(ctx, osClient)
		case osapi.DeleteByQueryRequest:
			resp, err = request.(osapi.DeleteByQueryRequest).Do(ctx, osClient)
		case osapi.ClusterRerouteRequest:
			resp, err = request.(osapi.ClusterRerouteRequest).Do(ctx, osClient)
		}
		if err == nil && resp.StatusCode < 400 {
			break
		}
		log.Error.Println("Retrying #", i+1, reflect.TypeOf(request), " the OS API due to the error: ", resp.Status(), err)
		time.Sleep(time.Duration(10) * time.Second)
	}
	return resp, err
}
