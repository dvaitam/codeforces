package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type instance struct {
	n    int
	k    int
	strs []string
}

type node struct {
	next [2]int
	fail int
	out  uint64
}

type solver struct {
	inst      instance
	trie      []node
	nodeCount int
	dpCurr    []bool
	dpNext    []bool
	avoidMemo map[uint64]bool
	order     []int
}

func newSolver(inst instance) *solver {
	s := &solver{
		inst: inst,
	}
	s.buildAutomaton()
	s.nodeCount = len(s.trie)
	dpSize := (inst.k + 1) * s.nodeCount
	s.dpCurr = make([]bool, dpSize)
	s.dpNext = make([]bool, dpSize)
	s.avoidMemo = make(map[uint64]bool)
	s.order = make([]int, inst.n)
	for i := 0; i < inst.n; i++ {
		s.order[i] = i
	}
	sort.Slice(s.order, func(i, j int) bool {
		li := len(inst.strs[s.order[i]])
		lj := len(inst.strs[s.order[j]])
		if li != lj {
			return li > lj
		}
		return inst.strs[s.order[i]] < inst.strs[s.order[j]]
	})
	return s
}

func (s *solver) buildAutomaton() {
	s.trie = []node{{next: [2]int{-1, -1}}}
	for idx, str := range s.inst.strs {
		cur := 0
		for _, ch := range str {
			id := 0
			if ch == ')' {
				id = 1
			}
			if s.trie[cur].next[id] == -1 {
				s.trie[cur].next[id] = len(s.trie)
				s.trie = append(s.trie, node{next: [2]int{-1, -1}})
			}
			cur = s.trie[cur].next[id]
		}
		s.trie[cur].out |= 1 << uint(idx)
	}
	queue := make([]int, 0)
	for c := 0; c < 2; c++ {
		nxt := s.trie[0].next[c]
		if nxt != -1 {
			s.trie[nxt].fail = 0
			queue = append(queue, nxt)
		} else {
			s.trie[0].next[c] = 0
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		fail := s.trie[v].fail
		s.trie[v].out |= s.trie[fail].out
		for c := 0; c < 2; c++ {
			nxt := s.trie[v].next[c]
			if nxt != -1 {
				s.trie[nxt].fail = s.trie[fail].next[c]
				queue = append(queue, nxt)
			} else {
				s.trie[v].next[c] = s.trie[fail].next[c]
			}
		}
	}
}

func (s *solver) canAvoid(mask uint64) bool {
	if val, ok := s.avoidMemo[mask]; ok {
		return val
	}
	blocked := make([]bool, s.nodeCount)
	for i := 0; i < s.nodeCount; i++ {
		if s.trie[i].out&mask != 0 {
			blocked[i] = true
		}
	}
	if blocked[0] {
		s.avoidMemo[mask] = false
		return false
	}
	for i := range s.dpCurr {
		s.dpCurr[i] = false
	}
	for i := range s.dpNext {
		s.dpNext[i] = false
	}
	s.dpCurr[0] = true
	k := s.inst.k
	for pos := 0; pos < k; pos++ {
		for i := range s.dpNext {
			s.dpNext[i] = false
		}
		for bal := 0; bal <= k; bal++ {
			base := bal * s.nodeCount
			for state := 0; state < s.nodeCount; state++ {
				if !s.dpCurr[base+state] {
					continue
				}
				if blocked[state] {
					continue
				}
				// add '('
				if bal+1 <= k {
					ns := s.trie[state].next[0]
					if !blocked[ns] {
						s.dpNext[(bal+1)*s.nodeCount+ns] = true
					}
				}
				// add ')'
				if bal > 0 {
					ns := s.trie[state].next[1]
					if !blocked[ns] {
						s.dpNext[(bal-1)*s.nodeCount+ns] = true
					}
				}
			}
		}
		s.dpCurr, s.dpNext = s.dpNext, s.dpCurr
	}
	ok := false
	for state := 0; state < s.nodeCount; state++ {
		if s.dpCurr[state] && !blocked[state] {
			ok = true
			break
		}
	}
	s.avoidMemo[mask] = ok
	return ok
}

func (s *solver) minimalGroups() (int, bool) {
	for i := 0; i < s.inst.n; i++ {
		mask := uint64(1) << uint(i)
		if !s.canAvoid(mask) {
			return 0, false
		}
	}
	for m := 1; m <= s.inst.n; m++ {
		if s.tryAssign(m) {
			return m, true
		}
	}
	return 0, false
}

func (s *solver) tryAssign(m int) bool {
	groupMasks := make([]uint64, m)
	assign := make([]int, s.inst.n)
	for i := range assign {
		assign[i] = -1
	}
	var dfs func(pos, used int) bool
	dfs = func(pos, used int) bool {
		if pos == len(s.order) {
			return true
		}
		idx := s.order[pos]
		bit := uint64(1) << uint(idx)
		for g := 0; g < used; g++ {
			oldMask := groupMasks[g]
			newMask := oldMask | bit
			if newMask == oldMask {
				continue
			}
			if !s.canAvoid(newMask) {
				continue
			}
			groupMasks[g] = newMask
			assign[idx] = g
			if dfs(pos+1, used) {
				return true
			}
			groupMasks[g] = oldMask
			assign[idx] = -1
		}
		if used < m {
			if s.canAvoid(bit) {
				groupMasks[used] = bit
				assign[idx] = used
				if dfs(pos+1, used+1) {
					return true
				}
				groupMasks[used] = 0
				assign[idx] = -1
			}
		}
		return false
	}
	return dfs(0, 0)
}

func isRegular(seq string) bool {
	bal := 0
	for _, ch := range seq {
		if ch == '(' {
			bal++
		} else if ch == ')' {
			bal--
		} else {
			return false
		}
		if bal < 0 {
			return false
		}
	}
	return bal == 0
}

func validateOutput(inst instance, out string, possible bool, optimal int) error {
	reader := bufio.NewReader(strings.NewReader(out))
	var first string
	if _, err := fmt.Fscan(reader, &first); err != nil {
		return fmt.Errorf("output incomplete")
	}
	if first == "-1" {
		if possible {
			return fmt.Errorf("solution exists but -1 printed")
		}
		if hasExtraNonSpace(reader) {
			return fmt.Errorf("extra tokens after -1")
		}
		return nil
	}
	mVal, err := strconv.Atoi(first)
	if err != nil {
		return fmt.Errorf("expected integer or -1, got %q", first)
	}
	if !possible {
		return fmt.Errorf("printed %d groups but answer is -1", mVal)
	}
	if mVal != optimal {
		return fmt.Errorf("number of groups %d differs from optimal %d", mVal, optimal)
	}
	if mVal <= 0 {
		return fmt.Errorf("invalid group count %d", mVal)
	}
	used := make([]bool, inst.n)
	for g := 0; g < mVal; g++ {
		var seq string
		if _, err := fmt.Fscan(reader, &seq); err != nil {
			return fmt.Errorf("missing sequence for group %d", g+1)
		}
		if len(seq) != inst.k {
			return fmt.Errorf("group %d sequence length %d != k", g+1, len(seq))
		}
		if !isRegular(seq) {
			return fmt.Errorf("group %d sequence is not regular", g+1)
		}
		var size int
		if _, err := fmt.Fscan(reader, &size); err != nil {
			return fmt.Errorf("missing size for group %d", g+1)
		}
		if size < 0 {
			return fmt.Errorf("negative group size %d", size)
		}
		if size > inst.n {
			return fmt.Errorf("group %d size too large", g+1)
		}
		members := make([]int, size)
		for i := 0; i < size; i++ {
			if _, err := fmt.Fscan(reader, &members[i]); err != nil {
				return fmt.Errorf("missing index %d in group %d", i+1, g+1)
			}
			if members[i] < 1 || members[i] > inst.n {
				return fmt.Errorf("index %d out of range", members[i])
			}
			if used[members[i]-1] {
				return fmt.Errorf("index %d appears multiple times", members[i])
			}
			used[members[i]-1] = true
			if strings.Contains(seq, inst.strs[members[i]-1]) {
				return fmt.Errorf("string %d is substring of group %d sequence", members[i], g+1)
			}
		}
	}
	for i, ok := range used {
		if !ok {
			return fmt.Errorf("string %d not assigned to any group", i+1)
		}
	}
	if hasExtraNonSpace(reader) {
		return fmt.Errorf("extra data after processing groups")
	}
	return nil
}

func hasExtraNonSpace(r *bufio.Reader) bool {
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return false
		}
		if !unicode.IsSpace(rune(ch)) {
			return true
		}
	}
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
	return out.String(), err
}

