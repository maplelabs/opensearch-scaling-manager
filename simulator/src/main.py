# This is the main module of OpenSearch Cluster simulator
# The intent of this module is to simulate the behaviour of a cluster
# under varies loading conditions


from config_parser import parse_config
from simulator import Simulator


if __name__ == '__main__':
    stats, events = parse_config('config.json')
    simulator = Simulator(events)
    data_points = simulator.aggregate_data(24*60)
    print(data_points)
