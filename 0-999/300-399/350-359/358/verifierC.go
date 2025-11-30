package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt.
const testcasesCData = `
2 0 5
10 0 0 0 13 16 17 16 0 0 0
15 12 16 0 0 0 0 0 11 17 0 13 11 11 14 0
13 14 7 15 11 14 14 17 15 10 5 19 0 9
10 16 16 18 6 11 19 10 6 3 0
2 0 0
4 4 7 1 1
2 5 0
1 0
3 0 0 4
6 16 12 7 0 0 19
4 0 9 19 8
13 4 7 0 3 0 4 12 10 0 8 0 20 0
18 1 0 0 0 7 1 7 2 0 19 19 8 8 0 0 5 0 7
4 0 0 0 0
17 14 20 6 13 0 1 16 3 15 16 3 9 11 13 0 0 1
14 14 0 2 0 0 2 0 3 12 4 12 8 0 0
11 6 0 15 15 0 4 8 15 13 9 5
16 17 2 18 0 17 0 2 20 0 12 0 10 16 0 0 13
4 7 8 0 0
8 11 18 14 0 12 5 1 12
9 13 15 17 2 7 6 12 20 10
15 14 3 0 18 6 3 17 6 0 18 0 4 0 15 0
19 0 2 15 2 15 10 16 2 11 0 12 0 4 0 14 4 0 0 17
13 6 9 4 17 9 0 9 9 12 2 6 16 4
17 9 0 17 0 0 3 12 6 2 11 0 9 4 14 6 10 0
8 11 0 0 0 0 12 10 0
20 10 20 9 20 11 15 0 11 2 5 12 0 0 11 0 20 17 19 9 4
12 15 0 0 0 0 8 0 0 19 0 20 1
5 0 0 8 0 3
11 19 7 12 15 10 19 18 15 14 0 14
2 14 4
4 16 2 14 0
9 0 5 0 2 20 15 0 17 0
17 1 0 16 16 4 0 0 0 0 0 0 0 19 12 16 5 4
7 0 0 4 12 9 17 0
9 0 8 0 17 0 6 0 1 20
1 15
18 0 7 19 6 4 17 2 14 6 0 8 16 13 1 3 8 0 0
11 0 18 6 4 14 11 19 0 0 6 18
11 5 20 7 20 4 9 17 9 0 0 0
11 0 10 2 0 0 6 0 12 20 8 15
16 7 0 0 7 0 10 12 7 17 0 4 12 3 12 12 0
8 0 2 17 1 13 16 0 11
17 15 14 7 0 0 0 6 17 0 0 0 20 2 4 4 0 8
16 0 10 3 3 11 20 0 9 4 0 20 0 17 2 16 13
14 5 0 0 11 9 0 7 7 11 0 2 0 3 0
13 0 18 1 16 13 4 0 0 1 18 0 13 20
15 0 4 4 0 4 18 14 6 8 11 1 19 13 0 17
16 8 0 11 19 7 17 0 0 11 16 7 20 12 0 19 2
13 11 17 4 0 5 15 5 10 12 0 8 0 1
7 3 0 10 12 0 10 16
20 19 14 19 0 0 1 0 20 0 20 6 0 5 12 15 8 0 0 11 5
8 20 7 0 12 18 0 0 15
12 6 19 15 4 6 0 0 6 8 0 0 9
2 15 2
15 0 14 20 0 0 2 16 9 0 0 6 0 5 14 1
11 1 0 0 20 20 0 0 7 0 0 0
3 19 0 0
14 6 0 10 2 1 0 13 15 9 0 3 0 0 0
5 18 17 0 7 7
11 5 6 1 10 15 11 16 1 0 4 7
9 3 0 19 0 19 2 10 0 0
15 14 0 0 12 0 0 11 10 1 0 5 3 13 3 3
2 0 0
13 4 0 8 16 9 13 12 12 0 9 14 18 17
3 8 18 2
17 16 5 11 7 13 0 6 5 9 4 2 5 0 1 1 17 7
6 7 6 7 8 19 0
1 12
17 1 0 17 9 20 0 14 0 14 0 8 0 13 5 15 0 9
3 5 1 2
18 0 0 0 18 15 19 10 0 0 0 16 0 0 0 0 11 3 20
6 7 0 0 13 1 3
8 1 0 5 0 0 0 0 0
11 20 0 0 7 11 5 0 0 0 18 11
10 9 1 2 16 19 9 19 19 20 17
6 7 20 0 19 0 0
13 0 20 0 3 9 7 15 0 1 0 0 0 16
8 0 1 0 13 18 16 0 0
7 2 8 17 3 5 9 0
20 0 0 3 14 20 0 7 10 7 8 3 0 0 8 3 0 10 18 6 14
4 20 1 13 7
18 20 1 16 8 0 0 13 5 7 0 0 12 18 4 1 0 0 8
18 5 0 2 17 5 0 0 0 0 0 6 16 14 0 0 16 16 0
14 0 5 9 4 3 0 1 19 0 10 19 4 0 0
19 10 1 6 1 0 0 15 4 19 2 16 3 0 10 8 6 0 0 18
4 13 14 8 0
9 0 19 6 0 17 10 0 9 0
5 0 17 17 0 16
19 14 12 0 0 4 0 18 0 14 8 15 0 17 13 6 2 0 8 9
2 9 7
14 8 12 0 0 0 0 0 8 10 18 0 0 0 12
10 0 0 1 13 10 0 4 20 12 8
20 13 20 15 6 12 0 9 18 20 2 0 0 10 0 20 0 2 0 0 16
5 7 0 15 16 4
15 13 12 0 19 13 0 0 0 1 0 16 18 18 10 8
5 9 7 2 0 15
1 6
1 0
`

