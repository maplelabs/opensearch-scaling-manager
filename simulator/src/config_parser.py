import json

from cluster import Cluster
from event import DataAggregation
import pattern


class Config:
    # Todo: create config class attributes to return from config parser
    instance_types = None


def parse_config(config_file_path):
    """
        Read and parse the config file into objects,
        that can work with simulator
        :param config_file_path: path of the json file
        :return: stats, events
    """
    # Define placeholders for cluster statistics and events to return
    stats, events = {}, []
    ingestion_patterns = []
    # read the config file
    file_object = open(config_file_path)
    all_configs = json.loads(file_object.read())
    file_object.close()

    for event in all_configs.get('events'):
        if event['type'] == 'ingestion':
            # print(event)
            if event['pattern'] == 'fixed':
                for state in event['pattern_params']['states']:
                    print(state)
                    ingestion_patterns.append(
                        pattern.Fixed(state['ingestion_rate_gb_per_hr'],
                                      state['time_hh_mm_ss'])
                    )
                events.append(DataAggregation(
                    name=event['type'],
                    patterns=ingestion_patterns
                )
                )
            elif event['pattern'] == 'random':
                events.append(DataAggregation(
                    name=event['type'],
                    patterns=[pattern.Random(0, event['pattern_params']['ingestion_rate_gb_per_hr'])]
                )
                )

    return stats, events
