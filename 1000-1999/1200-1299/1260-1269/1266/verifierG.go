package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod = 998244353

type test struct {
	input    string
	expected string
}

type samState struct {
	next   map[int]int
	link   int
	length int
}

type sam struct {
	states []samState
	last   int
}

func newSAM(cap int) *sam {
	s := &sam{states: make([]samState, 1, cap)}
	s.states[0] = samState{next: make(map[int]int), link: -1}
	return s
}

func (s *sam) extend(c int) {
	cur := len(s.states)
	s.states = append(s.states, samState{next: make(map[int]int), length: s.states[s.last].length + 1})
	p := s.last
	for p != -1 && s.states[p].next[c] == 0 {
		s.states[p].next[c] = cur
		p = s.states[p].link
	}
	if p == -1 {
		s.states[cur].link = 0
	} else {
		q := s.states[p].next[c]
		if s.states[p].length+1 == s.states[q].length {
			s.states[cur].link = q
		} else {
			clone := len(s.states)
			cloneState := samState{next: make(map[int]int), length: s.states[p].length + 1, link: s.states[q].link}
			for k, v := range s.states[q].next {
				cloneState.next[k] = v
			}
			s.states = append(s.states, cloneState)
			for p != -1 && s.states[p].next[c] == q {
				s.states[p].next[c] = clone
				p = s.states[p].link
			}
			s.states[q].link = clone
			s.states[cur].link = clone
		}
	}
	s.last = cur
}

func (s *sam) distinct() int64 {
	var res int64
	for i := 1; i < len(s.states); i++ {
		link := s.states[i].link
		res += int64(s.states[i].length - s.states[link].length)
	}
	return res
}

func nextPermutation(a []int) bool {
	n := len(a)
	i := n - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for j > i && a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func solveBrute(n int) int64 {
	if n <= 0 || n > 8 {
		return -1
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	total := 1
	for i := 2; i <= n; i++ {
		total *= i
	}
	s := newSAM(n*total*2 + 5)
	for {
		for _, v := range perm {
			s.extend(v)
		}
		if !nextPermutation(perm) {
			break
		}
	}
	return s.distinct() % mod
}

func formatInput(n int) string {
	return fmt.Sprintf("%d\n", n)
}

func fixedTests() []test {
	tests := make([]test, 0, 7)
	for n := 1; n <= 7; n++ {
		ans := solveBrute(n)
		tests = append(tests, test{
			input:    formatInput(n),
			expected: fmt.Sprintf("%d", ans),
		})
	}
	return tests
}

func randomTests(rng *rand.Rand, count int) []test {
	tests := make([]test, 0, count)
	for len(tests) < count {
		n := rng.Intn(7) + 1
		ans := solveBrute(n)
		tests = append(tests, test{
			input:    formatInput(n),
			expected: fmt.Sprintf("%d", ans),
		})
	}
	return tests
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(1266))
	tests := fixedTests()
	tests = append(tests, randomTests(rng, 40)...)
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\nInput:%s\n", i+1, err, t.input)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
