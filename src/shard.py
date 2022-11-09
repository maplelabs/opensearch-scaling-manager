
class Shard:
    """
    Shard acts as the smallest unit of data aggregation.
    They can either be primary or replica.
    """

    def __init__(
            self,
            # initializing: int,
            # relocating: int,
            # unassigned: int,
            # active: int,
            size: int,
            type_of_shard: str,
            state: str
    ):
        # Todo: Decide initializing, relocating, unassigned and
        #  active shards are to be properties of individual shard
        #  or managed in cluster
        """
        initializes the shard object
        :param size: size of the shard in bytes
        :param type_of_shard: "primary" or "replica" shard
        :param state: Can have the following values only
            INITIALIZING: The shard is recovering from a peer shard or gateway.
            RELOCATING: The shard is relocating.
            STARTED: The shard has started.
            UNASSIGNED: The shard is not assigned to any node.
        """
        # self.initializing = initializing
        # self.relocating = relocating
        # self.unassigned = unassigned
        # self.active = active
        self.size = size
        self.type = type_of_shard
        self.state = state
        self.host = None

    # @property
    # def total(self):
    #     """
    #
    #     :return:
    #     """
    #     # Todo: Derive total number of shards from initializing,
    #     #  relocating, unassigned and active shards
    #     return self.initializing + self.unassigned + self.relocating

    @property
    def host(self):
        return self.host

    @host.setter
    def host(self, value):
        self._host = value
