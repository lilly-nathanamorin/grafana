package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/services/ngalert/state"

	"github.com/grafana/grafana/pkg/infra/log"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana/pkg/services/ngalert/eval"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/stretchr/testify/assert"
)

func TestProcessEvalResults(t *testing.T) {
	evaluationTime, err := time.Parse("2006-01-02", "2021-03-25")
	if err != nil {
		t.Fatalf("error parsing date format: %s", err.Error())
	}
	evaluationDuration := 10 * time.Millisecond

	testCases := []struct {
		desc           string
		alertRule      *models.AlertRule
		evalResults    []eval.Results
		expectedStates map[string]*state.State
	}{
		{
			desc: "a cache entry is correctly created",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Normal,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
					},
					LastEvaluationTime: evaluationTime,
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test"},
				},
			},
		},
		{
			desc: "two results create two correct cache entries",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label_1": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
					eval.Result{
						Instance:           data.Labels{"instance_label_2": "test"},
						State:              eval.Alerting,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid alertname:test_title instance_label_1:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid alertname:test_title instance_label_1:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label_1":             "test",
					},
					State: eval.Normal,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
					},
					LastEvaluationTime: evaluationTime,
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test"},
				},
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid alertname:test_title instance_label_2:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid alertname:test_title instance_label_2:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label_2":             "test",
					},
					State: eval.Alerting,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Alerting,
						},
					},
					StartsAt:           evaluationTime,
					EndsAt:             evaluationTime.Add(20 * time.Second),
					LastEvaluationTime: evaluationTime,
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test"},
				},
			},
		},
		{
			desc: "state is maintained",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_1",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime.Add(1 * time.Minute),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_1 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_1",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_1 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_1",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Normal,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(1 * time.Minute),
							EvaluationState: eval.Normal,
						},
					},
					LastEvaluationTime: evaluationTime.Add(1 * time.Minute),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test"},
				},
			},
		},
		{
			desc: "normal -> alerting transition when For is unset",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_2",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Alerting,
						EvaluatedAt:        evaluationTime.Add(1 * time.Minute),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_2",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_2",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Alerting,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(1 * time.Minute),
							EvaluationState: eval.Alerting,
						},
					},
					StartsAt:           evaluationTime.Add(1 * time.Minute),
					EndsAt:             evaluationTime.Add(1 * time.Minute).Add(time.Duration(20) * time.Second),
					LastEvaluationTime: evaluationTime.Add(1 * time.Minute),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test", "alerting_at": "2021-03-25 00:01:00 +0000 UTC"},
				},
			},
		},
		{
			desc: "normal -> alerting when For is set",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_2",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
				For:             1 * time.Minute,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Alerting,
						EvaluatedAt:        evaluationTime.Add(10 * time.Second),
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Alerting,
						EvaluatedAt:        evaluationTime.Add(80 * time.Second),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_2",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_2",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Alerting,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(10 * time.Second),
							EvaluationState: eval.Alerting,
						},
						{
							EvaluationTime:  evaluationTime.Add(80 * time.Second),
							EvaluationState: eval.Alerting,
						},
					},
					StartsAt:           evaluationTime.Add(80 * time.Second),
					EndsAt:             evaluationTime.Add(80 * time.Second).Add(1 * time.Minute),
					LastEvaluationTime: evaluationTime.Add(80 * time.Second),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test", "alerting_at": "2021-03-25 00:01:20 +0000 UTC"},
				},
			},
		},
		{
			desc: "normal -> pending when For is set but not exceeded",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_2",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
				For:             1 * time.Minute,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Alerting,
						EvaluatedAt:        evaluationTime.Add(10 * time.Second),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_2",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_2",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Pending,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(10 * time.Second),
							EvaluationState: eval.Alerting,
						},
					},
					StartsAt:           evaluationTime.Add(10 * time.Second),
					EndsAt:             evaluationTime.Add(10 * time.Second).Add(1 * time.Minute),
					LastEvaluationTime: evaluationTime.Add(10 * time.Second),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test"},
				},
			},
		},
		{
			desc: "normal -> alerting when result is NoData and NoDataState is alerting",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_2",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
				NoDataState:     models.Alerting,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.NoData,
						EvaluatedAt:        evaluationTime.Add(10 * time.Second),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_2",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_2",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Alerting,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(10 * time.Second),
							EvaluationState: eval.NoData,
						},
					},
					StartsAt:           evaluationTime.Add(10 * time.Second),
					EndsAt:             evaluationTime.Add(10 * time.Second).Add(20 * time.Second),
					LastEvaluationTime: evaluationTime.Add(10 * time.Second),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test", "no_data": "2021-03-25 00:00:10 +0000 UTC"},
				},
			},
		},
		{
			desc: "normal -> nodata when result is NoData and NoDataState is nodata",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_2",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
				NoDataState:     models.NoData,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.NoData,
						EvaluatedAt:        evaluationTime.Add(10 * time.Second),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_2",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_2",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.NoData,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(10 * time.Second),
							EvaluationState: eval.NoData,
						},
					},
					StartsAt:           evaluationTime.Add(10 * time.Second),
					EndsAt:             evaluationTime.Add(10 * time.Second).Add(20 * time.Second),
					LastEvaluationTime: evaluationTime.Add(10 * time.Second),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test", "no_data": "2021-03-25 00:00:10 +0000 UTC"},
				},
			},
		},
		{
			desc: "normal -> normal when result is NoData and NoDataState is ok",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_2",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
				NoDataState:     models.OK,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.NoData,
						EvaluatedAt:        evaluationTime.Add(10 * time.Second),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_2",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_2",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Normal,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(10 * time.Second),
							EvaluationState: eval.NoData,
						},
					},
					StartsAt:           evaluationTime.Add(10 * time.Second),
					EndsAt:             evaluationTime.Add(10 * time.Second).Add(20 * time.Second),
					LastEvaluationTime: evaluationTime.Add(10 * time.Second),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test", "no_data": "2021-03-25 00:00:10 +0000 UTC"},
				},
			},
		},
		{
			desc: "EndsAt set correctly. normal -> alerting when result is NoData and NoDataState is alerting and For is set and For is breached",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_2",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
				For:             1 * time.Minute,
				NoDataState:     models.Alerting,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.NoData,
						EvaluatedAt:        evaluationTime.Add(10 * time.Second),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_2",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_2",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Alerting,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(10 * time.Second),
							EvaluationState: eval.NoData,
						},
					},
					StartsAt:           evaluationTime.Add(10 * time.Second),
					EndsAt:             evaluationTime.Add(10 * time.Second).Add(1 * time.Minute),
					LastEvaluationTime: evaluationTime.Add(10 * time.Second),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test", "no_data": "2021-03-25 00:00:10 +0000 UTC"},
				},
			},
		},
		{
			desc: "normal -> normal when result is NoData and NoDataState is KeepLastState",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_2",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
				For:             1 * time.Minute,
				NoDataState:     models.KeepLastState,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.NoData,
						EvaluatedAt:        evaluationTime.Add(10 * time.Second),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_2",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_2",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Normal,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(10 * time.Second),
							EvaluationState: eval.NoData,
						},
					},
					StartsAt:           evaluationTime.Add(10 * time.Second),
					EndsAt:             evaluationTime.Add(10 * time.Second).Add(1 * time.Minute),
					LastEvaluationTime: evaluationTime.Add(10 * time.Second),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test", "no_data": "2021-03-25 00:00:10 +0000 UTC"},
				},
			},
		},
		{
			desc: "normal -> normal when result is NoData and NoDataState is KeepLastState",
			alertRule: &models.AlertRule{
				OrgID:           1,
				Title:           "test_title",
				UID:             "test_alert_rule_uid_2",
				NamespaceUID:    "test_namespace_uid",
				Annotations:     map[string]string{"annotation": "test"},
				Labels:          map[string]string{"label": "test"},
				IntervalSeconds: 10,
				For:             1 * time.Minute,
				NoDataState:     models.KeepLastState,
			},
			evalResults: []eval.Results{
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.Normal,
						EvaluatedAt:        evaluationTime,
						EvaluationDuration: evaluationDuration,
					},
				},
				{
					eval.Result{
						Instance:           data.Labels{"instance_label": "test"},
						State:              eval.NoData,
						EvaluatedAt:        evaluationTime.Add(10 * time.Second),
						EvaluationDuration: evaluationDuration,
					},
				},
			},
			expectedStates: map[string]*state.State{
				"map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]": {
					AlertRuleUID: "test_alert_rule_uid_2",
					OrgID:        1,
					CacheId:      "map[__alert_rule_namespace_uid__:test_namespace_uid __alert_rule_uid__:test_alert_rule_uid_2 alertname:test_title instance_label:test label:test]",
					Labels: data.Labels{
						"__alert_rule_namespace_uid__": "test_namespace_uid",
						"__alert_rule_uid__":           "test_alert_rule_uid_2",
						"alertname":                    "test_title",
						"label":                        "test",
						"instance_label":               "test",
					},
					State: eval.Normal,
					Results: []state.Evaluation{
						{
							EvaluationTime:  evaluationTime,
							EvaluationState: eval.Normal,
						},
						{
							EvaluationTime:  evaluationTime.Add(10 * time.Second),
							EvaluationState: eval.NoData,
						},
					},
					StartsAt:           evaluationTime.Add(10 * time.Second),
					EndsAt:             evaluationTime.Add(10 * time.Second).Add(1 * time.Minute),
					LastEvaluationTime: evaluationTime.Add(10 * time.Second),
					EvaluationDuration: evaluationDuration,
					Annotations:        map[string]string{"annotation": "test", "no_data": "2021-03-25 00:00:10 +0000 UTC"},
				},
			},
		},
	}

	for _, tc := range testCases {
		st := state.NewManager(log.New("test_state_manager"))
		t.Run(tc.desc, func(t *testing.T) {
			for _, res := range tc.evalResults {
				_ = st.ProcessEvalResults(tc.alertRule, res)
			}
			for id, s := range tc.expectedStates {
				cachedState, err := st.Get(id)
				require.NoError(t, err)
				assert.Equal(t, s, cachedState)
			}
		})
	}
}
