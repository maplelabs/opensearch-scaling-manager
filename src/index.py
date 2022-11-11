from shard import Shard


class Index:
    """
    Essentially a collection of shards along with other
    parameters and methods that transform and store the
    state of shards associated to this index.
    """

    def __init__(
            self,
            shards: list[Shard],
            max_shards_per_index: int,
            roll_over_size: int,
            clean_up_age_in_minutes: int,
            health: str,
            store_size: int = 0,
            age_in_minutes: int = 0
    ):
        """
        Initialize the index object
        :param shards: list of associated shard objects
        :param max_shards_per_index: maximum number of shards
        :param roll_over_size: size in bytes after which the index will be rolled over
        :param clean_up_age_in_minutes: size in bytes after which the index will be rolled over
        :param health: health of the index in "green", "yellow" or "red"
        :param store_size: size of the index in bytes, 0 in case of new index
        :param age_in_minutes: age of index, 0 in case of new index
        """
        self.shards = shards
        self.max_shards_per_index = max_shards_per_index
        self.store_size = store_size
        self.health = health
        self.age = age_in_minutes
        self.roll_over_size = roll_over_size
        self.clean_up_age = clean_up_age_in_minutes
        self.host = None

    @property
    def max_shard_size(self):
        """
        Returns the maximum size a shard can have, based on
        shards per and index and roll over size
        :return: shard size in bytes
        """
        # Todo: establish relation between shard size,
        #  shards per and index and roll over size
        return

    @property
    def host(self):
        return self.host

    @host.setter
    def host(self, value):
        self._host = value
