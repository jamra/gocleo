package main

import (
	"fmt"
	"testing"

	"github.com/jamra/gocleo/internal/fst"
)

var (
	benchmarkFST    *fst.FST
	benchmarkEngine *fst.SearchEngine
	benchmarkDocs   []string
)

func init() {
	// Create a larger dataset for meaningful benchmarks
	benchmarkDocs = []string{
		"apple fruit delicious red sweet",
		"banana yellow tropical fruit",
		"application software program development",
		"application form document business",
		"apprentice learning student education",
		"approach method systematic technique",
		"computer programming language technology",
		"language processing natural analysis",
		"processing data computation analysis",
		"analysis statistical mathematical method",
		"method systematic organized approach",
		"systematic structured organized process",
		"structure data organization framework",
		"organization business company enterprise",
		"company software development technology",
		"development programming coding process",
		"process workflow management system",
		"system computer network infrastructure",
		"network distributed computing cloud",
		"computing machine learning artificial",
		"machine learning algorithm intelligence",
		"algorithm computational mathematics logic",
		"mathematics numerical statistical analysis",
		"statistics probability mathematical models",
		"models predictive analytical frameworks",
		"framework architectural software design",
		"design pattern implementation structure",
		"implementation coding development practice",
		"practice methodology systematic approach",
		"methodology research scientific analysis",
		"research investigation study exploration",
		"study academic educational learning",
		"education training knowledge development",
		"knowledge information data wisdom",
		"information technology digital systems",
		"technology innovation advancement progress",
		"innovation creativity invention discovery",
		"creativity artistic imaginative expression",
		"expression communication language interaction",
		"communication messaging information transfer",
		"transfer movement transportation logistics",
		"transportation vehicle automotive systems",
		"vehicle automobile mechanical engineering",
		"engineering technical scientific application",
		"technical expertise specialized knowledge",
		"expertise professional competency skills",
		"professional career workplace employment",
		"career development growth advancement",
		"growth expansion increase improvement",
		"improvement enhancement optimization progress",
		"optimization efficiency performance tuning",
		"performance speed execution measurement",
	}

	var err error
	benchmarkFST, err = fst.BuildFSTFromDocuments(benchmarkDocs)
	if err != nil {
		panic(err)
	}

	benchmarkEngine = fst.NewSearchEngine(benchmarkFST, benchmarkDocs, nil)
}

// BenchmarkIntersectionVsNaive compares automata intersection with naive iteration
func BenchmarkIntersectionVsNaive(b *testing.B) {
	patterns := []struct {
		name    string
		pattern string
	}{
		{"SimpleSuffix", ".*ing$"},
		{"ComplexPrefix", "^app.*"},
		{"CharClass", "^[a-d].*"},
		{"Alternation", "(data|computer|system)"},
		{"Wildcard", "pro.*ing"},
	}

	for _, pattern := range patterns {
		b.Run("Intersection_"+pattern.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := benchmarkEngine.IntersectionRegexSearch(pattern.pattern)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run("Naive_"+pattern.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := benchmarkEngine.RegexSearch(pattern.pattern)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkAutomataConstruction measures the cost of building automata
func BenchmarkAutomataConstruction(b *testing.B) {
	patterns := []string{
		"simple",
		".*ing$",
		"^app.*",
		"(data|computer|system)",
		"pro.*ing.*tion",
	}

	for _, pattern := range patterns {
		b.Run("NFA_"+pattern, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				nfa, err := fst.RegexToNFA(pattern)
				if err != nil {
					b.Fatal(err)
				}
				_ = nfa
			}
		})

		b.Run("DFA_"+pattern, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				nfa, err := fst.RegexToNFA(pattern)
				if err != nil {
					b.Fatal(err)
				}
				dfa := fst.NFAtoDFA(nfa)
				_ = dfa
			}
		})
	}
}

// BenchmarkComplexPatterns tests performance on increasingly complex regex patterns
func BenchmarkComplexPatterns(b *testing.B) {
	complexPatterns := []struct {
		name    string
		pattern string
	}{
		{"Simple", "app"},
		{"SimpleWildcard", "app.*"},
		{"MultipleClauses", "(app|dev|sys).*"},
		{"NestedGroups", "((app|dev)|(sys|net)).*"},
		{"CharClasses", "[a-z]{3,8}"},
		{"ComplexCombination", "^[a-d].*(ing|tion|ment)$"},
	}

	for _, test := range complexPatterns {
		b.Run("Intersection_"+test.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				results, err := benchmarkEngine.IntersectionRegexSearch(test.pattern)
				if err != nil {
					b.Fatal(err)
				}
				_ = results
			}
		})

		b.Run("Naive_"+test.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				results, err := benchmarkEngine.RegexSearch(test.pattern)
				if err != nil {
					b.Fatal(err)
				}
				_ = results
			}
		})
	}
}

// BenchmarkScalability tests how well intersection scales with FST size
func BenchmarkScalability(b *testing.B) {
	sizes := []int{10, 50, 100}
	
	for _, size := range sizes {
		// Create FST of specified size
		docs := benchmarkDocs[:min(size, len(benchmarkDocs))]
		testFST, err := fst.BuildFSTFromDocuments(docs)
		if err != nil {
			b.Fatal(err)
		}
		testEngine := fst.NewSearchEngine(testFST, docs, nil)

		b.Run(fmt.Sprintf("Intersection_Size%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := testEngine.IntersectionRegexSearch(".*ing$")
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(fmt.Sprintf("Naive_Size%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := testEngine.RegexSearch(".*ing$")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}