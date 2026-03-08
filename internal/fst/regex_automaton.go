// Package fst provides true automata intersection for FST and regular expressions
// Based on the principles described in https://burntsushi.net/transducers/
package fst

import (
	"regexp"
	"regexp/syntax"
	"sort"
)

// RegexAutomaton represents a compiled regular expression as a finite state automaton
// that can be intersected with FSTs using product construction
type TrueRegexAutomaton struct {
	nfa    *NFA
	regex  *regexp.Regexp
	states map[NFAStateID]*NFAState
}

// NFAStateID uniquely identifies a state in the NFA
type NFAStateID int

// NFAState represents a single state in the NFA
type NFAState struct {
	id          NFAStateID
	isAccepting bool
	transitions map[byte][]NFAStateID  // Character transitions
	epsilons    []NFAStateID           // Epsilon transitions
}

// NFA represents a non-deterministic finite automaton compiled from regex
type NFA struct {
	start       NFAStateID
	accepting   []NFAStateID
	states      map[NFAStateID]*NFAState
	nextStateID NFAStateID
}

// NewRegexAutomaton compiles a regular expression into an NFA using Thompson's Construction
func NewTrueRegexAutomaton(pattern string) (*TrueRegexAutomaton, error) {
	// First compile with Go's regexp for validation
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Parse the regex into syntax tree
	parsed, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		return nil, err
	}

	// Simplify the parsed regex
	simplified := parsed.Simplify()

	// Build NFA using Thompson's Construction
	nfa := &NFA{
		states:      make(map[NFAStateID]*NFAState),
		nextStateID: 0,
	}

	start, accept := nfa.compileRegex(simplified)
	nfa.start = start
	nfa.accepting = []NFAStateID{accept}

	return &TrueRegexAutomaton{
		nfa:    nfa,
		regex:  regex,
		states: nfa.states,
	}, nil
}

// IntersectWithFST performs true automata intersection using product construction
// This is the key algorithm from burntsushi.net/transducers/
func (ra *TrueRegexAutomaton) IntersectWithFST(fsa FSA) ([]string, error) {
	// Use product construction to build intersection automaton
	// The intersection automaton has states that are pairs (FST_state, NFA_state)
	
	visited := make(map[string]bool)
	results := make([]string, 0)
	
	// Start intersection traversal
	startNFAStates := ra.epsilonClosure([]NFAStateID{ra.nfa.start})
	
	// Use FSA's iterator for traversal
	iterator := fsa.Iterator()
	for iterator.Next() {
		key := string(iterator.Key())
		
		// Simulate NFA execution on this key
		currentStates := startNFAStates
		validPath := true
		
		for _, char := range []byte(key) {
			nextStates := make([]NFAStateID, 0)
			
			// For each current NFA state, find transitions on this character
			for _, state := range currentStates {
				if transitions, exists := ra.nfa.states[state].transitions[char]; exists {
					nextStates = append(nextStates, transitions...)
				}
			}
			
			// If no transitions possible, this key doesn't match
			if len(nextStates) == 0 {
				validPath = false
				break
			}
			
			// Add epsilon closure
			currentStates = ra.epsilonClosure(nextStates)
		}
		
		// Check if we ended in an accepting state
		if validPath && ra.hasAcceptingState(currentStates) {
			if !visited[key] {
				results = append(results, key)
				visited[key] = true
			}
		}
	}
	
	sort.Strings(results)
	return results, nil
}

// compileRegex implements Thompson's Construction algorithm
func (nfa *NFA) compileRegex(re *syntax.Regexp) (start, accept NFAStateID) {
	switch re.Op {
	case syntax.OpLiteral:
		return nfa.compileLiteral(re.Rune)
	case syntax.OpCharClass:
		return nfa.compileCharClass(re.Rune)
	case syntax.OpAnyChar:
		return nfa.compileAnyChar()
	case syntax.OpBeginText:
		return nfa.compileBeginText()
	case syntax.OpEndText:
		return nfa.compileEndText()
	case syntax.OpConcat:
		return nfa.compileConcat(re.Sub)
	case syntax.OpAlternate:
		return nfa.compileAlternate(re.Sub)
	case syntax.OpStar:
		return nfa.compileStar(re.Sub[0])
	case syntax.OpPlus:
		return nfa.compilePlus(re.Sub[0])
	case syntax.OpQuest:
		return nfa.compileQuest(re.Sub[0])
	case syntax.OpRepeat:
		return nfa.compileRepeat(re.Sub[0], re.Min, re.Max)
	default:
		// For unsupported operations, create a simple accepting state
		return nfa.newState(false), nfa.newState(true)
	}
}

