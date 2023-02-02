# Open-search Scaling Manager

Open Search Simulator is an attempt to mimic to behavior of an AWS on which OpenSearch is deployed.



### Simulator Configurations

------

The user can specify some key features of an OpenSearch Cluster for simulator through the config file provided. The functionalities supported are:



#### 1.Cluster Stats

------

**cluster_name:** Name of the cluster that is to be used.

**cluster_hostname:** Host name of the cluster.

**cluster_ip_address:** IP address of the cluster.

**node_machine_type_identifier:** Defines the type of the instance or node deployed in a cluster.

**total_nodes_count:** Total number of nodes present in the cluster.

**active_data_nodes:** Number of active data nodes in total number of nodes present in cluster.

**min_nodes_in_cluster:** Minimum number of nodes that the cluster must have to perform the necessary tasks.

**master_eligible_nodes_count:** Nodes that are eligible to become master whenever the present master node goes down.  

**heap_memory_factor:**

**index_count:** Number of index that cluster must have.

**primary_shards_per_index:** Number of primary shards that is present in index.

**replica_shards_per_index:** Number of replica shards that is present in index(replica of data that represents each primary shard).

**index_roll_over_size_gb:** Specific size at where index will roll over to new index when it exceeds.

**index_clean_up_age_days:** Duration at which index cleanup happens.

**total_disk_size_gb:** Total number of size in GB that the disk should have.

**simulation_frequency_minutes:** Time interval that the simulator will run the data simulation.



#### 2.Data Ingestion

------

Specify data ingestion with respect to time of the day to represent pattern for entire day(24hrs).

**states:** Field that contains various positions that contains the data require for the simulator to run in particular time interval.

**day:** Field that contains multiple day positions.

**position:** For a day there can be any number of position where it contains time_hh_mm_ss, ingestion_rate_gb_per_hr, searches.

**time_hh_mm_ss:** Time interval of the position. 

**ingestion_rate_gb_per_hr:** Amount of data that has been ingested for the particular interval of time that is defined in time_hh_mm_ss.

**searches:** Contains the types of searches that needs to be made, if the config has certain searches it takes the corresponding values. Three types of searches are simple, medium, complex.



#### 3.Randomness Percentage

------

**randomness_percentage:**  Percentage at which the stats value needs to be differing while simulating.



#### 4.Search Description

------

**search_description:** Specify searches along with their type, probability and load inflected on the cluster. Three level of search_description are simple, medium, complex.

**simple:**

​	**cpu_load_percent:** Percentage at which cpu must be used if search_description is simple.

​	**memory_load_percent: **Percentage at which memory must be used if search_description is simple.

​	**heap_load_percent: **Percentage at which heap must be used  if search_description is simple.

**medium:**

​	**cpu_load_percent:** Percentage at which cpu must be used if search_description is medium.

​	**memory_load_percent: **Percentage at which memory must be used if search_description is medium.

​	**heap_load_percent: **Percentage at which heap must be used if search_description is medium.

**complex:**

​	**cpu_load_percent:** Percentage at which cpu must be used if search_description is complex.

​	**memory_load_percent: **Percentage at which memory must be used if search_description is complex.

​	**heap_load_percent: **Percentage at which heap must be used if search_description is complex.



### Sample cofig.yaml

------

