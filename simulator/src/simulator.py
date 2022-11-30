import copy

from cluster import Cluster
from data_ingestion import DataIngestion
from search import Search


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
            cluster: Cluster,
            data_ingestion: DataIngestion,
            searches: list[Search],
            frequency_minutes: int,
            elapsed_time_minutes: int = 0,
    ):
        """
        Initialize the Simulator object
        :param elapsed_time_minutes: if provided, the cluster
            represents the state after the elapsed minutes
        """

        self.cluster = cluster
        self.data_ingestion = data_ingestion
        self.searches = searches
        self.elapsed_time_minutes = elapsed_time_minutes
        self.frequency_minutes = frequency_minutes

    def aggregate_data(
            self,
            duration_minutes,
            start_time_hh_mm_ss: str = '00_00_00'
    ):
        # first collect all data aggregation events
        x, y = self.data_ingestion.aggregate_data(start_time_hh_mm_ss, duration_minutes, self.frequency_minutes)
        return x, y

    def cpu_used_for_ingestion(self, ingestion):
        return min(ingestion/self.cluster.total_nodes_count * 0.075 * 100, 100)

    def memory_used_for_ingestion(self, ingestion):
        return min(ingestion / self.cluster.total_nodes_count * 0.012 * 100, 100)

    def cluster_state_for_ingestion(self):
        pass

    def run(self, duration_minutes):
        resultant_cluster_objects = []
        data_x, data_y = self.aggregate_data(duration_minutes)
        for y in data_y:
            self.cluster.cpu_usage_percent = self.cpu_used_for_ingestion(y)
            self.cluster.memory_usage_percent = self.memory_used_for_ingestion(y)
            # Todo: simulate effect on remaining cluster parameters 
            resultant_cluster_objects.append(copy.deepcopy(self.cluster))
        return resultant_cluster_objects
