package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Suffix automaton for counting distinct substrings over {O, X}.
type samState struct {
	len, link int
	next      [2]int
}

func countDistinctSubstrings(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	states := make([]samState, 2*n+5)
	for i := range states {
		states[i].next[0] = -1
		states[i].next[1] = -1
		states[i].link = -1
	}
	sz, last := 1, 0
	for _, c := range s {
		ci := 0
		if c == 'X' {
			ci = 1
		}
		cur := sz
		sz++
		states[cur].len = states[last].len + 1
		p := last
		for p != -1 && states[p].next[ci] == -1 {
			states[p].next[ci] = cur
			p = states[p].link
		}
		if p == -1 {
			states[cur].link = 0
		} else {
			q := states[p].next[ci]
			if states[p].len+1 == states[q].len {
				states[cur].link = q
			} else {
				clone := sz
				sz++
				states[clone] = states[q]
				states[clone].len = states[p].len + 1
				for p != -1 && states[p].next[ci] == q {
					states[p].next[ci] = clone
					p = states[p].link
				}
				states[q].link = clone
				states[cur].link = clone
			}
		}
		last = cur
	}
	res := 0
	for i := 1; i < sz; i++ {
		res += states[i].len - states[states[i].link].len
	}
	return res
}

func runTest(bin string, n int) error {
	cmd := exec.Command(bin)
	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %v", err)
	}
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %v", err)
	}

	reader := bufio.NewReader(outPipe)

	// Send n.
	fmt.Fprintf(inPipe, "%d\n", n)

	// Read n magic words.
	words := make([]string, n)
	seen := make(map[string]bool)
	for i := 0; i < n; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading word %d: %v", i+1, err)
		}
		w := strings.TrimSpace(line)
		if len(w) == 0 {
			return fmt.Errorf("word %d is empty", i+1)
		}
		for _, c := range w {
			if c != 'X' && c != 'O' {
				return fmt.Errorf("word %d contains invalid char %c", i+1, c)
			}
		}
		if seen[w] {
			return fmt.Errorf("duplicate word %q at position %d", w, i+1)
		}
		seen[w] = true
		words[i] = w
	}

	// Compute spell power for every ordered pair.
	type pair struct{ i, j int }
	powers := make(map[int]pair)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			p := countDistinctSubstrings(words[i] + words[j])
			if prev, exists := powers[p]; exists {
				return fmt.Errorf("power %d is not unique: (%d,%d) and (%d,%d) both give it",
					p, prev.i, prev.j, i+1, j+1)
			}
			powers[p] = pair{i + 1, j + 1}
		}
	}

	// Send query count, then interleave each query with its answer.
	fmt.Fprintf(inPipe, "%d\n", n*n)
	for power, exp := range powers {
		fmt.Fprintf(inPipe, "%d\n", power)
		var a, b int
		if _, err := fmt.Fscan(reader, &a, &b); err != nil {
			return fmt.Errorf("reading answer for power %d: %v", power, err)
		}
		if a != exp.i || b != exp.j {
			return fmt.Errorf("power %d: expected (%d,%d) got (%d,%d)",
				power, exp.i, exp.j, a, b)
		}
	}

	inPipe.Close()
	return cmd.Wait()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for _, n := range []int{2, 3, 4, 5} {
		if err := runTest(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "n=%d failed: %v\n", n, err)
			os.Exit(1)
		}
		fmt.Printf("n=%d: OK\n", n)
	}
	fmt.Println("All tests passed")
}
