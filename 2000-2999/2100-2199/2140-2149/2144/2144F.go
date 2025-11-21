package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type node struct {
	next [2]int
	fail int
	out  uint64
}

var (
	n, k           int
	strs           []string
	trie           []node
	nodeCount      int
	maxBalance     int
	dpCurr         []bool
	dpNext         []bool
	avoidMemo      map[uint64]bool
	seqMemo        map[uint64]string
	order          []int
	solutionMasks  []uint64
	solutionAssign []int
	solutionCount  int
)

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	strs = make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &strs[i])
	}

	buildAutomaton()
	nodeCount = len(trie)
	maxBalance = k
	dpSize := (maxBalance + 1) * nodeCount
	dpCurr = make([]bool, dpSize)
	dpNext = make([]bool, dpSize)
	avoidMemo = make(map[uint64]bool)
	seqMemo = make(map[uint64]string)

	for i := 0; i < n; i++ {
		mask := uint64(1) << uint(i)
		if !canAvoid(mask) {
			fmt.Println(-1)
			return
		}
	}

	order = make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		li := len(strs[order[i]])
		lj := len(strs[order[j]])
		if li != lj {
			return li > lj
		}
		return strs[order[i]] < strs[order[j]]
	})

	found := false
	for m := 1; m <= n; m++ {
		if attempt(m) {
			found = true
			break
		}
	}
	if !found {
		fmt.Println(-1)
		return
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, solutionCount)
	for g := 0; g < solutionCount; g++ {
		mask := solutionMasks[g]
		seq, ok := constructSequence(mask)
		if !ok {
			// should not happen because mask already verified
			fmt.Fprintln(out, "()")
		} else {
			fmt.Fprintln(out, seq)
		}
		var members []int
		for i := 0; i < n; i++ {
			if solutionAssign[i] == g {
				members = append(members, i+1)
			}
		}
		fmt.Fprintln(out, len(members))
		for i, idx := range members {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, idx)
		}
		fmt.Fprintln(out)
	}
	out.Flush()
}

func buildAutomaton() {
	trie = []node{{next: [2]int{-1, -1}}}
	for idx, s := range strs {
		cur := 0
		for _, ch := range s {
			var id int
			if ch == '(' {
				id = 0
			} else {
				id = 1
			}
			if trie[cur].next[id] == -1 {
				trie[cur].next[id] = len(trie)
				trie = append(trie, node{next: [2]int{-1, -1}})
			}
			cur = trie[cur].next[id]
		}
		trie[cur].out |= uint64(1) << uint(idx)
	}
	queue := make([]int, 0)
	for c := 0; c < 2; c++ {
		nxt := trie[0].next[c]
		if nxt != -1 {
			trie[nxt].fail = 0
			queue = append(queue, nxt)
		} else {
			trie[0].next[c] = 0
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		fail := trie[v].fail
		trie[v].out |= trie[fail].out
		for c := 0; c < 2; c++ {
			nxt := trie[v].next[c]
			if nxt != -1 {
				trie[nxt].fail = trie[fail].next[c]
				queue = append(queue, nxt)
			} else {
				trie[v].next[c] = trie[fail].next[c]
			}
		}
	}
}

func canAvoid(mask uint64) bool {
	if val, ok := avoidMemo[mask]; ok {
		return val
	}
	blocked := make([]bool, nodeCount)
	for i := 0; i < nodeCount; i++ {
		if trie[i].out&mask != 0 {
			blocked[i] = true
		}
	}
	if blocked[0] {
		avoidMemo[mask] = false
		return false
	}
	for i := range dpCurr {
		dpCurr[i] = false
	}
	for i := range dpNext {
		dpNext[i] = false
	}
	dpCurr[0] = true
	for pos := 0; pos < k; pos++ {
		for i := range dpNext {
			dpNext[i] = false
		}
		for bal := 0; bal <= maxBalance; bal++ {
			base := bal * nodeCount
			for state := 0; state < nodeCount; state++ {
				if !dpCurr[base+state] {
					continue
				}
				if blocked[state] {
					continue
				}
				if bal+1 <= maxBalance {
					ns := trie[state].next[0]
					if !blocked[ns] {
						dpNext[(bal+1)*nodeCount+ns] = true
					}
				}
				if bal > 0 {
					ns := trie[state].next[1]
					if !blocked[ns] {
						dpNext[(bal-1)*nodeCount+ns] = true
					}
				}
			}
		}
		dpCurr, dpNext = dpNext, dpCurr
	}
	ok := false
	for state := 0; state < nodeCount; state++ {
		if dpCurr[state] {
			ok = true
			break
		}
	}
	avoidMemo[mask] = ok
	return ok
}

func attempt(m int) bool {
	groupMasks := make([]uint64, m)
	assign := make([]int, n)
	for i := range assign {
		assign[i] = -1
	}
	var dfs func(pos, used int) bool
	dfs = func(pos, used int) bool {
		if pos == len(order) {
			solutionCount = used
			solutionMasks = append([]uint64(nil), groupMasks[:used]...)
			solutionAssign = append([]int(nil), assign...)
			return true
		}
		idx := order[pos]
		bit := uint64(1) << uint(idx)
		for g := 0; g < used; g++ {
			oldMask := groupMasks[g]
			newMask := oldMask | bit
			if newMask == oldMask {
				continue
			}
			if !canAvoid(newMask) {
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
			if canAvoid(bit) {
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

type choiceInfo struct {
	ch        byte
	nextState int
}

func constructSequence(mask uint64) (string, bool) {
	if seq, ok := seqMemo[mask]; ok {
		return seq, true
	}
	blocked := make([]bool, nodeCount)
	for i := 0; i < nodeCount; i++ {
		if trie[i].out&mask != 0 {
			blocked[i] = true
		}
	}
	if blocked[0] {
		return "", false
	}
	choice := make(map[int]choiceInfo)
	memo := make(map[int]bool)
	var dfs func(pos, bal, state int) bool
	keyFor := func(pos, bal, state int) int {
		return ((pos*(maxBalance+1))+bal)*nodeCount + state
	}
	dfs = func(pos, bal, state int) bool {
		if blocked[state] {
			return false
		}
		if pos == k {
			return bal == 0
		}
		key := keyFor(pos, bal, state)
		if val, ok := memo[key]; ok {
			return val
		}
		if bal+1 <= maxBalance {
			ns := trie[state].next[0]
			if !blocked[ns] && dfs(pos+1, bal+1, ns) {
				choice[key] = choiceInfo{'(', ns}
				memo[key] = true
				return true
			}
		}
		if bal > 0 {
			ns := trie[state].next[1]
			if !blocked[ns] && dfs(pos+1, bal-1, ns) {
				choice[key] = choiceInfo{')', ns}
				memo[key] = true
				return true
			}
		}
		memo[key] = false
		return false
	}
	if !dfs(0, 0, 0) {
		return "", false
	}
	var sb []byte
	bal := 0
	state := 0
	for pos := 0; pos < k; pos++ {
		key := keyFor(pos, bal, state)
		info := choice[key]
		sb = append(sb, info.ch)
		if info.ch == '(' {
			bal++
		} else {
			bal--
		}
		state = info.nextState
	}
	seq := string(sb)
	seqMemo[mask] = seq
	return seq, true
}