```yaml
---
cluster_name: test
cluster_hostname: test-cluster.aws
cluster_ip_address: 10.0.0.1
node_machine_type_identifier: m5-12xlarge
total_nodes_count: 7
active_data_nodes: 7
min_nodes_in_cluster: 3
master_eligible_nodes_count: 7
heap_memory_factor: 0.5
index_count: 100
primary_shards_per_index: 2
replica_shards_per_index: 1
index_roll_over_size_gb: 10
index_clean_up_age_days: 20
total_disk_size_gb: 14000
simulation_frequency_minutes: 5

# Specify data ingestion with respect to time of the day to represent pattern for entire day(24hrs).
states:
- Day: 1
  pattern:
    - position: 1
      time_hh_mm_ss: '00_00_00'
      ingestion_rate_gb_per_hr: 0
      searches:
        simple: 50000
        medium: 500
    - position: 2
      time_hh_mm_ss: '02_00_00'
      ingestion_rate_gb_per_hr: 0
      searches:
        simple: 30000
        medium: 1000
    - position: 3
      time_hh_mm_ss: '04_00_00'
      ingestion_rate_gb_per_hr: 0
      searches:
        simple: 50000
        medium: 2000
    - position: 4
      time_hh_mm_ss: '06_00_00'
      ingestion_rate_gb_per_hr: 3
      searches:
        simple: 50000
        medium: 2000
    - position: 5
      time_hh_mm_ss: '08_00_00'
      ingestion_rate_gb_per_hr: 70
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 6
      time_hh_mm_ss: '09_00_00'
      ingestion_rate_gb_per_hr: 60
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 7
      time_hh_mm_ss: '10_00_00'
      ingestion_rate_gb_per_hr: 80
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 8
      time_hh_mm_ss: '11_00_00'
      ingestion_rate_gb_per_hr: 24
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 9
      time_hh_mm_ss: '12_00_00'
      ingestion_rate_gb_per_hr: 50
      searches:
        simple: 110000
        medium: 80000
        complex: 55000
    - position: 10
      time_hh_mm_ss: '13_00_00'
      ingestion_rate_gb_per_hr: 20
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 11
      time_hh_mm_ss: '14_00_00'
      ingestion_rate_gb_per_hr: 73
      searches:
        simple: 60000
        medium: 50000
        complex: 25000
    - position: 12
      time_hh_mm_ss: '15_00_00'
      ingestion_rate_gb_per_hr: 60
      searches:
        simple: 30000
        medium: 50000
        complex: 10000
    - position: 13
      time_hh_mm_ss: '16_00_00'
      ingestion_rate_gb_per_hr: 90
      searches:
        simple: 55000
        medium: 45000
        complex: 20000
    - position: 14
      time_hh_mm_ss: '17_00_00'
      ingestion_rate_gb_per_hr: 56
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 15
      time_hh_mm_ss: '18_00_00'
      ingestion_rate_gb_per_hr: 70
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 16
      time_hh_mm_ss: '19_00_00'
      ingestion_rate_gb_per_hr: 40
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 17
      time_hh_mm_ss: '20_00_00'
      ingestion_rate_gb_per_hr: 26
      searches:
        simple: 20000
        medium: 10000
    - position: 18
      time_hh_mm_ss: '21_00_00'
      ingestion_rate_gb_per_hr: 10
      searches:
        simple: 20000
        medium: 10000
    - position: 19
      time_hh_mm_ss: '22_00_00'
      ingestion_rate_gb_per_hr: 7
      searches:
        simple: 60000
        medium: 10000
    - position: 20
      time_hh_mm_ss: '23_00_00'
      ingestion_rate_gb_per_hr: 1
      searches:
        simple: 10000
        medium: 2000      
- Day: 2
  pattern:
    - position: 1
      time_hh_mm_ss: '00_00_00'
      ingestion_rate_gb_per_hr: 0
      searches:
        simple: 50000
        medium: 500
    - position: 2
      time_hh_mm_ss: '02_00_00'
      ingestion_rate_gb_per_hr: 0
      searches:
        simple: 30000
        medium: 1000
    - position: 3
      time_hh_mm_ss: '04_00_00'
      ingestion_rate_gb_per_hr: 0
      searches:
        simple: 50000
        medium: 2000
    - position: 4
      time_hh_mm_ss: '06_00_00'
      ingestion_rate_gb_per_hr: 3
      searches:
        simple: 50000
        medium: 2000
    - position: 5
      time_hh_mm_ss: '08_00_00'
      ingestion_rate_gb_per_hr: 70
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 6
      time_hh_mm_ss: '09_00_00'
      ingestion_rate_gb_per_hr: 60
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 7
      time_hh_mm_ss: '10_00_00'
      ingestion_rate_gb_per_hr: 80
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 8
      time_hh_mm_ss: '11_00_00'
      ingestion_rate_gb_per_hr: 24
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 9
      time_hh_mm_ss: '12_00_00'
      ingestion_rate_gb_per_hr: 50
      searches:
        simple: 110000
        medium: 80000
        complex: 55000
    - position: 10
      time_hh_mm_ss: '13_00_00'
      ingestion_rate_gb_per_hr: 20
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 11
      time_hh_mm_ss: '14_00_00'
      ingestion_rate_gb_per_hr: 73
      searches:
        simple: 60000
        medium: 50000
        complex: 25000
    - position: 12
      time_hh_mm_ss: '15_00_00'
      ingestion_rate_gb_per_hr: 60
      searches:
        simple: 30000
        medium: 50000
        complex: 10000
    - position: 13
      time_hh_mm_ss: '16_00_00'
      ingestion_rate_gb_per_hr: 90
      searches:
        simple: 55000
        medium: 45000
        complex: 20000
    - position: 14
      time_hh_mm_ss: '17_00_00'
      ingestion_rate_gb_per_hr: 56
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 15
      time_hh_mm_ss: '18_00_00'
      ingestion_rate_gb_per_hr: 70
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 16
      time_hh_mm_ss: '19_00_00'
      ingestion_rate_gb_per_hr: 40
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 17
      time_hh_mm_ss: '20_00_00'
      ingestion_rate_gb_per_hr: 26
      searches:
        simple: 20000
        medium: 10000
    - position: 18
      time_hh_mm_ss: '21_00_00'
      ingestion_rate_gb_per_hr: 10
      searches:
        simple: 20000
        medium: 10000
    - position: 19
      time_hh_mm_ss: '22_00_00'
      ingestion_rate_gb_per_hr: 7
      searches:
        simple: 60000
        medium: 10000
    - position: 20
      time_hh_mm_ss: '23_00_00'
      ingestion_rate_gb_per_hr: 1
      searches:
        simple: 10000
        medium: 2000
- Day: 3
  pattern:
    - position: 1
      time_hh_mm_ss: '00_00_00'
      ingestion_rate_gb_per_hr: 0
      searches:
        simple: 50000
        medium: 500
    - position: 2
      time_hh_mm_ss: '02_00_00'
      ingestion_rate_gb_per_hr: 0
      searches:
        simple: 30000
        medium: 1000
    - position: 3
      time_hh_mm_ss: '04_00_00'
      ingestion_rate_gb_per_hr: 0
      searches:
        simple: 50000
        medium: 2000
    - position: 4
      time_hh_mm_ss: '06_00_00'
      ingestion_rate_gb_per_hr: 3
      searches:
        simple: 50000
        medium: 2000
    - position: 5
      time_hh_mm_ss: '08_00_00'
      ingestion_rate_gb_per_hr: 90
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 6
      time_hh_mm_ss: '09_00_00'
      ingestion_rate_gb_per_hr: 12
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 7
      time_hh_mm_ss: '10_00_00'
      ingestion_rate_gb_per_hr: 80
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 8
      time_hh_mm_ss: '11_00_00'
      ingestion_rate_gb_per_hr: 124
      searches:
        simple: 100000
        medium: 80000
        complex: 50000
    - position: 9
      time_hh_mm_ss: '12_00_00'
      ingestion_rate_gb_per_hr: 10
      searches:
        simple: 110000
        medium: 80000
        complex: 55000
    - position: 10
      time_hh_mm_ss: '13_00_00'
      ingestion_rate_gb_per_hr: 90
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 11
      time_hh_mm_ss: '14_00_00'
      ingestion_rate_gb_per_hr: 73
      searches:
        simple: 60000
        medium: 50000
        complex: 25000
    - position: 12
      time_hh_mm_ss: '15_00_00'
      ingestion_rate_gb_per_hr: 60
      searches:
        simple: 30000
        medium: 50000
        complex: 10000
    - position: 13
      time_hh_mm_ss: '16_00_00'
      ingestion_rate_gb_per_hr: 90
      searches:
        simple: 55000
        medium: 45000
        complex: 20000
    - position: 14
      time_hh_mm_ss: '17_00_00'
      ingestion_rate_gb_per_hr: 56
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 15
      time_hh_mm_ss: '18_00_00'
      ingestion_rate_gb_per_hr: 70
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 16
      time_hh_mm_ss: '19_00_00'
      ingestion_rate_gb_per_hr: 40
      searches:
        simple: 55000
        medium: 40000
        complex: 20000
    - position: 17
      time_hh_mm_ss: '20_00_00'
      ingestion_rate_gb_per_hr: 26
      searches:
        simple: 20000
        medium: 10000
    - position: 18
      time_hh_mm_ss: '21_00_00'
      ingestion_rate_gb_per_hr: 10
      searches:
        simple: 20000
        medium: 10000
    - position: 19
      time_hh_mm_ss: '22_00_00'
      ingestion_rate_gb_per_hr: 7
      searches:
        simple: 60000
        medium: 10000
    - position: 20
      time_hh_mm_ss: '23_00_00'
      ingestion_rate_gb_per_hr: 1
      searches:
        simple: 10000
        medium: 2000
randomness_percentage: 35


# Specify searches along with their type, probability and load inflected on the cluster.
search_description:
  simple:
    cpu_load_percent: 0.001
    memory_load_percent: 1
    heap_load_percent: 0.01
  medium:
    cpu_load_percent: 0.0015
    memory_load_percent: 1.5
    heap_load_percent: 0.01
  complex:
    cpu_load_percent: 0.002
    memory_load_percent: 2
    heap_load_percent: 0.01

```



