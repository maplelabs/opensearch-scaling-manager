class Cluster:
    """
    Acts as an interface for simulation of all associated nodes
    Performs and simulates the output of all operations performed
    by the master node.
    """
    def __init__(
            self,
            name: str,
            host_name: str,
            ip_address: str,
            status: str,
            cpu_usage_percent: float = 0,
            memory_usage_percent: float = 0,
            disk_usage_percent: float = 0,
            heap_usage_percent: float = 0,
            total_nodes_count: int = 0,
            data_nodes_count: int = 0,
            master_nodes_count: int = 0,
            index_count: int = 0,
            index_roll_over_size: int = 0,
            index_clean_up_age_in_minutes: int = 0,
            total_shard_count: int = 0,
            shards_per_index: int = 0,
            initializing_shards_count: int = 0,
            relocating_shards_count: int = 0,
            unassigned_shards_count: int = 0,
            active_shards_count: int = 0,
    ):
        """

        :param name:
        :param host_name:
        :param ip_address:
        :param status:
        :param cpu_usage_percent:
        :param memory_usage_percent:
        :param disk_usage_percent:
        :param heap_usage_percent:
        :param total_nodes_count:
        :param data_nodes_count:
        :param master_nodes_count:
        :param index_count:
        :param index_roll_over_size:
        :param index_clean_up_age_in_minutes:
        :param total_shard_count:
        :param shards_per_index:
        :param initializing_shards_count:
        :param relocating_shards_count:
        :param unassigned_shards_count:
        :param active_shards_count:
        """
        """
        Initialize the cluster object
        :param name: name of the cluster
        :param host_name: name of the cluster host
        :param ip_address: ip address of the cluster
        :param status: status of the cluster from "green", "yellow" or "red"
        :param active_data_nodes: count of data nodes, 0 in case of new cluster,
            original value calculated form associated nodes after assignment
        :param initializing_shards_count: count of shards in initializing state,
            0 in case of new cluster
        :param relocating_shards_count: count of shards in relocating state,
            0 in case of new cluster
        :param unassigned_shards_count: count of shards in unassigned state,
            0 in case of new cluster
        :param active_shards_count: count of shards in active state,
            0 in case of new cluster
        :param active_primary_shards: count of primary shards in active state,
            0 in case of new cluster
        :param active_shards_percent:
        """
        self.name = name
        self.host_name = host_name
        self.ip_address = ip_address
        self.status = status
    #    self.nodes = nodes
        self.cpu_usage_percent = cpu_usage_percent
        self.memory_usage_percent =  memory_usage_percent
        self.disk_usage_percent = disk_usage_percent
        self.heap_usage_percent = heap_usage_percent
        self.total_nodes_count = total_nodes_count
        self.data_nodes_count = data_nodes_count
        self.master_nodes_count = master_nodes_count
        self.index_count = index_count
        self.index_roll_over_size = index_roll_over_size
        self.index_clean_up_age_in_minutes = index_clean_up_age_in_minutes
        self.total_shard_count =total_shard_count
        self.shards_per_index = shards_per_index 
        self.initializing_shards = initializing_shards_count
        self.relocating_shards = relocating_shards_count
        self.unassigned_shards = unassigned_shards_count
        self.active_shards = active_shards_count
    #   self.active_primary_shards = active_primary_shards
    #   self.active_shards_percent = active_shards_percent

    # TODO: Define methods for controlling cluster behaviour,
    #  node addition, removal etc
    # def add_node(self):
    #     pass
    #
    # def remove_node(self):
    #     pass
