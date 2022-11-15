// This package consists of all the data structure required for defining a task.
// Tasks are set of Actions.
// The actions can have list of rules.
// The recommendation engine will parse these rules and recommend the action if rules meets the criteria.
// Multiple rules can be added inside an action and like wise multiple actions can be added inside a task.
package task

// This struct contains the action to be perforrmed by the recommendation and set of rules wrt the action.
type Action struct {
	// ActionName indicates the name of the action to recommend by the recommendation engine.
	ActionName string `yaml:"action_name"`
	// Rules indicates list of rules to evaluate the criteria for the recommendation engine.
	Rules []Rule `yaml:"rules"`
	// Operator indicates the logical operation needs to be performed while executing the rules
	Operator string `yaml:"operator"`
}

// This struct contains the rule.
type Rule struct {
	// Metic indicates the name of the metric. These can be:
	// 	Cpu
	//	Mem
	//	Shard
	Metric string `yaml:"metric"`
	// Limit indicates the threshold value for a metric.
	// If this threshold is achieved for a given metric for the decision periond then the rule will be activated.
	Limit float32 `yaml:"limit"`
	// Stat indicates the statistics on which the evaluation of the rule will happen.
	// For Cpu and Mem the values can be:
	// 	Avg: The average CPU or MEM value will be calculated for a given decision period.
	//  Count: The number of occurences where CPU or MEM value crossed the threshold limit.
	// For rule: Shard, the stat will not be applicable as the shard will be calculated across the cluster and is not a statistical value.
	Stat string `yaml:"stat"`
	// DecisionPeriod indicates the time in minutes for which a rule is evalated.
	DecisionPeriod int `yaml:"decision_period"`
	// Occurences indicate the number of time a rule reached the threshold limit for a give decision period.
	// It will be applicable only when the Stat is set to Count.
	Occurences int `yaml:"occurences"`
}

// This struct contains the task details which is set of actions.
type Task struct {
	// Actions indicates list of actions.
	// An action indicates what operation needs to recommended by recommendation engine.
	// As of now actions can be of two types:
	//
	//	scale_up_by_1
	//	scale_down_by_1
	Actions []Action `yaml:"action"`
}
