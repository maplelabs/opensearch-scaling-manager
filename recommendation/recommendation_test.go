package recommendation

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/opensearch-project/opensearch-go/opensearchtransport"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type Client struct {
	*opensearchapi.API // Embeds the API methods
	Transport          opensearchtransport.Interface
}

func TestTaskNotRecommendedOr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 1,
				"Min": 0,
				"Max": 1,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 30,
				"Min": 20,
				"Max": 80,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskRecommendedOr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 1, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 4,
				"max": 12,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, true, isRecommendedTask)
}

func TestTaskRecommendedOr1(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_down_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 1,
				"Min": 0,
				"Max": 1,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, true, isRecommendedTask)
}

func TestTaskNotRecommendedOr1(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_down_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 1, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 29, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 4,
				"max": 12,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 30,
				"Min": 20,
				"Max": 80,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskNotRecommendedAnd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 1,
				"Min": 0,
				"Max": 1,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 1,
				"Min": 0,
				"Max": 1,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskNotRecommendedAnd1(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 1, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 59, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 70, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 0,
				"max": 1,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 60,
				"min": 0,
				"max": 1,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskNotRecommendedAnd2(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_down_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 10, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 61, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 50, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 5,
				"min": 5,
				"max": 10,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 60,
				"min": 0,
				"max": 1,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskRecommendedAnd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 1, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 10, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 12,
				"max": 8,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 20,
				"min": 80,
				"max": 40,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, true, isRecommendedTask)
}

func TestTaskRecommendedAnd1(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_down_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 5, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 30, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 12,
				"max": 8,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 20,
				"min": 80,
				"max": 40,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, true, isRecommendedTask)
}

func TestTaskNotEnoughDataAnd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 1, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 10, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 10, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 12,
				"max": 8,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(400, "Not enough Data points")
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskNotEnoughDataAnd1(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 1, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 10, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 10, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 12,
				"max": 8,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 20,
				"min": 80,
				"max": 40,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(400, "Not enough Data points")
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskNotEnoughDataOr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 5, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 12,
				"max": 8,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(400, "Not enough Data points")
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskDecisionPeriodSmallAnd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 1, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 10, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 10, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 12,
				"max": 8,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(400, "Decision period too small")
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskDecisionPeriodSmallOr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 5, stat: AVG, decision_period: 9}, {metric: MemUtil, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=CpuUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 12,
				"max": 8,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg?metric=MemUtil&duration=9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(400, "Decision period too small")
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskNotRecommendedOrCountTerm(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: MemUtil, limit: 59, stat: COUNT, occurences: 12, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=CpuUtil&duration=9&threshold=1.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 3,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=MemUtil&duration=9&threshold=59.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 4,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskRecommendedOrCountTerm(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: MemUtil, limit: 59, stat: COUNT, occurences: 12, decision_period: 9}]}`

	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=CpuUtil&duration=9&threshold=1.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 6,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=MemUtil&duration=9&threshold=59.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 13,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskNotRecommendedOrCountTerm1(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_down_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 1, stat: COUNT, occurences: 2, decision_period: 9}, {metric: MemUtil, limit: 59, stat: COUNT, occurences: 12, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=CpuUtil&duration=9&threshold=1.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 3,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=MemUtil&duration=9&threshold=59.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 13,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskRecommendedOrCountTerm1(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_down_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 1, stat: COUNT, occurences: 5, decision_period: 9}, {metric: MemUtil, limit: 59, stat: COUNT, occurences: 12, decision_period: 9}]}`

	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=CpuUtil&duration=9&threshold=1.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 6,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=MemUtil&duration=9&threshold=59.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 10,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskRecommendedAndCountTerm(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 1.0, stat: COUNT, occurences: 10, decision_period: 9}, {metric: MemUtil, limit: 59.0, stat: COUNT, occurences: 12, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=CpuUtil&duration=9&threshold=1.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 11,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=MemUtil&duration=9&threshold=59.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 13,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskNotRecommendedAndCountTerm(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 1, stat: COUNT, occurences: 3, decision_period: 9}, {metric: MemUtil, limit: 59, stat: COUNT, occurences: 12, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=CpuUtil&duration=9&threshold=1.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 4,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=MemUtil&duration=9&threshold=59.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 4,
			})
			return resp, err
		},
	)
	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskRecommendedAndCountTerm1(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_down_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 1.0, stat: COUNT, occurences: 10, decision_period: 9}, {metric: MemUtil, limit: 59.0, stat: COUNT, occurences: 12, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=CpuUtil&duration=9&threshold=1.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 9,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=MemUtil&duration=9&threshold=59.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 11,
			})
			return resp, err
		},
	)

	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskNotRecommendedAndCountTerm1(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_down_by_1, operator: AND, rules: [{metric: CpuUtil, limit: 1, stat: COUNT, occurences: 3, decision_period: 9}, {metric: MemUtil, limit: 59, stat: COUNT, occurences: 12, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=CpuUtil&duration=9&threshold=1.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 2,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/violated?metric=MemUtil&duration=9&threshold=59.000000",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"violated_count": 13,
			})
			return resp, err
		},
	)
	isRecommendedTask, _ := task.GetNextTask(5, true, false)
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}
