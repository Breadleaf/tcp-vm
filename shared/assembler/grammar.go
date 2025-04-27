package assembler

import (
	"slices"
	"tcp-vm/shared/util"

	"bufio"
	"fmt"
	"os"
	"strings"
)

const FILE_LOG_TAG = "tcp-vm/shared/assembler - grammar.go"

type grammarType int

const (
	NonTerminal grammarType = iota
	Terminal
	Lambda
)

func (gt grammarType) String() string {
	switch gt {
	case NonTerminal:
		return "NonTerminal"
	case Terminal:
		return "Terminal"
	case Lambda:
		return "Lambda"
	}
	return "impossible"
}

type grammarItem struct {
	Type  grammarType
	Value string
}

func (gi *grammarItem) Equal(gi2 *grammarItem) bool {
	return (gi.Type == gi2.Type &&
		gi.Value == gi2.Value)
}

type grammar struct {
	Rules        map[grammarItem][][]grammarItem
	NonTerminals util.Set[grammarItem]
	Terminals    util.Set[grammarItem]
}

func newGrammar() (*grammar, error) {
	g := grammar{
		Rules:        make(map[grammarItem][][]grammarItem),
		NonTerminals: util.NewSet[grammarItem](),
		Terminals:    util.NewSet[grammarItem](),
	}

	if err := g.parseGrammarText(); err != nil {
		return nil, fmt.Errorf("error parsing grammar text: %v", err)
	}

	return &g, nil
}

func (g *grammar) parseGrammarText() error {
	logTag := fmt.Sprintf("%s - parseGrammarText()", FILE_LOG_TAG)
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	reader := strings.NewReader(grammarText)
	scanner := bufio.NewScanner(reader)

	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading grammar text: %v", err)
	}

	util.LogMessage(func() {
		fmt.Println("Parsed from grammar text:")
		fmt.Println("lines:")
		for _, line := range lines {
			fmt.Printf("'%+v'\n", line)
		}
	})

	// model of a non fully parsed grammar
	// each value in groups will represent a rule
	nonTerminals := util.NewSet[string]()
	groups := make(map[string][]string)

	for _, line := range lines {
		parts := strings.SplitN(line, "->", 2)
		if len(parts) != 2 {
			return fmt.Errorf("parsing line: %s: missing '->'", line)
		}
		lhs := strings.TrimSpace(parts[0])
		rhs := strings.TrimSpace(parts[1])
		nonTerminals.Add(lhs)
		groups[lhs] = append(groups[lhs], rhs)
	}

	util.LogMessage(func() {
		fmt.Println("Parsed file prior to grammar conversion:")
		fmt.Printf("non terminals: %+v\n", nonTerminals)
		fmt.Println("groups:")
		for lhs, group := range groups {
			fmt.Printf("%s: [", lhs)
			for i, s := range group {
				fmt.Printf("'%s'", s)
				if i < len(group)-1 {
					fmt.Printf(" ")
				}
			}
			fmt.Println("]")
		}
	})

	for ntName, rhsList := range groups {
		ntItem := grammarItem{Type: NonTerminal, Value: ntName}

		var rules [][]grammarItem
		for _, rhs := range rhsList {
			tokStrings := strings.Fields(rhs)
			var rule []grammarItem
			for _, tk := range tokStrings {
				typ := Terminal
				if nonTerminals.Contains(tk) {
					typ = NonTerminal
				} else if tk == "lambda" {
					typ = Lambda
				}
				rule = append(rule, grammarItem{Type: typ, Value: tk})
			}
			rules = append(rules, rule)
		}

		g.Rules[ntItem] = rules
		g.NonTerminals.Add(ntItem)
	}

	for _, rhs := range g.Rules {
		for _, rule := range rhs {
			for _, el := range rule {
				if !g.NonTerminals.Contains(el) {
					g.Terminals.Add(el)
				}
			}
		}
	}

	if LOG_PARSED_GRAMMAR_OBJECT {
		util.LogMessage(func() {
			fmt.Println("Parsed file post grammar conversion:")
			fmt.Printf("%+v\n", g)
		})
	}

	return nil
}

