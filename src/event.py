from pattern import Pattern


class Event:
    """
    Base class for event generation
    """
    def __init__(
            self,
            name: str,
    ):
        """
        initializes the event object
        :param name: identifier of the event
        """
        self.name = name


class CertainEvent(Event):
    """
    Certain events are the events that occur
    at the particular time of any day
    """
    def __init__(self, name: str, time_of_day: str):
        """
        initializes the certain event object
        :param name: identifier of the event
        :param time_of_day: time of the day when the event is supposed to occur
        """
        super().__init__(name)
        self.time_of_day = time_of_day

    @property
    def time_until_event(self):
        """
        tells the time before next occurrence of the event
        :return: time in seconds
        """
        # Todo: develop the logic
        return


class ProbableEvent(Event):
    """
    Events that are governed by probability
    """
    def __init__(self, name: str, probability: float):
        """
        initializes the probable event object
        :param name: identifier of the event
        :param probability: probability of occurrence between 0 and 1
        """
        super().__init__(name)
        self.probability = probability


class DataAggregation(CertainEvent):
    """
    Data aggregation event bound to happen at a fixed period of time
    """
    def __init__(
            self,
            name: str,
            time_of_day: str,
            ingestion_rate_gb_per_hour: float,
            pattern: Pattern
    ):
        """
        initializes the data aggregation object
        :param name: identifier of the event
        :param time_of_day: time of the day when the event is supposed to occur
        :param ingestion_rate_gb_per_hour: data ingestion rate
        :param pattern: pattern object to govern the event
        """
        super().__init__(name, time_of_day)
        self.ingestion_rate_gb_per_hour = ingestion_rate_gb_per_hour
        self.pattern = pattern


class SearchEvent(ProbableEvent):
    """
    Search event based on probability
    """
    def __init__(
            self,
            name: str,
            probability: float,
            cpu_load_percent: int,
            memory_load_percent: int):
        """
        initializes the search event object
        :param name: identifier of the event
        :param probability: probability of occurrence between 0 and 1
        :param cpu_load_percent: effect of event on cluster cpu
        :param memory_load_percent: effect of event on cluster memory
        """
        super().__init__(name, probability)
        self.cpu_load_percent = cpu_load_percent
        self.memory_load_percent = memory_load_percent


class SimpleSearch(SearchEvent):
    def __init__(
            self,
            name: str,
            probability: float,
            cpu_load_percent: int = 10,
            memory_load_percent: int = 10):
        """
        initializes the simple search event object
        :param name: identifier of the event
        :param probability: probability of occurrence between 0 and 1
        :param cpu_load_percent: effect of event on cluster cpu
        :param memory_load_percent: effect of event on cluster memory
        """
        super().__init__(name, probability, cpu_load_percent, memory_load_percent)


class MediumSearch(SearchEvent):
    def __init__(
            self,
            name: str,
            probability: float,
            cpu_load_percent: int = 20,
            memory_load_percent: int = 20):
        """
        initializes the medium search event object
        :param name: identifier of the event
        :param probability: probability of occurrence between 0 and 1
        :param cpu_load_percent: effect of event on cluster cpu
        :param memory_load_percent: effect of event on cluster memory
        """
        super().__init__(name, probability, cpu_load_percent, memory_load_percent)


class ComplexSearch(SearchEvent):
    def __init__(
            self,
            name: str,
            probability: float,
            cpu_load_percent: int = 30,
            memory_load_percent: int = 30):
        """
        initializes the complex search event object
        :param name: identifier of the event
        :param probability: probability of occurrence between 0 and 1
        :param cpu_load_percent: effect of event on cluster cpu
        :param memory_load_percent: effect of event on cluster memory
        """
        super().__init__(name, probability, cpu_load_percent, memory_load_percent)


class IndexRollOver(Event):
    # Todo: Define new type of event
    pass
