package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007

type edge struct {
	to     byte
	vision []byte
}

type state struct {
	cur         byte
	pendingType int8
	pending     string
}

type solveResult struct {
	n     int
	lines []int
}

type testCase struct {
	input string
	desc  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	for i, tc := range tests {
		expect, err := solveCase(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to solve internal test %d (%s): %v\n", i+1, tc.desc, err)
			os.Exit(1)
		}
		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d (%s): %v\ninput:\n%s\noutput:\n%s\n", i+1, tc.desc, err, tc.input, got)
			os.Exit(1)
		}
		if err := compareOutput(expect, got); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.desc, err, tc.input, formatLines(expect.lines), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func compareOutput(expect solveResult, got string) error {
	tokens := strings.Fields(got)
	if len(tokens) != len(expect.lines) {
		return fmt.Errorf("expected %d numbers, got %d", len(expect.lines), len(tokens))
	}
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return fmt.Errorf("line %d: invalid integer %q (%v)", i+1, tok, err)
		}
		norm := int((val%mod + mod) % mod)
		if norm != expect.lines[i] {
			return fmt.Errorf("line %d: expected %d, got %d", i+1, expect.lines[i], norm)
		}
	}
	return nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func solveCase(input string) (solveResult, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return solveResult{}, err
	}
	graph := make([][]edge, n+1)
	for i := 0; i < m; i++ {
		var x, y, k int
		if _, err := fmt.Fscan(reader, &x, &y, &k); err != nil {
			return solveResult{}, err
		}
		seq := make([]byte, k)
		for j := 0; j < k; j++ {
			var v int
			if _, err := fmt.Fscan(reader, &v); err != nil {
				return solveResult{}, err
			}
			seq[j] = byte(v)
		}
		graph[x] = append(graph[x], edge{to: byte(y), vision: seq})
	}
	limit := 2 * n
	ans := make([]int, limit+1)
	curStates := make(map[state]int)
	for start := 1; start <= n; start++ {
		st := state{cur: byte(start), pendingType: 1, pending: string([]byte{byte(start)})}
		curStates[st] = (curStates[st] + 1) % mod
	}
	for length := 1; length < limit && len(curStates) > 0; length++ {
		nextStates := make(map[state]int)
		for st, cnt := range curStates {
			if cnt == 0 {
				continue
			}
			baseQueue := []byte(st.pending)
			queueLen := len(baseQueue)
			lenVision := length
			switch st.pendingType {
			case 1:
				lenVision = length - queueLen
			case -1:
				lenVision = length + queueLen
			}
			if lenVision < 0 || lenVision > limit {
				continue
			}
			for _, e := range graph[int(st.cur)] {
				if lenVision+len(e.vision) > limit {
					continue
				}
				queue := append([]byte(nil), baseQueue...)
				ptype := st.pendingType
				ok := true
				for _, val := range e.vision {
					if ptype == 1 {
						if len(queue) == 0 || queue[0] != val {
							ok = false
							break
						}
						queue = queue[1:]
						if len(queue) == 0 {
							ptype = 0
						}
					} else {
						queue = append(queue, val)
						ptype = -1
					}
				}
				if !ok {
					continue
				}
				dest := e.to
				if ptype == -1 {
					if len(queue) == 0 || queue[0] != dest {
						continue
					}
					queue = queue[1:]
					if len(queue) == 0 {
						ptype = 0
					}
				} else {
					queue = append(queue, dest)
					ptype = 1
				}
				newState := state{cur: dest, pendingType: ptype, pending: string(queue)}
				newCnt := nextStates[newState] + cnt
				if newCnt >= mod {
					newCnt -= mod
				}
				nextStates[newState] = newCnt
				if ptype == 0 {
					ans[length+1] += cnt
					if ans[length+1] >= mod {
						ans[length+1] -= mod
					}
				}
			}
		}
		curStates = nextStates
	}
	res := solveResult{n: n, lines: make([]int, limit)}
	copy(res.lines, ans[1:])
	return res, nil
}

func formatLines(vals []int) string {
	var sb strings.Builder
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{
		input: sampleInput(),
		desc:  "sample",
	})
	tests = append(tests, testCase{
		input: "1 0\n",
		desc:  "single node no edges",
	})
	tests = append(tests, testCase{
		input: "2 1\n1 2 1 1\n",
		desc:  "single edge simple vision",
	})
	tests = append(tests, testCase{
		input: "3 2\n1 2 2 1 2\n2 3 1 3\n",
		desc:  "chain vision prefix",
	})
	tests = append(tests, testCase{
		input: "4 3\n1 2 1 1\n2 3 2 2 2\n3 4 1 4\n",
		desc:  "repeated node vision",
	})
	rng := rand.New(rand.NewSource(3240324))
	for i := 0; i < 40; i++ {
		n := rng.Intn(6) + 2
		tests = append(tests, testCase{
			input: randomCase(rng, n),
			desc:  fmt.Sprintf("random-%d", i+1),
		})
	}
	tests = append(tests, testCase{
		input: denseCase(8),
		desc:  "dense-8",
	})
	tests = append(tests, testCase{
		input: sparseCycleCase(),
		desc:  "cycle",
	})
	return tests
}

func sampleInput() string {
	return "6 6\n1 2 2 1 2\n2 3 1 3\n3 4 2 4 5\n4 5 0\n5 3 1 3\n6 1 1 6\n"
}

func randomCase(rng *rand.Rand, n int) string {
	type pair struct{ a, b int }
	var pairs []pair
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			pairs = append(pairs, pair{i, j})
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})
	maxEdges := len(pairs)
	m := 0
	if maxEdges > 0 {
		m = rng.Intn(maxEdges + 1)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		p := pairs[i]
		x, y := p.a, p.b
		if rng.Intn(2) == 0 {
			x, y = y, x
		}
		k := rng.Intn(n + 1)
		fmt.Fprintf(&sb, "%d %d %d", x, y, k)
		for j := 0; j < k; j++ {
			fmt.Fprintf(&sb, " %d", rng.Intn(n)+1)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func denseCase(n int) string {
	type pair struct{ a, b int }
	var edges []pair
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			edges = append(edges, pair{i, j})
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	flag := false
	for _, e := range edges {
		x, y := e.a, e.b
		if flag {
			x, y = y, x
		}
		flag = !flag
		fmt.Fprintf(&sb, "%d %d %d", x, y, 2)
		fmt.Fprintf(&sb, " %d %d\n", x, y)
	}
	return sb.String()
}

func sparseCycleCase() string {
	return "5 5\n1 2 1 1\n2 3 0\n3 4 2 3 4\n4 5 1 5\n5 1 1 5\n"
}