// Replace oracle-based check with semantic verification that accepts any valid plan.
func verify(vals []int, out string) error {
	// Normalize lines (drop trailing empty lines)
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) != len(vals) {
		return fmt.Errorf("expected %d lines, got %d", len(vals), len(lines))
	}
	// containers
	stack := make([]int, 0)
	queue := make([]int, 0)
	deque := make([]int, 0)
	block := make([]int, 0)

	pushKinds := map[string]bool{
		"pushStack": true,
		"pushQueue": true,
		"pushFront": true,
		"pushBack":  true,
	}
	popKinds := map[string]bool{
		"popStack": true,
		"popQueue": true,
		"popFront": true,
		"popBack":  true,
	}
	containerOf := func(op string) string {
		switch op {
		case "popStack":
			return "stack"
		case "popQueue":
			return "queue"
		case "popFront", "popBack":
			return "deque"
		}
		return ""
	}

	// helpers
	pushFront := func(v int) { deque = append([]int{v}, deque...) }
	pushBack := func(v int) { deque = append(deque, v) }
	popStack := func() (int, bool) {
		if len(stack) == 0 {
			return 0, false
		}
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v, true
	}
	popQueue := func() (int, bool) {
		if len(queue) == 0 {
			return 0, false
		}
		v := queue[0]
		queue = queue[1:]
		return v, true
	}
	popFront := func() (int, bool) {
		if len(deque) == 0 {
			return 0, false
		}
		v := deque[0]
		deque = deque[1:]
		return v, true
	}
	popBack := func() (int, bool) {
		if len(deque) == 0 {
			return 0, false
		}
		v := deque[len(deque)-1]
		deque = deque[:len(deque)-1]
		return v, true
	}

	for i, x := range vals {
		line := strings.TrimSpace(lines[i])
		if x != 0 {
			if !pushKinds[line] {
				return fmt.Errorf("line %d: expected push*, got %q", i+1, line)
			}
			// simulate push
			switch line {
			case "pushStack":
				stack = append(stack, x)
			case "pushQueue":
				queue = append(queue, x)
			case "pushFront":
				pushFront(x)
			case "pushBack":
				pushBack(x)
			}
			block = append(block, x)
		} else {
			// zero: must output k ops, where k = min(3, len(block))
			if line == "0" {
				expected := len(block)
				if expected > 3 {
					expected = 3
				}
				if expected != 0 {
					return fmt.Errorf("line %d: reported 0 pops, but need %d", i+1, expected)
				}
			} else {
				parts := strings.Fields(line)
				var k int
				if _, err := fmt.Sscan(parts[0], &k); err != nil {
					return fmt.Errorf("line %d: bad count: %q", i+1, parts[0])
				}
				expected := len(block)
				if expected > 3 {
					expected = 3
				}
				if k != expected {
					return fmt.Errorf("line %d: expected %d pops, got %d", i+1, expected, k)
				}
				if len(parts) != 1+k {
					return fmt.Errorf("line %d: expected %d pop tokens, got %d", i+1, k, len(parts)-1)
				}
				seenContainer := map[string]bool{}
				popped := make([]int, 0, k)
				for j := 0; j < k; j++ {
					op := parts[1+j]
					if !popKinds[op] {
						return fmt.Errorf("line %d: bad op %q", i+1, op)
					}
					c := containerOf(op)
					if seenContainer[c] {
						return fmt.Errorf("line %d: repeated container via %q", i+1, op)
					}
					seenContainer[c] = true
					var v int
					var ok bool
					switch op {
					case "popStack":
						v, ok = popStack()
					case "popQueue":
						v, ok = popQueue()
					case "popFront":
						v, ok = popFront()
					case "popBack":
						v, ok = popBack()
					}
					if !ok {
						return fmt.Errorf("line %d: pop from empty %s", i+1, op)
					}
					popped = append(popped, v)
				}
				// compute expected top k values from current block
				if k > 0 {
					tmp := append([]int(nil), block...)
					// sort desc
					for a := 0; a < len(tmp); a++ {
						for b := a + 1; b < len(tmp); b++ {
							if tmp[b] > tmp[a] {
								tmp[a], tmp[b] = tmp[b], tmp[a]
							}
						}
					}
					expVals := append([]int(nil), tmp[:k]...)
					// sort both popped and expected asc for multiset compare
					asc := func(s []int) {
						for a := 0; a < len(s); a++ {
							for b := a + 1; b < len(s); b++ {
								if s[b] < s[a] {
									s[a], s[b] = s[b], s[a]
								}
							}
						}
					}
					asc(popped)
					asc(expVals)
					if len(popped) != len(expVals) {
						return fmt.Errorf("line %d: internal size mismatch", i+1)
					}
					for j := 0; j < k; j++ {
						if popped[j] != expVals[j] {
							return fmt.Errorf("line %d failed: expected top %v, got %v", i+1, expVals, popped)
						}
					}
				}
			}
			// clear containers and block for next segment
			stack = stack[:0]
			queue = queue[:0]
			deque = deque[:0]
			block = block[:0]
		}
	}
	return nil
}

func run358C(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

type testCase struct {
	vals []int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesCData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d values, got %d", idx+1, n, len(fields)-1)
		}
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value %d: %w", idx+1, i, err)
			}
			vals[i] = val
		}
		cases = append(cases, testCase{vals: vals})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(tc.vals))
		for i, v := range tc.vals {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		got, err := run358C(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := verify(tc.vals, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nexpected %s\ngot %s\n", idx+1, err, "(any valid sequence)", strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
