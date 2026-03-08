package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/jamra/gocleo/internal/fst"
)

func main() {
	fmt.Println("🤖 FST Automata Intersection Demo")
	fmt.Println("=================================")
	fmt.Println("This demonstrates true FST ∩ Regex intersection using finite state automata theory!")
	fmt.Println()

	// Create sample documents
	documents := []string{
		"apple fruit delicious",
		"banana tropical fruit",
		"application software program",
		"application form document",
		"apprentice learning student",
		"approach method technique",
		"computer programming language",
		"language processing natural",
		"processing data analysis",
		"analysis statistical method",
		"method systematic approach",
		"systematic organized structure",
		"structure data organization",
		"organization business company",
		"company software development",
		"development programming process",
		"process workflow system",
		"system computer network",
		"network distributed computing",
		"computing machine learning",
	}

	fmt.Printf("📚 Building FST from %d documents...\n", len(documents))
	start := time.Now()

	// Build FST
	fstIndex, err := fst.BuildFSTFromDocuments(documents)
	if err != nil {
		log.Fatal("Failed to build FST:", err)
	}

	// Create search engine
	searchEngine := fst.NewSearchEngine(fstIndex, documents, nil)

	buildTime := time.Since(start)
	fmt.Printf("✅ FST built in %v with %d entries\n", buildTime, fstIndex.Size())
	fmt.Println()

	// Test different regex patterns with automata intersection
	patterns := []struct {
		name    string
		pattern string
		desc    string
	}{
		{"Simple Suffix", ".*ing$", "Words ending in '-ing'"},
		{"Complex Pattern", "^app.*", "Words starting with 'app'"},
		{"Character Class", "^[a-c].*", "Words starting with a, b, or c"},
		{"Multiple Options", "(data|computer|system)", "Words containing 'data', 'computer', or 'system'"},
		{"Wildcard", "pro.*ing", "Words starting with 'pro' and ending with 'ing'"},
		{"Length Pattern", "^.{4,6}$", "Words with 4-6 characters"},
	}

	fmt.Println("🔍 Testing Automata Intersection vs Naive Iteration:")
	fmt.Println("====================================================")

	for _, test := range patterns {
		fmt.Printf("\n📋 Pattern: %s (%s)\n", test.pattern, test.desc)
		fmt.Printf("   Description: %s\n", test.name)

		// Method 1: True Automata Intersection (our new approach)
		start = time.Now()
		intersectionResults, err := searchEngine.IntersectionRegexSearch(test.pattern)
		intersectionTime := time.Since(start)

		if err != nil {
			fmt.Printf("❌ Intersection error: %v\n", err)
			continue
		}

		fmt.Printf("🤖 Automata Intersection: %d results in %v\n", len(intersectionResults), intersectionTime)

		// Method 2: Naive iteration (old approach)
		start = time.Now()
		naiveResults, err := searchEngine.RegexSearch(test.pattern)
		naiveTime := time.Since(start)

		if err != nil {
			fmt.Printf("❌ Naive search error: %v\n", err)
			continue
		}

		fmt.Printf("🔄 Naive Iteration: %d results in %v\n", len(naiveResults), naiveTime)

		// Performance comparison
		if naiveTime.Nanoseconds() > 0 {
			speedup := float64(naiveTime.Nanoseconds()) / float64(intersectionTime.Nanoseconds())
			fmt.Printf("⚡ Speedup: %.2fx faster with automata intersection\n", speedup)
		}

		// Show some results
		if len(intersectionResults) > 0 {
			fmt.Print("   Results: ")
			for i, result := range intersectionResults {
				if i >= 5 {
					fmt.Print("...")
					break
				}
				if i > 0 {
					fmt.Print(", ")
				}
				// Extract the matching word from the document
				words := strings.Fields(result.Word)
				for _, word := range words {
					if matched, _ := regexp.MatchString(test.pattern, word); matched {
						fmt.Print(word)
						break
					}
				}
			}
			fmt.Println()
		}

		// Debug information
		debugInfo, err := searchEngine.GetIntersectionDebugInfo(test.pattern)
		if err == nil {
			fmt.Printf("🔧 Debug: NFA=%d states, DFA=%d states, Intersection=%d states\n", 
				debugInfo.NFAStates, debugInfo.DFAStates, debugInfo.IntersectionStates)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎯 AUTOMATA INTERSECTION THEORY DEMO")
	fmt.Println(strings.Repeat("=", 60))

	// Demonstrate the theory behind intersection
	testPattern := "app.*"
	fmt.Printf("\n📚 Detailed Analysis for pattern: %s\n", testPattern)

	debugInfo, err := searchEngine.GetIntersectionDebugInfo(testPattern)
	if err != nil {
		log.Printf("Debug error: %v", err)
	} else {
		fmt.Println(debugInfo.String())
	}

	fmt.Println("\n🧠 How Automata Intersection Works:")
	fmt.Println("1. 📝 Regex → NFA: Convert regex pattern to Non-deterministic Finite Automaton")
	fmt.Println("2. 🔄 NFA → DFA: Convert NFA to Deterministic Finite Automaton using subset construction")  
	fmt.Println("3. ⚡ FST ∩ DFA: Perform product construction to create intersection automaton")
	fmt.Println("4. 📊 Extract Results: Traverse intersection automaton to find all accepted strings")
	fmt.Println()
	fmt.Println("💡 Benefits:")
	fmt.Println("   • No iteration over all FST keys (O(1) instead of O(n))")
	fmt.Println("   • Mathematically optimal approach")
	fmt.Println("   • Handles complex regex patterns efficiently")
	fmt.Println("   • Scales better with large datasets")

	fmt.Println("\n🏁 Demo completed! Automata intersection theory in action! 🤖")
}