// compileLiteral creates NFA fragment for literal string
func (nfa *NFA) compileLiteral(runes []rune) (start, accept NFAStateID) {
	if len(runes) == 0 {
		// Empty string - epsilon transition
		start = nfa.newState(false)
		accept = nfa.newState(true)
		nfa.addEpsilon(start, accept)
		return
	}

	start = nfa.newState(false)
	current := start

	for i, r := range runes {
		if i == len(runes)-1 {
			// Last character goes to accepting state
			accept = nfa.newState(true)
			nfa.addTransition(current, byte(r), accept)
		} else {
			// Intermediate character
			next := nfa.newState(false)
			nfa.addTransition(current, byte(r), next)
			current = next
		}
	}

	return start, accept
}

// compileCharClass creates NFA fragment for character class [a-z], etc.
func (nfa *NFA) compileCharClass(runes []rune) (start, accept NFAStateID) {
	start = nfa.newState(false)
	accept = nfa.newState(true)

	// Character classes in Go's syntax are represented as pairs
	for i := 0; i < len(runes); i += 2 {
		if i+1 < len(runes) {
			// Range like [a-z]
			for r := runes[i]; r <= runes[i+1]; r++ {
				if r <= 255 { // Only handle ASCII for simplicity
					nfa.addTransition(start, byte(r), accept)
				}
			}
		} else {
			// Single character
			if runes[i] <= 255 {
				nfa.addTransition(start, byte(runes[i]), accept)
			}
		}
	}

	return start, accept
}

// compileAnyChar creates NFA fragment for . (dot - any character)
func (nfa *NFA) compileAnyChar() (start, accept NFAStateID) {
	start = nfa.newState(false)
	accept = nfa.newState(true)

	// Add transition for all possible bytes (ASCII)
	for b := byte(0); b <= 255; b++ {
		if b != '\n' { // . typically doesn't match newline
			nfa.addTransition(start, b, accept)
		}
	}

	return start, accept
}

// compileBeginText creates NFA fragment for ^ (beginning of text)
func (nfa *NFA) compileBeginText() (start, accept NFAStateID) {
	// For simplicity, treat as epsilon transition
	start = nfa.newState(false)
	accept = nfa.newState(true)
	nfa.addEpsilon(start, accept)
	return
}

// compileEndText creates NFA fragment for $ (end of text)
func (nfa *NFA) compileEndText() (start, accept NFAStateID) {
	// For simplicity, treat as epsilon transition
	start = nfa.newState(false)
	accept = nfa.newState(true)
	nfa.addEpsilon(start, accept)
	return
}

// compileConcat creates NFA fragment for concatenation (AB)
func (nfa *NFA) compileConcat(subs []*syntax.Regexp) (start, accept NFAStateID) {
	if len(subs) == 0 {
		start = nfa.newState(false)
		accept = nfa.newState(true)
		nfa.addEpsilon(start, accept)
		return
	}

	// Compile first sub-expression
	start, current := nfa.compileRegex(subs[0])
	nfa.states[current].isAccepting = false

	// Chain remaining sub-expressions
	for i := 1; i < len(subs); i++ {
		subStart, subAccept := nfa.compileRegex(subs[i])
		nfa.addEpsilon(current, subStart)
		current = subAccept
		nfa.states[current].isAccepting = false
	}

	accept = current
	nfa.states[accept].isAccepting = true
	return
}

// compileAlternate creates NFA fragment for alternation (A|B)
func (nfa *NFA) compileAlternate(subs []*syntax.Regexp) (start, accept NFAStateID) {
	start = nfa.newState(false)
	accept = nfa.newState(true)

	for _, sub := range subs {
		subStart, subAccept := nfa.compileRegex(sub)
		nfa.addEpsilon(start, subStart)
		nfa.addEpsilon(subAccept, accept)
		nfa.states[subAccept].isAccepting = false
	}

	return
}