func formatInput(inst instance) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", inst.n, inst.k))
	for _, s := range inst.strs {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func fixedTests() []instance {
	return []instance{
		{n: 1, k: 2, strs: []string{"()"}},
		{n: 1, k: 4, strs: []string{"(("}},
		{n: 2, k: 4, strs: []string{"()", "()"}},
		{n: 3, k: 6, strs: []string{"(()", ")()", ")))"}},
	}
}

func randomTests() []instance {
	rng := rand.New(rand.NewSource(2144))
	var tests []instance
	for t := 0; t < 40; t++ {
		n := rng.Intn(6) + 1
		kOpts := []int{2, 4, 6, 8}
		k := kOpts[rng.Intn(len(kOpts))]
		strs := make([]string, n)
		for i := 0; i < n; i++ {
			length := rng.Intn(k-1) + 2
			var sb strings.Builder
			for j := 0; j < length; j++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('(')
				} else {
					sb.WriteByte(')')
				}
			}
			strs[i] = sb.String()
		}
		tests = append(tests, instance{n: n, k: k, strs: strs})
	}
	return tests
}

func generateTests() []instance {
	tests := fixedTests()
	tests = append(tests, randomTests()...)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, inst := range tests {
		input := formatInput(inst)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\nInput:%s\n", i+1, err, input)
			os.Exit(1)
		}
		s := newSolver(inst)
		minGroups, possible := s.minimalGroups()
		if err := validateOutput(inst, output, possible, minGroups); err != nil {
			expected := "-1"
			if possible {
				expected = fmt.Sprintf("minimal groups = %d", minGroups)
			}
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\nReason:%v\n", i+1, input, expected, output, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
