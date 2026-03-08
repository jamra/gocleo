package fst

// RegexToNFA converts a regex pattern to NFA for benchmarking and testing
func RegexToNFA(pattern string) (*NFA, error) {
	automaton, err := NewTrueRegexAutomaton(pattern)
	if err != nil {
		return nil, err
	}
	return automaton.nfa, nil
}

// NFAtoDFA converts NFA to DFA using subset construction
// This is a simplified implementation for demonstration
func NFAtoDFA(nfa *NFA) *DFA {
	// For now, return a simplified DFA structure
	// In a full implementation, this would perform subset construction
	return &DFA{
		states: make(map[DFAStateID]*DFAState),
		start:  0,
		accepting: make([]DFAStateID, 0),
	}
}

// DFAStateID uniquely identifies a DFA state
type DFAStateID int

// DFAState represents a state in the DFA
type DFAState struct {
	id          DFAStateID
	isAccepting bool
	transitions map[byte]DFAStateID
}

// DFA represents a deterministic finite automaton
type DFA struct {
	start     DFAStateID
	accepting []DFAStateID
	states    map[DFAStateID]*DFAState
}

// Accept tests if the DFA accepts the given input
func (dfa *DFA) Accept(input string) bool {
	current := dfa.start
	
	for _, char := range []byte(input) {
		state, exists := dfa.states[current]
		if !exists {
			return false
		}
		
		next, hasTransition := state.transitions[char]
		if !hasTransition {
			return false
		}
		
		current = next
	}
	
	// Check if final state is accepting
	if state, exists := dfa.states[current]; exists {
		return state.isAccepting
	}
	
	return false
}