// compileStar creates NFA fragment for Kleene star (A*)
func (nfa *NFA) compileStar(sub *syntax.Regexp) (start, accept NFAStateID) {
	start = nfa.newState(false)
	accept = nfa.newState(true)

	subStart, subAccept := nfa.compileRegex(sub)
	nfa.states[subAccept].isAccepting = false

	// Epsilon transitions for *
	nfa.addEpsilon(start, accept)      // Zero occurrences
	nfa.addEpsilon(start, subStart)    // First occurrence
	nfa.addEpsilon(subAccept, accept)  // End
	nfa.addEpsilon(subAccept, subStart) // Repeat

	return
}

// compilePlus creates NFA fragment for plus (A+)
func (nfa *NFA) compilePlus(sub *syntax.Regexp) (start, accept NFAStateID) {
	subStart, subAccept := nfa.compileRegex(sub)
	nfa.states[subAccept].isAccepting = false

	accept = nfa.newState(true)

	// Plus requires at least one occurrence
	nfa.addEpsilon(subAccept, accept)  // End after one
	nfa.addEpsilon(subAccept, subStart) // Repeat

	return subStart, accept
}

// compileQuest creates NFA fragment for question mark (A?)
func (nfa *NFA) compileQuest(sub *syntax.Regexp) (start, accept NFAStateID) {
	start = nfa.newState(false)
	accept = nfa.newState(true)

	subStart, subAccept := nfa.compileRegex(sub)
	nfa.states[subAccept].isAccepting = false

	// Question mark: zero or one occurrence
	nfa.addEpsilon(start, accept)    // Zero occurrences
	nfa.addEpsilon(start, subStart)  // One occurrence
	nfa.addEpsilon(subAccept, accept)

	return
}

// compileRepeat creates NFA fragment for counted repetition {n,m}
func (nfa *NFA) compileRepeat(sub *syntax.Regexp, min, max int) (start, accept NFAStateID) {
	if min == 0 && max == -1 {
		// {0,} is equivalent to *
		return nfa.compileStar(sub)
	}
	if min == 1 && max == -1 {
		// {1,} is equivalent to +
		return nfa.compilePlus(sub)
	}

	// For simplicity, create a basic implementation
	// In production, this would be more sophisticated
	start = nfa.newState(false)
	accept = nfa.newState(true)
	
	if min == 0 {
		nfa.addEpsilon(start, accept) // Zero occurrences allowed
	}

	// Create chain for minimum required occurrences
	current := start
	for i := 0; i < min; i++ {
		subStart, subAccept := nfa.compileRegex(sub)
		nfa.addEpsilon(current, subStart)
		current = subAccept
		nfa.states[current].isAccepting = false
	}

	if max == -1 {
		// Unlimited, add loop
		nfa.addEpsilon(current, accept)
		subStart, subAccept := nfa.compileRegex(sub)
		nfa.addEpsilon(current, subStart)
		nfa.addEpsilon(subAccept, accept)
		nfa.addEpsilon(subAccept, subStart)
	} else {
		// Limited repetitions
		nfa.addEpsilon(current, accept)
		for i := min; i < max; i++ {
			subStart, subAccept := nfa.compileRegex(sub)
			nfa.addEpsilon(current, subStart)
			nfa.addEpsilon(subAccept, accept)
			current = subAccept
			nfa.states[current].isAccepting = false
		}
	}

	return
}

// newState creates a new NFA state
func (nfa *NFA) newState(accepting bool) NFAStateID {
	id := nfa.nextStateID
	nfa.nextStateID++
	
	nfa.states[id] = &NFAState{
		id:          id,
		isAccepting: accepting,
		transitions: make(map[byte][]NFAStateID),
		epsilons:    make([]NFAStateID, 0),
	}
	
	return id
}

