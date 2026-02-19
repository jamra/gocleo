package fst

// MinimizeAutomaton is a convenience function that minimizes an automaton
func MinimizeAutomaton(automaton *Automaton) *Automaton {
	if len(automaton.States) <= 1 {
		return automaton // Already minimal
	}

	// Simple minimization: remove duplicate states with same behavior
	stateGroups := make(map[string][]int)
	
	for i, state := range automaton.States {
		// Create signature based on accepting status and transitions
		signature := ""
		if state.IsFinal {
			signature = "F"
		} else {
			signature = "N"
		}
		
		// Add transition signature
		for _, trans := range state.Transitions {
			signature += string(trans.Label) + string(rune(trans.Target))
		}
		
		stateGroups[signature] = append(stateGroups[signature], i)
	}

	// If no groups can be merged, return original
	if len(stateGroups) == len(automaton.States) {
		return automaton
	}

	// Build minimized automaton
	newStates := make([]State, 0, len(stateGroups))
	stateMapping := make(map[uint32]uint32)
	newStateID := uint32(0)

	for _, group := range stateGroups {
		// Use first state as representative
		representative := group[0]
		oldState := automaton.States[representative]
		
		newState := State{
			ID:          newStateID,
			IsFinal:     oldState.IsFinal,
			Output:      oldState.Output,
			Transitions: make([]Transition, len(oldState.Transitions)),
		}
		
		copy(newState.Transitions, oldState.Transitions)
		newStates = append(newStates, newState)
		
		// Map all states in group to new state
		for _, oldStateID := range group {
			stateMapping[uint32(oldStateID)] = newStateID
		}
		
		newStateID++
	}

	// Update transition targets
	for i := range newStates {
		for j := range newStates[i].Transitions {
			if newTarget, exists := stateMapping[newStates[i].Transitions[j].Target]; exists {
				newStates[i].Transitions[j].Target = newTarget
			}
		}
	}

	return &Automaton{
		States:     newStates,
		StartState: stateMapping[automaton.StartState],
		NumStates:  uint32(len(newStates)),
	}
}

// CalculateCompressionRatio calculates the compression achieved by minimization
func CalculateCompressionRatio(original, minimized *Automaton) float64 {
	if len(original.States) == 0 {
		return 0.0
	}
	return float64(len(minimized.States)) / float64(len(original.States))
}

// MinimizationStats provides statistics about the minimization process
type MinimizationStats struct {
	OriginalStates    int
	MinimizedStates   int
	StatesRemoved     int
	CompressionRatio  float64
	SpaceSavingPct    float64
}

// GetMinimizationStats returns detailed statistics about minimization
func GetMinimizationStats(original, minimized *Automaton) MinimizationStats {
	originalCount := len(original.States)
	minimizedCount := len(minimized.States)
	removed := originalCount - minimizedCount
	ratio := CalculateCompressionRatio(original, minimized)
	savings := (1.0 - ratio) * 100.0

	return MinimizationStats{
		OriginalStates:   originalCount,
		MinimizedStates:  minimizedCount,
		StatesRemoved:    removed,
		CompressionRatio: ratio,
		SpaceSavingPct:   savings,
	}
}