func (g *grammar) derivesToLambda(nonTerminal grammarItem) bool {
	logTag := fmt.Sprintf("%s - derivesToLambda()", FILE_LOG_TAG)
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	util.LogMessage(func() {
		fmt.Printf(
			"Does %+v derive to lambda with rules:\n",
			nonTerminal,
		)
		for nt, rules := range g.Rules {
			for _, rule := range rules {
				fmt.Printf("%s: %v\n", nt, rule)
			}
		}
	})

	T := util.Stack[[]grammarItem]{}

	var helper func(grammarItem) bool
	helper = func(nt grammarItem) bool {
		for _, p := range g.Rules[nt] {
			if T.Contains(p) {
				continue
			}

			if len(p) == 1 && p[0].Type == Lambda {
				return true
			}

			// if anything in p is terminal then continue
			anythingTerminal := false
			for _, el := range p {
				if el.Type == Terminal {
					anythingTerminal = true
					break
				}
			}
			if anythingTerminal {
				continue
			}

			allDeriveLambda := true
			for _, Xi := range p {
				if Xi.Type == Terminal {
					continue
				}

				T.Push(p)

				allDeriveLambda = helper(Xi)

				_, err := T.Pop()
				if err != nil {
					fmt.Printf(
						"%s - unrecoverable state: %v\n",
						logTag,
						err,
					)
					os.Exit(1)
				}

				if !allDeriveLambda {
					break
				}
			}

			if allDeriveLambda {
				return true
			}
		}

		return false
	}

	return helper(nonTerminal)
}

func (g *grammar) firstSet(nonTerminal grammarItem) util.Set[grammarItem] {
	logTag := fmt.Sprintf("%s - firstSet()", FILE_LOG_TAG)
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	memo := make(map[grammarItem]util.Set[grammarItem])

	var first func(grammarItem) util.Set[grammarItem]
	first = func(gi grammarItem) util.Set[grammarItem] {
		if fs, ok := memo[gi]; ok {
			return fs
		}

		result := util.NewSet[grammarItem]()
		memo[gi] = result

		switch gi.Type {
		case NonTerminal:
			firstNT := func(
				seq []grammarItem,
			) util.Set[grammarItem] {
				s := util.NewSet[grammarItem]()

				for _, sym := range seq {
					// add terminal symbols <- first(sym)
					for f, _ := range first(sym) {
						if f.Type == Terminal {
							s.Add(f)
						}
					}

					// if symbol doesnt derive to lambda:
					// - we are done
					if !g.derivesToLambda(sym) {
						break
					}

					// continue to next symbol
				}

				return s
			}

			for _, prod := range g.Rules[gi] {
				if len(prod) == 1 && prod[0].Type == Lambda {
					continue
				}
				for f, _ := range firstNT(prod) {
					result.Add(f)
				}
			}
		case Terminal:
			result.Add(gi)
		case Lambda:
			// lambda -> {}
		}

		return result
	}

	return first(nonTerminal)
}