### Simulator Behavior

------

As simulator starts, it generates and stores the data points corresponding to the entire day and stores them in a internal database. Based on the user inputs (through APIs), the data points are fetched or re-generated.



### Installation and Executing Simulator

------

To install the simulator please download the source code using following command:

```
git clone https://github.com/maplelabs/opensearch-scaling-manager.git -b release_v0.1_dev
```



Execute the following commands to run and install the simulator

```python
cd opensearch-scaling-manager/simulator
# Path to simulator module.

python -m venv venv
# Creating virtual environment.

.\venv\Scripts\activate
# Activatinng virtual environment.

pip install -r .\requirements.txt
# Install every requirements for simulator.

cd src
# Path to execute simulator.

python app.py
# Run entire simulator module.
```



### APIs

------

Simulator provide the following APIs to interact with it

| Path               | Query Parameters                                             | Description                                                  | Method | Request Body       | Response                                   |
| :----------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------ | ------------------ | ------------------------------------------ |
| /stats/avg         | {key,value} = {metric:string},{duration:int}                 | Returns the average value of a stat for the last specified duration. | GET    | None               | {"avg": float, "min": float, "max": float} |
| /stats/violated    | {key,value} = {metric:string},{duration:int},{threshold:float} | Returns the number of time, a stat crossed the threshold duration the specified duration. | GET    | None               | {"ViolatedCount": int}                     |
| /stats/current     | {key,value} = {metric:string},{duration:int}                 | Returns the most recent value of a stat.                     | GET    | None               | {"current": float}                         |
| /provision/addnode | None                                                         | Ask the simulator to perform a node addition.                | POST   | {"nodes": integer} | {"nodes": int}                             |
| /provision/remnode | None                                                         | Ask the simulator to perform a node removal.                 | POST   | {"nodes": integer} | {"nodes": int}                             |



