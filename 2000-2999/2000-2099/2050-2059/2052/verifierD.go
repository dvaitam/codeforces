package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource       = "2000-2999/2000-2099/2050-2059/2052/2052D.go"
	randomCaseLimit = 120
	totalOpLimit    = 4000
)

type operation struct {
	typ string
	res bool
}

type testCase struct {
	ops   []operation
	edges [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}
		refOrder, refNegOne, err := parseAnswer(refOut, len(tc.ops))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}
		if !refNegOne {
			if err := validateOrder(tc, refOrder); err != nil {
				fmt.Fprintf(os.Stderr, "reference produced invalid order on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
				os.Exit(1)
			}
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}
		candOrder, candNegOne, err := parseAnswer(candOut, len(tc.ops))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		if refNegOne {
			if !candNegOne {
				fmt.Fprintf(os.Stderr, "test %d expected -1, candidate produced an order\n", idx+1)
				os.Exit(1)
			}
			continue
		}

		if candNegOne {
			fmt.Fprintf(os.Stderr, "test %d has a valid order but candidate printed -1\n", idx+1)
			os.Exit(1)
		}
		if err := validateOrder(tc, candOrder); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid candidate order: %v\n", idx+1, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2052D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func parseAnswer(out string, n int) ([]int, bool, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return nil, false, fmt.Errorf("empty output")
	}
	if tokens[0] == "-1" {
		if len(tokens) > 1 {
			return nil, false, fmt.Errorf("extra tokens after -1")
		}
		return nil, true, nil
	}
	if len(tokens) != n {
		return nil, false, fmt.Errorf("expected %d numbers, got %d", n, len(tokens))
	}
	order := make([]int, n)
	for i, t := range tokens {
		val, err := strconv.Atoi(t)
		if err != nil {
			return nil, false, fmt.Errorf("failed to parse number %q: %v", t, err)
		}
		order[i] = val
	}
	return order, false, nil
}

