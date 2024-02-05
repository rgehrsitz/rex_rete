package rete_test

import (
	"encoding/json"
	"rgehrsitz/rexrete/pkg/rete" // Adjust to your module's import path
	"rgehrsitz/rexrete/pkg/rules"
	"testing"
)

func TestRuleEvaluation(t *testing.T) {
	// Example JSON rule string (replace with actual JSON rule)
	jsonRule := `{
        "name": "AdultUser",
        "priority": 1,
        "conditions": {
            "all": [{"fact": "age", "operator": "greaterThanOrEqual", "value": 18}]
        },
        "event": {"eventType": "UserIsAdult", "customProperty": "User has reached adulthood."}
    }`

	// Step 1: Parse JSON rule into a Rule object
	var rule rules.Rule
	if err := json.Unmarshal([]byte(jsonRule), &rule); err != nil {
		t.Fatalf("Failed to parse rule: %v", err)
	}

	// Step 2: Construct the Rete network based on the rule
	network := rete.NewNetwork()
	network.LoadRule(rule) // Assuming LoadRule handles the construction

	// Step 3: Insert facts and run the evaluation
	network.AddFact(rete.NewWME("age", 20)) // Example fact
	triggeredEvents := network.Evaluate()   // Assuming Evaluate runs the network and returns triggered events

	// Step 4: Verify the expected event is triggered
	if len(triggeredEvents) != 1 || triggeredEvents[0].EventType != "UserIsAdult" {
		t.Errorf("Expected UserIsAdult event to be triggered")
	}
}