func (g *grammar) followSet(nonTerminal grammarItem) util.Set[grammarItem] {
	logTag := fmt.Sprintf("%s - followSet()", FILE_LOG_TAG)
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	if nonTerminal.Type != NonTerminal {
		return util.NewSet[grammarItem]()
	}

	memo := make(map[string]util.Set[grammarItem])
	var follow func(grammarItem, util.Set[grammarItem]) util.Set[grammarItem]
	follow = func(
		gi grammarItem,
		visited util.Set[grammarItem],
	) util.Set[grammarItem] {
		memoKey := gi.Value
		if cached, exists := memo[memoKey]; exists && !visited.Contains(gi) {
			return cached
		}

		if visited.Contains(gi) {
			return util.NewSet[grammarItem]()
		}

		newVisited := util.NewSet[grammarItem]()
		for v := range visited {
			newVisited.Add(v)
		}
		newVisited.Add(gi)

		f := util.NewSet[grammarItem]()

		for lhs, productions := range g.Rules {
			for _, prod := range productions {
				for i, item := range prod {
					// check if we found our non-terminal in the production
					if item.Equal(&gi) {
						// if it's not the last symbol in the production
						if i < len(prod)-1 {
							// calculate FIRST of everything that follows
							j := i + 1
							beta := prod[j:]

							// process each symbol in beta
							allDeriveToLambda := true
							for _, symbol := range beta {
								if symbol.Type == Terminal {
									f.Add(symbol)
									allDeriveToLambda = false
									break
								} else if symbol.Type == NonTerminal {
									// add all terminals from FIRST(symbol) to follow set
									firstSet := g.firstSet(symbol)
									for t := range firstSet {
										if t.Type == Terminal {
											f.Add(t)
										}
									}

									// check if this symbol can derive lambda
									if !g.derivesToLambda(symbol) {
										allDeriveToLambda = false
										break
									}
								}
							}

							// if all symbols in beta can derive lambda, add FOLLOW(lhs)
							if allDeriveToLambda && !lhs.Equal(&gi) {
								followLHS := follow(lhs, newVisited)
								for t := range followLHS {
									f.Add(t)
								}
							}
						} else {
							// if gi is the last symbol in the production
							if !lhs.Equal(&gi) {
								followLHS := follow(lhs, newVisited)
								for t := range followLHS {
									f.Add(t)
								}
							}
						}
					}
				}
			}
		}

		memo[memoKey] = f
		return f
	}

	return follow(nonTerminal, util.NewSet[grammarItem]())
}

func (g *grammar) predictSet(
	lhs grammarItem,
	rhs []grammarItem,
) util.Set[grammarItem] {
	logTag := fmt.Sprintf("%s - predictSet()", FILE_LOG_TAG)
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	// special case: if the production is A -> lambda
	if len(rhs) == 1 && rhs[0].Type == Lambda {
		return g.followSet(lhs)
	}

	result := util.NewSet[grammarItem]()

	// process each symbol in the rhs
	allDeriveToLambda := true
	for i, sym := range rhs {
		if sym.Type == Terminal {
			// if we find a terminal symbol, add it to the predict set and return
			// (only if it's the first symbol or all previous symbols can derive lambda)
			if i == 0 || allDeriveToLambda {
				result.Add(sym)
				return result
			}
			allDeriveToLambda = false
			break
		} else if sym.Type == NonTerminal {
			// add all terminals from FIRST(sym) to the predict set
			firstSet := g.firstSet(sym)
			for t := range firstSet {
				if t.Type == Terminal {
					result.Add(t)
				}
			}

			// check if this symbol can derive lambda
			if !g.derivesToLambda(sym) {
				allDeriveToLambda = false
				break
			}
		}
	}

	// if all symbols in rhs can derive lambda, add FOLLOW(lhs)
	if allDeriveToLambda {
		followSet := g.followSet(lhs)
		for t := range followSet {
			result.Add(t)
		}
	}

	return result
}

// verify all predict sets for given NT are pairwise disjoint
func (g *grammar) verifyPredictPairwiseDisjoint(gi grammarItem) bool {
	logTag := fmt.Sprintf("%s - verifyPredictPairwiseDisjoint()", FILE_LOG_TAG)
	util.LogStart(logTag)
	defer util.LogEnd(logTag)

	predictSets := []util.Set[grammarItem]{}
	for _, rule := range g.Rules[gi] {
		predictSets = append(predictSets, g.predictSet(gi, rule))
	}

	seen := []grammarItem{}
	for _, s := range predictSets {
		for el, _ := range s {
			if slices.Contains(seen, el) {
				return false
			}
			seen = append(seen, el)
		}
	}

	return true
}
