# This is the main module of OpenSearch Cluster simulator
# The intent of this module is to simulate the behaviour of a cluster
# under varies loading conditions

import sys
from config_parser import parse_config
from simulator import Simulator
from cluster import Cluster
from data_ingestion import State, DataIngestion


if __name__ == '__main__':
    print(sys.argv)
    input('...')
    configs = parse_config('config.yaml')
    all_states = [State(**state) for state in configs.data_ingestion.get('states')]
    randomness_percentage = configs.data_ingestion.get('randomness_percentage')

    data_function = DataIngestion(all_states, randomness_percentage)

    cluster = Cluster(**configs.stats)

    sim = Simulator(cluster, data_function, configs.searches, configs.data_generation_interval_minutes, 0)

    result = sim.run(24*60)
    for res in result:
        print(res.cpu_usage_percent, res.memory_usage_percent)