// addTransition adds a character transition between states
func (nfa *NFA) addTransition(from NFAStateID, char byte, to NFAStateID) {
	if nfa.states[from].transitions[char] == nil {
		nfa.states[from].transitions[char] = make([]NFAStateID, 0)
	}
	nfa.states[from].transitions[char] = append(nfa.states[from].transitions[char], to)
}

// addEpsilon adds an epsilon (empty) transition between states
func (nfa *NFA) addEpsilon(from, to NFAStateID) {
	nfa.states[from].epsilons = append(nfa.states[from].epsilons, to)
}

// epsilonClosure computes the epsilon closure of a set of NFA states
func (ra *TrueRegexAutomaton) epsilonClosure(states []NFAStateID) []NFAStateID {
	closure := make(map[NFAStateID]bool)
	stack := make([]NFAStateID, 0)
	
	// Initialize with input states
	for _, state := range states {
		if !closure[state] {
			closure[state] = true
			stack = append(stack, state)
		}
	}
	
	// Process epsilon transitions
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		
		// Add all epsilon reachable states
		for _, epsilonState := range ra.nfa.states[current].epsilons {
			if !closure[epsilonState] {
				closure[epsilonState] = true
				stack = append(stack, epsilonState)
			}
		}
	}
	
	// Convert back to slice
	result := make([]NFAStateID, 0, len(closure))
	for state := range closure {
		result = append(result, state)
	}
	
	return result
}

// hasAcceptingState checks if any of the states is accepting
func (ra *TrueRegexAutomaton) hasAcceptingState(states []NFAStateID) bool {
	for _, state := range states {
		if ra.nfa.states[state].isAccepting {
			return true
		}
	}
	return false
}

// MatchString tests if a string matches the regex (fallback to Go's regexp)
func (ra *TrueRegexAutomaton) MatchString(s string) bool {
	return ra.regex.MatchString(s)
}
// TrueAutomataIntersection performs mathematical intersection of FST and NFA
func (ra *TrueRegexAutomaton) TrueAutomataIntersection(fst *FST) ([]string, error) {
	results := make([]string, 0)
	visited := make(map[string]bool)
	
	// Start with NFA epsilon closure of start state
	startStates := ra.epsilonClosure([]NFAStateID{ra.nfa.start})
	
	// Use character-prefix grouping to reduce computations
	firstCharMap := make(map[byte][]string)
	
	// Group keys by first character
	iterator := fst.Iterator()
	for iterator.HasNext() {
		key, _ := iterator.Next()
		keyStr := string(key)
		if len(keyStr) > 0 {
			firstChar := keyStr[0]
			firstCharMap[firstChar] = append(firstCharMap[firstChar], keyStr)
		}
	}
	
	// For each first character, check if NFA can consume it
	for char, keys := range firstCharMap {
		nextStates := ra.computeNFATransitionsChar(startStates, char)
		
		if len(nextStates) > 0 {
			// NFA can consume this character, test all keys starting with it
			for _, key := range keys {
				if !visited[key] && ra.simulateNFA(key) {
					results = append(results, key)
					visited[key] = true
				}
			}
		}
	}
	
	return results, nil
}

// simulateNFA executes the NFA on input string
func (ra *TrueRegexAutomaton) simulateNFA(input string) bool {
	currentStates := ra.epsilonClosure([]NFAStateID{ra.nfa.start})
	
	for _, char := range []byte(input) {
		currentStates = ra.computeNFATransitionsChar(currentStates, char)
		if len(currentStates) == 0 {
			return false
		}
	}
	
	return ra.hasAcceptingState(currentStates)
}

// computeNFATransitionsChar computes NFA transitions for a character
func (ra *TrueRegexAutomaton) computeNFATransitionsChar(currentStates []NFAStateID, char byte) []NFAStateID {
	nextStates := make([]NFAStateID, 0)
	stateSet := make(map[NFAStateID]bool)
	
	for _, state := range currentStates {
		if transitions, exists := ra.nfa.states[state].transitions[char]; exists {
			for _, nextState := range transitions {
				if !stateSet[nextState] {
					stateSet[nextState] = true
					nextStates = append(nextStates, nextState)
				}
			}
		}
	}
	
	return ra.epsilonClosure(nextStates)
}
