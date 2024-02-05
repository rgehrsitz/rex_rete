package rules

import (
	"encoding/json"
)

// Rule represents a rule defined in the system, including its conditions and associated event.
type Rule struct {
	Name       string     `json:"name"`
	Priority   int        `json:"priority"`
	Conditions Conditions `json:"conditions"`
	Event      RuleEvent  `json:"event"`
}

// Conditions holds all the conditions for a rule, including "all" and "any" conditions.
type Conditions struct {
	All []Condition `json:"all"`
	Any []Condition `json:"any"`
}

// Condition represents a single condition within a rule, specifying the fact, operator, and value.
type Condition struct {
	Fact     string      `json:"fact"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"` // Use interface{} to allow different types of values
}

// RuleEvent represents the event that is triggered when a rule's conditions are met.
type RuleEvent struct {
	EventType      string        `json:"eventType"`
	CustomProperty string        `json:"customProperty"`
	Facts          []interface{} `json:"facts"`  // Placeholder for facts that triggered the event
	Values         []interface{} `json:"values"` // Placeholder for values corresponding to the facts
}

func ParseRuleFromJSON(jsonData string) (Rule, error) {
	var rule Rule
	err := json.Unmarshal([]byte(jsonData), &rule)
	if err != nil {
		return Rule{}, err
	}
	return rule, nil
}
