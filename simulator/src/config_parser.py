import yaml


class Config:
    def __init__(
            self,
            stats: dict,
            data_ingestion: dict,
            searches: dict,
            data_generation_interval_minutes: int
    ):
        self.stats = stats
        self.data_generation_interval_minutes = data_generation_interval_minutes
        self.data_ingestion = data_ingestion
        self.searches = searches


def parse_config(config_file_path):
    """
        Read and parse the config file into objects,
        that can work with simulator
        :param config_file_path: path of the yaml file
        :return: stats, events
    """
    # read the config file
    file_object = open(config_file_path)
    all_configs = yaml.safe_load(file_object)
    file_object.close()
    data_generation_interval_minutes = all_configs.pop('data_generation_interval_minutes')
    data_ingestion = all_configs.pop('data_ingestion')
    searches = all_configs.pop('searches')
    stats = all_configs

    # Todo: add validations
    return Config(stats, data_ingestion, searches, data_generation_interval_minutes)


if __name__ == '__main__':
    parse_config('config.yaml')
