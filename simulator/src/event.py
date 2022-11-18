import numpy as np
from scipy.interpolate import InterpolatedUnivariateSpline
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


class DataAggregation(Event):
    """
    Data aggregation event bound to happen at a fixed period of time
    """
    def __init__(
            self,
            name: str,
            patterns: [Pattern]
    ):
        """
        initializes the data aggregation object
        :param name: identifier of the event
        :param ingestion_rate_gb_per_hour: data ingestion rate
        :param patterns: pattern object to govern the event
        """
        super().__init__(name)
        self.patterns = patterns

    def data_aggregation_points(
            self,
            start_time_hh_mm_ss: str,
            duration_minutes: int,
            frequency_minutes: int):
        """
        Produce cumulative data points of all events and return an array of resultant aggregation
        :param start_time_hh_mm_ss: start time in hh_mm_ss in 24-hour format, eg. '080000'
        :param duration_minutes: duration of point generation in minutes
        :param frequency_minutes: gap between the resultant points
        :return: array of float containing resultant data aggregation points
        """

        # fits
        if all(pattern.name == 'fixed' for pattern in self.patterns):
            print('aggregating fixed patterns')
            time_of_day, ingestion_rate_gb_per_hour = [], []

            for pat in self.patterns:
                time_of_day.append(pat.time_of_day_hh_mm_ss)
                ingestion_rate_gb_per_hour.append(pat.max)

            font1 = {'family': 'serif', 'color': 'red', 'size': 15}
            font2 = {'family': 'serif', 'color': 'darkred', 'size': 10}

            time_of_day = [int(i.split('_')[0]) * 60 for i in time_of_day]

            print(ingestion_rate_gb_per_hour)
            print(time_of_day)

            # plt.subplot(3, 1, 1)
            # # Plot the user defined inflection points first
            # plt.plot(time_of_day, ingestion_rate_gb_per_hour, 'o')
            # plt.xlabel('Time (in minutes)', font2)
            # # plt.ylabel('Ingestion Rate (in GB/Hr)', font2)
            # plt.title('User Defined Inflection Points', fontdict=font1)
            # plt.xlim([0, duration_minutes])  # set range for x axis
            # plt.ylim([0, 60])  # set range for x axis
            # plt.grid()

            # add missing value of 0th hour
            if time_of_day[0] != 0:
                time_of_day.insert(0, 0)
                ingestion_rate_gb_per_hour.insert(0, 5)

            # positions to inter/extrapolate
            intervals = int(duration_minutes / frequency_minutes) + 1

            x = np.linspace(0, duration_minutes, intervals)
            # spline order: 1 linear, 2 quadratic, 3 cubic ... 
            order = 1
            # do inter/extrapolation
            s = InterpolatedUnivariateSpline(time_of_day, ingestion_rate_gb_per_hour, k=order)
            y = s(x)

            print(y)
            print(x)
            return x, y

        if all(pattern.name == 'random' for pattern in self.patterns):
            print('aggregating random patterns')

            # positions to inter/extrapolate
            intervals = int(duration_minutes / frequency_minutes) + 1
            random_set = np.random.randint(0, 5, size=intervals)
            x = np.linspace(0, duration_minutes, intervals)
            return x, random_set

            # plt.subplot(3, 1, 2)
            # # Plot the random load
            # plt.plot(x, random_set)
            # plt.xlim([0, duration_minutes])  # set range for x axis
            # plt.ylim([0, 60])  # set range for x axis
            # plt.xlabel('Time (in minutes) -->', font2)
            # plt.ylabel('Ingestion Rate (in GB/Hr) -->', font2)
            # plt.title('Random Load', font1)
            # plt.grid()

            # plt.subplot(3, 1, 3)
            # Plot the resultant wave form
            # plt.plot(x, y + random_set)
            # plt.xlim([0, duration_minutes])  # set range for x axis
            # plt.ylim([0, 60])  # set range for x axis
            # plt.xlabel('Time (in minutes) -->', font2)
            # # plt.ylabel('Ingestion Rate (in GB/hr)', font2)
            # plt.title('Resultant Load', font1)
            # plt.grid()
            #
            # plt.show()


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