## Scaling Manager Configuration

------

The user can specify some key features of an OpenSearch Cluster for simulator through the config file provided. The functionalities supported are:

**user_config:**

​	**monitor_with_logs:** Field that contains bool value which specifies whether to monitor with logs or not

​	**monitor_with_simulator:** Field that contains bool value which specifies whether to monitor with simulator or not

​	**polling_interval_in_secs:**  polling_interval_in_secs indicates the time in seconds for which polling will be repeated

​	**is_accelerated:** Field that contains bool value which accelerates the time

**cluster_details:**

​	**ip_address:** IP address of the cluster 

​	**cluster_name:** Name of the cluster 

​	**os_credentials:** 

 		**os_admin_username:** Username for admin

 		**os_admin_password:** Password for admin

 	**cloud_type:** Type of cloud used in cluster

​	 **cloud_credentials:**

​		 **secret_key:** Secret key for cluster

​		**access_key:** Access key for cluster

​	 **base_node_type:** t2x.large

​	 **number_cpus_per_node:** Total number of CPU present per node

​	 **ram_per_node_in_gb:** Size of RAM used per node (GB)

​	 **disk_per_node_in_gb:** Size of DISK used per node (GB)

 	**number_max_nodes_allowed:** Maximum number of nodes allowed for the cluster

