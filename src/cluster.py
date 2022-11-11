from node import Node


class Cluster:
    """
    Acts as an interface for simulation of all associated nodes
    Performs and simulates the output of all operations performed
    by the master node.
    """
    def __init__(
            self,
            name: str,
            host: str,
            status: str,
            nodes: list[Node],
            active_data_nodes: int = 0,
            master_eligible_nodes: int = 0,
            initializing_shards: int = 0,
            relocating_shards: int = 0,
            unassigned_shards: int = 0,
            active_shards: int = 0,
            active_primary_shards: int = 0,
            active_shards_percent: float = 0

    ):
        """
        Initialize the cluster object
        :param name: name of the cluster
        :param host: ip address of the cluster
        :param status: status of the cluster from "green", "yellow" or "red"
        :param nodes: list of associated node objects
        :param active_data_nodes: count of data nodes, 0 in case of new cluster,
            original value calculated form associated nodes after assignment
        :param initializing_shards: count of shards in initializing state,
            0 in case of new cluster
        :param relocating_shards: count of shards in relocating state,
            0 in case of new cluster
        :param unassigned_shards: count of shards in unassigned state,
            0 in case of new cluster
        :param active_shards: count of shards in active state,
            0 in case of new cluster
        :param active_primary_shards: count of primary shards in active state,
            0 in case of new cluster
        :param active_shards_percent:
        """
        self.name = name
        self.host = host
        self.status = status
        self.nodes = nodes
        self.master_eligible_nodes = master_eligible_nodes
        self.initializing_shards = initializing_shards
        self.relocating_shards = relocating_shards
        self.unassigned_shards = unassigned_shards
        self.active_shards = active_shards
        self.active_data_nodes = active_data_nodes
        self.active_primary_shards = active_primary_shards
        self.active_shards_percent = active_shards_percent

    # TODO: Define methods for controlling cluster behaviour,
    #  node addition, removal etc
    # def add_node(self):
    #     pass
    #
    # def remove_node(self):
    #     pass
