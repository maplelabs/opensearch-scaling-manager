// This packkage provide the data structure needed to get the metrics.
// There are two kind of metrics:
//
//	Cluster metrics: This data structure will provide cluster level metrics.
//	Node metrics: This data structure will provide node level metrics.
//
// The cluster metrics contains list of the node metrics collected over all the nodes present in a cluster.
// Also package contains a data structure called Metric which will calculate the statistics over a period of time.
// This will be used by recommendation module.
package cluster

// This struct will contain node metrics for a node in the OpenSearch cluster.
type Node struct {
	// NodeId indicates a unique ID of the node given by OpenSearch.
	NodeId string
	// NodeName indicates human-readable identifier for a particular instance of OpenSearch which is a configurable input.
	NodeName string
	// HostIp indicates the IP address of the node.
	HostIp string
	// IsMater indicates if the node is a master node.
	IsMaster bool
	// IsData indicates if the node is a data node.
	IsData bool
	// CpuUtil indicates the overall CPU Utilization in percentage for a node.
	CpuUtil float32
	// MemUtil indicates the overall Memory Utilization in percentage for a node.
	RamUtil float32
	// HeapUtil indicates the overall Java Heap Utilization in percentage for a node.
	HeapUtil float32
	// DiskUtil indicates the overall Disk Utilization in percentage for a node.
	DiskUtil float32
	// NumShards Number of shards present on a node.
	NumShards int
}

// This struct will contain the static metrics of the cluster.
type ClusterStatic struct {
	// ClusterName indicates the Cluster name for the OpenSearch cluster.
	ClusterName string
	// IpAddress indicate the master node IP for the OpenSearch cluster.
	IpAddress string
	// CloudType indicate the type of the cloud service where the OpenSearch cluster is deployed.
	CloudType string
	// BaseNodeType indicate the instance type of the node.
	// This parameters depends on the cloud service.
	BaseNodeType string
	// NumCpusPerNode indicates the number of the CPU core running on a node in a cluster.
	NumCpusPerNode string
	// RAMPerNodeInGB indicates the RAM size in GB running on a node in a cluster.
	RAMPerNodeInGB string
	// DiskPerNodeInGB indicates the Disk size in GB running on a node in a cluster.
	DiskPerNodeInGB string
	// NumMaxNodesAllowed indicates the number of maximum allowed node present in the cluster.
	// Based on this value we will determine whether to scale out further or not.
	NumMaxNodesAllowed int
}

// This struct will contain the dynamic metrics of the cluster.
type ClusterDynamic struct {
	// NumNodes indicates the number of nodes present in the OpenSearch cluster at any time.
	NumNodes int
	// ClusterStatus indicates the present state of a cluster.
	// 	red: One or more primary shards are unassigned, so some data is unavailable.
	//		This can occur briefly during cluster startup as primary shards are assigned.
	//  yellow: All primary shards are assigned, but one or more replica shards are unassigned.
	//		If a node in the cluster fails, some data could be unavailable until that node is repaired.
	//  green: All shards are assigned.
	ClusterStatus string
	// NumActiveShards indicates the total number of active primary and replica shards.
	NumActiveShards int
	// NumActivePrimaryShards indicates the number of active primary shards.
	NumActivePrimaryShards int
	// NumInitializingShards indicates the number of shards that are under initialization.
	NumInitializingShards int
	// NumUnassignedShards indicats the number of shards that are not allocated.
	NumUnassignedShards int
	// NumRelocatingShards indicates the number of shards that are under relocation.
	NumRelocatingShards int
	// NumMasterNodes indicates the number of master eligible nodes present in the cluster.
	NumMasterNodes int
	// NumActiveDataNodes indicates the number of active data nodes present in the cluster.
	NumActiveDataNodes int
}

// This struct will provide the overall cluster metrcis for a OpenSearch cluster.
type Cluster struct {
	// ClusterStatic indicates the static set of data present for a cluster.
	ClusterStatic ClusterStatic
	// ClusterDyanamic indicates the dynamic set of data present for a cluster.
	ClusterDynamic ClusterDynamic
	// NodeList indicates node metrics for all the nodes.
	NodeList []Node
}

// This struct used by the recommendation engine to find the statistics of a metrics for a given period.(CPU, MEM, HEAP, DISK).
type Metric struct {
	// Avg indicates the average for a metric for a time period.
	Avg float32
	// Min indicates the minimum value for a metric for a time period.
	Min float32
	// Max indicates the maximum value for a metric for a time period.
	Max float32
}

// This struct contains statistics for cluster and node for an evaluation period.
type MetricValues struct {
	// ClusterLevel indicates the statistical data for a metric for a time period.
	ClusterLevel Metric
	// NodeLevel indicates the list of statistical data for a metric for a node for a time period.
	NodeLevel []Metric
}

// This contains the count voilated for cluster and node for an evaluation period.
type MetricCount struct {
	// ClusterLevel indicates the count voilated for a cluster for a time period.
	ClusterLevel int
	// NodeLevel indicates the list of the count voilated for a node for a time period.
	NodeLevel []int
}
