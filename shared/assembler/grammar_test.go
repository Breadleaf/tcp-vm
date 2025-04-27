package assembler

import (
	"fmt"
	"testing"
)

func Test_parseGrammarText(t *testing.T) {
	_, err := newGrammar()
	if err != nil {
		t.Fatalf("newGrammar() failed: %v", err)
	}
}

func Test_derivesToLambda(t *testing.T) {
	g, err := newGrammar()
	if err != nil {
		t.Fatalf("newGrammar() failed: %v", err)
	}

	fmt.Println("Testing Non Terminals")
	for nt, _ := range g.NonTerminals {
		fmt.Printf("derivesToLambda(%s): %+v\n", nt.Value, g.derivesToLambda(nt))
	}

	fmt.Println("Testing Terminals")
	for t, _ := range g.Terminals {
		fmt.Printf("derivesToLambda(%s): %+v\n", t.Value, g.derivesToLambda(t))
	}
}

func Test_firstSet(t *testing.T) {
	g, err := newGrammar()
	if err != nil {
		t.Fatalf("newGrammar() failed: %v", err)
	}

	fmt.Println("Testing Non Terminals")
	for nt, _ := range g.NonTerminals {
		fmt.Printf("first(%s): %+v\n", nt.Value, g.firstSet(nt))
	}

	fmt.Println("Testing Terminals")
	for t, _ := range g.Terminals {
		fmt.Printf("first(%s): %+v\n", t.Value, g.firstSet(t))
	}
}

func Test_followSet(t *testing.T) {
	g, err := newGrammar()
	if err != nil {
		t.Fatalf("newGrammar() failed: %v", err)
	}

	fmt.Println("Testing Non Terminals")
	for nt, _ := range g.NonTerminals {
		fmt.Printf("follow(%s): %+v\n", nt.Value, g.followSet(nt))
	}

	fmt.Println("Testing Terminals")
	for t, _ := range g.Terminals {
		fmt.Printf("follow(%s): %+v\n", t.Value, g.followSet(t))
	}
}

func Test_predictSet(t *testing.T) {
	g, err := newGrammar()
	if err != nil {
		t.Fatalf("newGrammar() failed: %v", err)
	}

	for nt, prod := range g.Rules {
		for _, rule := range prod {
			fmt.Printf(
				"predict(%+v -> %+v): %+v\n",
				nt,
				rule,
				g.predictSet(nt, rule),
			)
		}
	}
}

func Test_verifyPredictPairwiseDisjoint(t *testing.T) {
	g, err := newGrammar()
	if err != nil {
		t.Fatalf("newGrammar() failed: %v", err)
	}

	for nt, _ := range g.NonTerminals {
		fmt.Printf("Checking if %s is pairwise disjoint...\n", nt)
		if g.verifyPredictPairwiseDisjoint(nt) {
			fmt.Printf("%s is pairwise disjoint\n", nt)
		} else {
			fmt.Printf("%s is NOT pairwise disjoint\n", nt)
		}
	}
}
