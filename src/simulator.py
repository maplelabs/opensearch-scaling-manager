from event import Event


class Stat:
    """
    Base class for defining Stats
    """
    def __init__(
            self,
            tag: str,
            timestamp: str
    ):
        """
        Initialize the stat object
        :param tag: custom tag provided to the stat
        :param timestamp: time of stat creation in ISO 8601 format
        """
        self.tag = tag
        self.timestamp = timestamp


class NodeStat(Stat):
    """
    Class for storing statistics of a node
    """
    def __init__(
            self,
            node_id: str,
            name: str,
            host: str,
            roles: list,
            cpu_util: float,
            ram_util: float,
            heap_util: float,
            disk_util: float,
            shard_count: int,
            timestamp: str,
    ):
        """
        Initialize the NodeStat object
        :param node_id: node identifier
        :param name: name of the node
        :param host: ip address of the node
        :param roles: roles associated to the node
        :param cpu_util: cpu utilization of node in percentage
        :param ram_util: ram utilization of node in percentage
        :param heap_util: heap utilization of node in percentage
        :param disk_util: disk utilization of node in percentage
        :param shard_count: count of shards on the node
        :param timestamp: time of stat creation in ISO 8601 format
        """
        super().__init__('NodeStats', timestamp)
        self.node_id = node_id
        self.name = name
        self.host = host
        self.roles = roles
        self.cpu_util = cpu_util
        self.ram_util = ram_util
        self.heap_util = heap_util
        self.disk_util = disk_util
        self.shard_count = shard_count


class ClusterStat(Stat):
    def __init__(
            self,
            name: str,
            node_count: int,
            master_nodes: int,
            active_data_nodes: int,
            cluster_state: str,
            initializing_shards:  int,
            unassigned_shards: int,
            relocating_shards: int,
            timestamp: str
    ):
        """
        Initialize the ClusterStat object
        :param name: name of the cluster
        :param node_count: count of nodes in the cluster
        :param master_nodes: count of master nodes in the cluster
        :param active_data_nodes: count of data nodes in the cluster
        :param cluster_state: state of the cluster, "green", "yellow" or "red"
        :param initializing_shards: number of shards in the initializing state
        :param unassigned_shards: number of shards in the unassigned state
        :param relocating_shards: number of shards in the relocating state
        :param timestamp: time of stat creation in ISO 8601 format
        """
        super().__init__('ClusterStats', timestamp)
        self.name = name
        self.node_count = node_count
        self.master_nodes = master_nodes
        self.active_data_nodes = active_data_nodes
        self.cluster_state = cluster_state
        self.initializing_shards = initializing_shards
        self.unassigned_shards = unassigned_shards
        self.relocating_shards = relocating_shards


class Simulator:
    """
    Takes care of:
        - creation of cluster
        - triggering of events
        - altering the states of nodes and cluster based on events
        - providing the current statistics of node and cluster
    """
    def __init__(
            self,
            events: list[Event],
            node_stat: NodeStat,
            cluster_stat: ClusterStat,
            elapsed_time_minutes: int = 0,
    ):
        """
        Initialize the Simulator object
        :param events: list of event objects
        :param node_stat: statistics of the node as object
        :param cluster_stat: statistics of the cluster as object
        :param elapsed_time_minutes: if provided, the cluster
            represents the state after the elapsed minutes
        """
        self.events = events
        self.elapsed_time_minutes = elapsed_time_minutes
        self.node_stat = node_stat
        self.cluster_stat = cluster_stat

    def simulate_data(self, duration):
        pass
