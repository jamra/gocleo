package fst

import (
	"fmt"
	"sort"
)

// State represents a single state in the automaton
type State struct {
	ID          uint32
	IsFinal     bool
	Output      uint64  // For FST - the value associated with this state
	Transitions []Transition
}

// Transition represents a labeled edge between states
type Transition struct {
	Label byte   // The input character/byte
	Target uint32 // Target state ID
	Output uint64 // Output for this transition (FST only)
}

// Automaton represents a finite state automaton/transducer
type Automaton struct {
	States    []State
	StartState uint32
	NumStates uint32
}

// NewAutomaton creates a new empty automaton
func NewAutomaton() *Automaton {
	return &Automaton{
		States:    make([]State, 0),
		StartState: 0,
		NumStates: 0,
	}
}

// AddState adds a new state to the automaton
func (a *Automaton) AddState(isFinal bool, output uint64) uint32 {
	stateID := a.NumStates
	state := State{
		ID:          stateID,
		IsFinal:     isFinal,
		Output:      output,
		Transitions: make([]Transition, 0),
	}
	a.States = append(a.States, state)
	a.NumStates++
	return stateID
}

// AddTransition adds a transition from one state to another
func (a *Automaton) AddTransition(fromState uint32, label byte, toState uint32, output uint64) {
	if fromState >= a.NumStates {
		panic(fmt.Sprintf("invalid from state: %d", fromState))
	}
	
	transition := Transition{
		Label:  label,
		Target: toState,
		Output: output,
	}
	
	a.States[fromState].Transitions = append(a.States[fromState].Transitions, transition)
	
	// Keep transitions sorted by label for binary search
	sort.Slice(a.States[fromState].Transitions, func(i, j int) bool {
		return a.States[fromState].Transitions[i].Label < a.States[fromState].Transitions[j].Label
	})
}

// GetState returns the state with the given ID
func (a *Automaton) GetState(stateID uint32) *State {
	if stateID >= a.NumStates {
		return nil
	}
	return &a.States[stateID]
}

// FindTransition finds a transition from the given state with the given label
func (a *Automaton) FindTransition(stateID uint32, label byte) *Transition {
	if stateID >= a.NumStates {
		return nil
	}
	
	state := &a.States[stateID]
	transitions := state.Transitions
	
	// Binary search for the transition
	i := sort.Search(len(transitions), func(i int) bool {
		return transitions[i].Label >= label
	})
	
	if i < len(transitions) && transitions[i].Label == label {
		return &transitions[i]
	}
	
	return nil
}

// Accept tests if the automaton accepts the given input
func (a *Automaton) Accept(input []byte) bool {
	currentState := a.StartState
	
	for _, b := range input {
		transition := a.FindTransition(currentState, b)
		if transition == nil {
			return false
		}
		currentState = transition.Target
	}
	
	finalState := a.GetState(currentState)
	return finalState != nil && finalState.IsFinal
}

// AcceptWithOutput tests if the automaton accepts the input and returns the output
func (a *Automaton) AcceptWithOutput(input []byte) (bool, uint64) {
	currentState := a.StartState
	totalOutput := uint64(0)
	
	for _, b := range input {
		transition := a.FindTransition(currentState, b)
		if transition == nil {
			return false, 0
		}
		totalOutput += transition.Output
		currentState = transition.Target
	}
	
	finalState := a.GetState(currentState)
	if finalState == nil || !finalState.IsFinal {
		return false, 0
	}
	
	totalOutput += finalState.Output
	return true, totalOutput
}

// AutomatonBuilder helps build automata efficiently
type AutomatonBuilder struct {
	automaton *Automaton
	registry  map[string]uint32 // For state deduplication
}

// NewAutomatonBuilder creates a new automaton builder
func NewAutomatonBuilder() *AutomatonBuilder {
	return &AutomatonBuilder{
		automaton: NewAutomaton(),
		registry:  make(map[string]uint32),
	}
}

// Build returns the constructed automaton
func (ab *AutomatonBuilder) Build() *Automaton {
	return ab.automaton
}

// BuildFromStrings builds an automaton from a sorted list of strings
func (ab *AutomatonBuilder) BuildFromStrings(keys []string) *Automaton {
	if len(keys) == 0 {
		// Empty automaton
		ab.automaton.AddState(false, 0)
		return ab.automaton
	}
	
	// Add initial state
	startState := ab.automaton.AddState(false, 0)
	ab.automaton.StartState = startState
	
	// Build trie-like structure
	ab.buildRecursive(keys, 0, startState)
	
	return ab.automaton
}

// buildRecursive recursively builds the automaton from sorted strings
// buildRecursive recursively builds the automaton from sorted strings
func (ab *AutomatonBuilder) buildRecursive(keys []string, depth int, stateID uint32) {
	if len(keys) == 0 {
		return
	}
	
	// Group keys by their character at current depth
	groups := make(map[byte][]string)
	var hasEmptyKey bool
	
	for _, key := range keys {
		if depth >= len(key) {
			hasEmptyKey = true
			continue
		}
		
		char := key[depth]
		groups[char] = append(groups[char], key)
	}
	
	// Mark state as final if we have an empty key
	if hasEmptyKey {
		ab.automaton.States[stateID].IsFinal = true
	}
	
	// Process each character group (only process existing characters)
	for char, group := range groups {
		// Filter group to only include keys that continue past this character
		var filteredGroup []string
		var hasTerminatingKey bool
		
		for _, key := range group {
			if depth+1 < len(key) {
				// Key continues beyond this character
				filteredGroup = append(filteredGroup, key)
			} else if depth+1 == len(key) {
				// This key ends exactly at the next depth
				hasTerminatingKey = true
			}
		}
		
		// Create target state
		targetState := ab.automaton.AddState(false, 0)
		ab.automaton.AddTransition(stateID, char, targetState, 0)
		
		// Mark target as final if any key terminates there
		if hasTerminatingKey {
			ab.automaton.States[targetState].IsFinal = true
		}
		
		// Recursively build for filtered group (only keys that continue)
		if len(filteredGroup) > 0 {
			ab.buildRecursive(filteredGroup, depth+1, targetState)
		}
	}
}