**task_details:** Field that contains details on what task should be performed i.e scale_up_by_1 or scale_down_by_1

- **task_name:** Task name indicates the name of the task to recommend by the recommendation engine.
  **operator:** Operator indicates the logical operation needs to be performed while executing the rules
  **rules:** Rules indicates list of rules to evaluate the criteria for the recommendation engine.

  - **metric:** Metric indicates the name of the metric. These can be CpuUtil, MemUtil, ShardUtil, DiskUtil
    **limit: **Limit indicates the threshold value for a metric.
    **stat:** Stat indicates the statistics on which the evaluation of the rule will happen. These can be AVG, COUNT
    **decision_period:** Decision Period indicates the time in minutes for which a rule is evaluated.

  

### Sample config.yaml

------

```yaml
user_config:
  monitor_with_logs: true
  monitor_with_simulator: false
  polling_interval_in_secs: 10
  is_accelerated: false
cluster_details:
  ip_address: 10.81.1.250
  cluster_name: cluster.1
  os_credentials: 
    os_admin_username: elastic
    os_admin_password: changeme
  cloud_type: AWS
  cloud_credentials:
    secret_key: secret_key
    access_key: access_key
  base_node_type: t2x.large
  number_cpus_per_node: 5
  ram_per_node_in_gb: 10
  disk_per_node_in_gb: 100
  number_max_nodes_allowed: 2
task_details:
- task_name: scale_up_by_1
  operator: OR
  rules:
  - metric: CpuUtil
    limit: 55
    stat: AVG
    decision_period: 300
  - metric: CpuUtil
    limit: 50
    stat: COUNT
    occurrences: 10
    decision_period: 300
  - metric: RamUtil
    limit: 90
    stat: AVG
    decision_period: 300
- task_name: scale_down_by_1
  operator: AND
  rules:
  - metric: CpuUtil
    limit: 90
    stat: AVG
    decision_period: 300
  - metric: CpuUtil
    limit: 90
    stat: COUNT
    occurrences: 6
    decision_period: 300
  - metric: RamUtil
    limit: 99
    stat: AVG
    decision_period: 300
```



### Build, Packaging and installation

------

To install the scaling manager please download the source code using following command:

```
git clone https://github.com/maplelabs/opensearch-scaling-manager.git -b release_v0.1_dev
```



Run the following commands to build and install the scaling manager

```
cd opensearch-scaling-manager/
# Build the scaling_manager module.
sudo make build
# Package the scaling_manager module and create a tarball.
sudo make pack
# Install the scaling_manager module and create systemd service.
sudo make install
```



To start scaling manager run the following command:

```
sudo systemctl start scaling_manager
```



To stop the scaling manager run the following command:

```
sudo systemctl stop scaling_manager
```