from index import Index


class Node:
    """
    Representation of a OpenSearch node that essentially
    is a collection of node resources like shards, cpu, memory
    and the corresponding methods altering these parameters
    """
    def __init__(
            self,
            id_: str,
            name: str,
            instance_type: str,
            cpu_cores: int,
            memory: int,
            host: str,
            roles: list,
            indices: list[Index],
            cpu_usage_percent: int = 0,
            memory_usage_percent: int = 0,
            disk_usage_percent: int = 0,
            heap_usage_percentage: int = 0,
            current_ingestion_rate: int = 0
    ):
        """
        Initialize the node object
        :param id_
        :param name: name of the node
        :param instance_type: AWS instance type,
            Derive the value of cpus and memory from internal mapping
        :param cpu_cores: number of cpu cores
        :param memory: RAM of the instance in bytes
        :param host: ip address of the node
        :param roles: list of roles (in string)
        :param indices: list of associated index objects
        :param cpu_usage_percent: current cpu usage of the node
        :param memory_usage_percent: current memory usage of the node
        :param disk_usage_percent: current disk usage of the node
        :param heap_usage_percentage: current heap usage of the node
        :param current_ingestion_rate: current ingestion rate of the node
        """
        self.id = id_
        self.name = name
        self.instance_type = instance_type
        self.memory = memory
        self.cpu_cores = cpu_cores
        self.host = host
        self.roles = roles
        self.cpu_usage_percent = cpu_usage_percent
        self.memory_usage_percent = memory_usage_percent
        self.indices = indices
        self.disk_usage_percent = disk_usage_percent
        self.heap_usage_percentage = heap_usage_percentage
        self.current_ingestion_rate = current_ingestion_rate

    @property
    def is_master(self):
        """
        States whether the node is mater or not
        :return: bool
        """
        return 'master' in self.roles

    @property
    def is_data(self):
        """
        States whether the node is data or not
        :return: bool
        """
        return 'data' in self.roles



    #Todo: create method for calculating ingestion rate