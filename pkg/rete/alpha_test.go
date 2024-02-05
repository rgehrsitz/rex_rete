package rete

import (
	"container/list"
	"fmt"
	"testing"
)

func TestAlphaMemoryActivation(t *testing.T) {
	// Initialize a new AlphaMemory instance.
	alphaMem := &AlphaMemory{
		items:      list.New(),
		successors: list.New(),
	}

	// Initialize a new WME and activate the AlphaMemory with this WME.
	wme := NewWME("ClassName", "ID", "Attribute", "Value")

	alphaMem.activation(wme)

	// Assert that the AlphaMemory's items list correctly contains the WME.
	if alphaMem.items.Len() != 1 {
		t.Errorf("Expected 1 WME in AlphaMemory, got %d", alphaMem.items.Len())
	} else {
		firstItem := alphaMem.items.Front().Value.(*WME)
		if firstItem != wme {
			t.Errorf("Expected AlphaMemory items to contain the activated WME")
		}
	}

	// Assert that the WME's alphaMems list correctly includes the AlphaMemory.
	if wme.alphaMems.Len() != 1 {
		t.Errorf("Expected AlphaMemory to be in WME's alphaMems, got %d", wme.alphaMems.Len())
	} else {
		firstAlphaMem := wme.alphaMems.Front().Value.(*AlphaMemory)
		if firstAlphaMem != alphaMem {
			t.Errorf("Expected WME's alphaMems to contain the AlphaMemory")
		}
	}
}

func TestConstantTestNodeActivationSuccessPath(t *testing.T) {
	// Correct initialization of AlphaMemory within ConstantTestNode
	outputMemory := &AlphaMemory{
		items:      list.New(),
		successors: list.New(),
	}

	node := ConstantTestNode{
		fieldToTest:    3, // Corrected field index
		fieldMustEqual: "Value",
		outputMemory:   outputMemory, // Use the properly initialized AlphaMemory
		children:       list.New(),
	}

	// Initialize matching and non-matching WMEs
	matchingWME := NewWME("ClassName", "ID", "Attribute", "Value")
	nonMatchingWME := NewWME("ClassName", "ID", "Attribute", "DifferentValue")

	// Activate the node with the matching WME
	node.activation(matchingWME)
	// Verify that the matching WME is added to the node's output memory
	if node.outputMemory.items.Len() != 1 {
		t.Errorf("Expected 1 matching WME in the output memory, got %d", node.outputMemory.items.Len())
	} else {
		firstItem := node.outputMemory.items.Front().Value.(*WME)
		if firstItem != matchingWME {
			t.Errorf("Expected the matching WME to be in the output memory")
		}
	}

	// Reset the output memory for the next test
	node.outputMemory.items = list.New()

	// Activate the node with the non-matching WME
	node.activation(nonMatchingWME)
	// Verify that the non-matching WME is not added to the node's output memory
	if node.outputMemory.items.Len() != 0 {
		t.Errorf("Expected 0 non-matching WMEs in the output memory, got %d", node.outputMemory.items.Len())
	}
}

func TestConstantTestNodeActivationFailurePath(t *testing.T) {
	// Initialize the ConstantTestNode with a specific fieldToTest and fieldMustEqual.
	fieldIndex := 3 // Index for "Value" field
	expectedValue := "Value"
	node := ConstantTestNode{
		fieldToTest:    fieldIndex,
		fieldMustEqual: expectedValue,
		outputMemory: &AlphaMemory{
			items:      list.New(),
			successors: list.New(),
		},
		children: list.New(),
	}

	// Create a child node to test if it gets activated (should not for the failure path).
	childNode := &ConstantTestNode{
		outputMemory: &AlphaMemory{
			items:      list.New(),
			successors: list.New(),
		},
	}
	node.children.PushBack(childNode) // Add child node to test its activation indirectly.

	// Create a non-matching WME.
	nonMatchingWME := NewWME("ClassName", "ID", "Attribute", "DifferentValue")

	// Activate the node with the non-matching WME.
	node.activation(nonMatchingWME)

	// Verify that the non-matching WME does not activate the node's output memory.
	if node.outputMemory.items.Len() != 0 {
		t.Errorf("Expected 0 WMEs in the output memory for a non-matching activation, got %d", node.outputMemory.items.Len())
	}

	// Verify that the non-matching WME does not activate the child nodes.
	if childNode.outputMemory.items.Len() != 0 {
		t.Errorf("Expected 0 WMEs in the child node's output memory for a non-matching activation, got %d", childNode.outputMemory.items.Len())
	}
}

type TestSuccessor struct {
	activatedWith *WME
}

func (ts *TestSuccessor) GetNodeType() string {
	return "TestSuccessor"
}

func (ts *TestSuccessor) GetItems() *list.List {
	return list.New()
}

func (ts *TestSuccessor) GetParent() IReteNode {
	return nil
}

func (ts *TestSuccessor) GetChildren() *list.List {
	return list.New()
}

func (ts *TestSuccessor) LeftActivation(t *Token, w *WME, b Env) {
	// Not needed for this test
}

func (ts *TestSuccessor) RightActivation(w *WME) {
	ts.activatedWith = w
}

func TestAlphaMemoryPropagationToSuccessors(t *testing.T) {
	// Create an AlphaMemory.
	alphaMem := &AlphaMemory{
		items:      list.New(),
		successors: list.New(),
	}

	// Attach a mock IReteNode as a successor.
	testSuccessor := &TestSuccessor{}
	alphaMem.successors.PushBack(testSuccessor)

	// Create a WME and activate the AlphaMemory with it.
	wme := NewWME("ClassName", "ID", "Attribute", "Value")
	alphaMem.activation(wme)

	// Verify that the RightActivation method was called on the successor with the correct WME.
	if testSuccessor.activatedWith != wme {
		t.Errorf("Expected RightActivation to be called with the correct WME")
	}
}

func TestChainActivationThroughConstantTestNodeSimplified(t *testing.T) {
	// Initialize the ConstantTestNode with an AlphaMemory that has properly initialized lists.
	rootNode := ConstantTestNode{
		fieldToTest:    3, // Assuming "Value" field index
		fieldMustEqual: "Value1",
		outputMemory: &AlphaMemory{
			items:      list.New(), // Ensure items list is initialized
			successors: list.New(), // Ensure successors list is initialized
		},
		children: list.New(),
	}

	// Similar initialization for childNode with proper AlphaMemory setup
	childNode := ConstantTestNode{
		outputMemory: &AlphaMemory{
			items:      list.New(), // Ensure items list is initialized
			successors: list.New(), // Ensure successors list is initialized
		},
	}

	for e := rootNode.children.Front(); e != nil; e = e.Next() {
		child := e.Value.(*ConstantTestNode)
		fmt.Println("Child node found with criteria:", child.fieldToTest, child.fieldMustEqual)
	}

	// Attach childNode directly to rootNode
	rootNode.children.PushBack(&childNode)
	// Verify the attachment
	if rootNode.children.Back().Value.(*ConstantTestNode) != &childNode {
		t.Error("childNode was not properly attached to rootNode")
	}

	// Create a matching WME
	matchingWME := NewWME("ClassName", "ID", "Attribute", "Value1")

	// Activate rootNode with the matching WME
	rootNode.activation(matchingWME)

	// Verify propagation to childNode's outputMemory
	if childNode.outputMemory.items.Len() != 1 {
		t.Errorf("Expected the matching WME to propagate to the child node, but it did not")
	}
}