func validateOrder(tc testCase, order []int) error {
	n := len(tc.ops)
	if len(order) != n {
		return fmt.Errorf("order length mismatch: got %d, expected %d", len(order), n)
	}
	pos := make([]int, n+1)
	seen := make([]bool, n+1)
	for i, v := range order {
		if v < 1 || v > n {
			return fmt.Errorf("position %d uses invalid index %d", i+1, v)
		}
		if seen[v] {
			return fmt.Errorf("duplicate index %d in order", v)
		}
		seen[v] = true
		pos[v] = i
	}
	for i := 1; i <= n; i++ {
		if !seen[i] {
			return fmt.Errorf("missing index %d in order", i)
		}
	}

	for _, e := range tc.edges {
		a, b := e[0], e[1]
		if pos[a] >= pos[b] {
			return fmt.Errorf("edge %d -> %d violated", a, b)
		}
	}

	state := false
	for i, idx := range order {
		op := tc.ops[idx-1]
		var got bool
		if op.typ == "set" {
			if state {
				got = false
			} else {
				got = true
				state = true
			}
		} else { // unset
			if state {
				got = true
				state = false
			} else {
				got = false
			}
		}
		if got != op.res {
			return fmt.Errorf("operation %d (index %d) expected %v, got %v", i+1, idx, op.res, got)
		}
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	total := 0

	add := func(tc testCase) {
		if len(tc.ops) == 0 {
			return
		}
		if total+len(tc.ops) > totalOpLimit {
			return
		}
		tests = append(tests, tc)
		total += len(tc.ops)
	}

	// Statement samples.
	add(manualCase(
		[]operation{
			{"set", true}, {"unset", true}, {"set", false}, {"unset", false}, {"unset", false},
		},
		[][2]int{{1, 4}, {5, 2}},
	))
	add(manualCase(
		[]operation{{"unset", true}, {"unset", false}, {"set", true}},
		nil,
	))
	add(manualCase(
		[]operation{{"unset", false}, {"set", true}},
		[][2]int{{2, 1}},
	))
	add(manualCase(
		[]operation{{"unset", false}, {"set", false}},
		nil,
	))

	// Small crafted impossible case: two set true without a reset.
	add(manualCase(
		[]operation{{"set", true}, {"set", true}},
		nil,
	))

	// Small valid case with no true operations.
	add(manualCase(
		[]operation{{"unset", false}, {"unset", false}, {"unset", false}},
		[][2]int{{1, 2}, {2, 3}},
	))

	// Random cases.
	for attempts := 0; attempts < 200 && total < totalOpLimit; attempts++ {
		n := rng.Intn(randomCaseLimit) + 1
		if total+n > totalOpLimit {
			n = totalOpLimit - total
		}
		if n <= 0 {
			break
		}
		if rng.Intn(100) < 65 {
			add(randomValidCase(n, rng))
		} else {
			add(randomInvalidCase(n, rng))
		}
	}

	if len(tests) == 0 {
		add(randomValidCase(1, rng))
	}
	return tests
}

func manualCase(ops []operation, edges [][2]int) testCase {
	cpyOps := make([]operation, len(ops))
	copy(cpyOps, ops)
	cpyEdges := make([][2]int, len(edges))
	copy(cpyEdges, edges)
	return testCase{ops: cpyOps, edges: cpyEdges}
}

func randomValidCase(n int, rng *rand.Rand) testCase {
	trueCnt := rng.Intn(3)
	// Ensure feasibility with the register model.
	order := make([]int, n)
	for i := range order {
		order[i] = i + 1
	}
	rng.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

	ops := make([]operation, n)
	switch trueCnt {
	case 0:
		for i := range ops {
			ops[i] = operation{typ: "unset", res: false}
		}
	case 1:
		pos := rng.Intn(n)
		for i := 0; i < pos; i++ {
			ops[i] = operation{typ: "unset", res: false}
		}
		ops[pos] = operation{typ: "set", res: true}
		for i := pos + 1; i < n; i++ {
			ops[i] = operation{typ: "set", res: false}
		}
	default:
		pos1 := rng.Intn(n)
		pos2 := rng.Intn(n)
		for pos2 == pos1 {
			pos2 = rng.Intn(n)
		}
		if pos1 > pos2 {
			pos1, pos2 = pos2, pos1
		}
		for i := 0; i < pos1; i++ {
			ops[i] = operation{typ: "unset", res: false}
		}
		ops[pos1] = operation{typ: "set", res: true}
		for i := pos1 + 1; i < pos2; i++ {
			ops[i] = operation{typ: "set", res: false}
		}
		ops[pos2] = operation{typ: "unset", res: true}
		for i := pos2 + 1; i < n; i++ {
			ops[i] = operation{typ: "unset", res: false}
		}
	}

	edges := randomEdges(order, rng, 0.35)
	return testCase{ops: ops, edges: edges}
}

func randomInvalidCase(n int, rng *rand.Rand) testCase {
	choice := rng.Intn(3)
	switch choice {
	case 0:
		// No set true operations, but several operations require state true.
		ops := make([]operation, n)
		for i := range ops {
			if rng.Intn(100) < 50 {
				ops[i] = operation{typ: "set", res: false}
			} else {
				ops[i] = operation{typ: "unset", res: false}
			}
		}
		order := naturalOrder(n)
		return testCase{ops: ops, edges: randomEdges(order, rng, 0.4)}
	case 1:
		// Two set true operations without a way to return to false.
		ops := make([]operation, n)
		truePlaced := 0
		for i := range ops {
			if truePlaced < 2 && rng.Intn(100) < 30 {
				ops[i] = operation{typ: "set", res: true}
				truePlaced++
			} else {
				ops[i] = operation{typ: "set", res: false}
			}
		}
		for truePlaced < 2 && truePlaced < n {
			ops[truePlaced] = operation{typ: "set", res: true}
			truePlaced++
		}
		order := naturalOrder(n)
		return testCase{ops: ops, edges: randomEdges(order, rng, 0.25)}
	default:
		// Start from a valid case and add an edge that breaks feasibility.
		base := randomValidCase(n, rng)
		setPos := -1
		unsetTruePos := -1
		for idx, op := range base.ops {
			if op.typ == "set" && op.res {
				setPos = idx
			}
			if op.typ == "unset" && op.res {
				unsetTruePos = idx
			}
		}
		edges := make([][2]int, len(base.edges))
		copy(edges, base.edges)
		if setPos != -1 {
			if unsetTruePos != -1 && rng.Intn(2) == 0 {
				// Force unset true before the initial set true.
				edges = append(edges, [2]int{unsetTruePos + 1, setPos + 1})
			} else {
				// Pick an operation that requires state true and force it before the first set true.
				target := -1
				for i, op := range base.ops {
					if i == setPos {
						continue
					}
					if op.typ == "set" && !op.res {
						target = i
						break
					}
				}
				if target == -1 {
					target = 0
				}
				edges = append(edges, [2]int{target + 1, setPos + 1})
			}
		} else {
			// Fallback impossible case without true operations.
			ops := make([]operation, n)
			for i := range ops {
				ops[i] = operation{typ: "set", res: false}
			}
			order := naturalOrder(n)
			return testCase{ops: ops, edges: randomEdges(order, rng, 0.3)}
		}
		return testCase{ops: base.ops, edges: edges}
	}
}

func randomEdges(order []int, rng *rand.Rand, prob float64) [][2]int {
	n := len(order)
	index := make([]int, n+1)
	for pos, v := range order {
		index[v] = pos
	}
	var edges [][2]int
	seen := make(map[[2]int]struct{})
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Float64() > prob {
				continue
			}
			a := order[i]
			b := order[j]
			key := [2]int{a, b}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			edges = append(edges, key)
		}
	}
	return edges
}

func naturalOrder(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	return arr
}

func buildInput(tc testCase) string {
	var b strings.Builder
	n := len(tc.ops)
	fmt.Fprintf(&b, "%d\n", n)
	for _, op := range tc.ops {
		resStr := "false"
		if op.res {
			resStr = "true"
		}
		fmt.Fprintf(&b, "%s %s\n", op.typ, resStr)
	}
	fmt.Fprintf(&b, "%d\n", len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
	}
	return b.String()
}
