package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const alphabet = 14

type node struct {
	next [alphabet]int
	fail int
	val  int64
}

func newNode() node {
	n := node{fail: 0, val: 0}
	for i := 0; i < alphabet; i++ {
		n.next[i] = -1
	}
	return n
}

func insert(trie *[]node, s string, cost int64) {
	cur := 0
	for i := 0; i < len(s); i++ {
		idx := int(s[i] - 'a')
		if (*trie)[cur].next[idx] == -1 {
			*trie = append(*trie, newNode())
			(*trie)[cur].next[idx] = len(*trie) - 1
		}
		cur = (*trie)[cur].next[idx]
	}
	(*trie)[cur].val += cost
}

func build(trie []node) {
	queue := make([]int, 0)
	for c := 0; c < alphabet; c++ {
		if trie[0].next[c] != -1 {
			child := trie[0].next[c]
			trie[child].fail = 0
			queue = append(queue, child)
		} else {
			trie[0].next[c] = 0
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		f := trie[v].fail
		trie[v].val += trie[f].val
		for c := 0; c < alphabet; c++ {
			if trie[v].next[c] != -1 {
				child := trie[v].next[c]
				trie[child].fail = trie[f].next[c]
				queue = append(queue, child)
			} else {
				trie[v].next[c] = trie[f].next[c]
			}
		}
	}
}

func segTransition(trie []node, seg string) ([]int, []int64) {
	n := len(trie)
	nxt := make([]int, n)
	val := make([]int64, n)
	for s := 0; s < n; s++ {
		cur := s
		sum := int64(0)
		for i := 0; i < len(seg); i++ {
			ch := int(seg[i] - 'a')
			cur = trie[cur].next[ch]
			sum += trie[cur].val
		}
		nxt[s] = cur
		val[s] = sum
	}
	return nxt, val
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solve(k int, ts []string, cs []int64, S string) int64 {
	trie := []node{newNode()}
	for i := 0; i < k; i++ {
		insert(&trie, ts[i], cs[i])
	}
	build(trie)
	segments := make([]string, 0)
	sb := make([]byte, 0, len(S))
	qm := 0
	for i := 0; i < len(S); i++ {
		if S[i] == '?' {
			segments = append(segments, string(sb))
			sb = sb[:0]
			qm++
		} else {
			sb = append(sb, S[i])
		}
	}
	segments = append(segments, string(sb))
	m := qm
	nStates := len(trie)
	segNext := make([][]int, len(segments))
	segVal := make([][]int64, len(segments))
	for i := 0; i < len(segments); i++ {
		segNext[i], segVal[i] = segTransition(trie, segments[i])
	}
	startState := segNext[0][0]
	startVal := segVal[0][0]
	if m == 0 {
		return startVal
	}
	masksByCnt := make([][]int, m+1)
	for mask := 0; mask < (1 << alphabet); mask++ {
		cnt := bits.OnesCount(uint(mask))
		if cnt <= m {
			masksByCnt[cnt] = append(masksByCnt[cnt], mask)
		}
	}
	maskIdx := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		maskIdx[i] = make([]int, 1<<alphabet)
		for j := range maskIdx[i] {
			maskIdx[i][j] = -1
		}
		for idx, mask := range masksByCnt[i] {
			maskIdx[i][mask] = idx
		}
	}
	negInf := int64(-1 << 60)
	dpPrev := make([][]int64, len(masksByCnt[0]))
	for idx := range dpPrev {
		dpPrev[idx] = make([]int64, nStates)
		for j := range dpPrev[idx] {
			dpPrev[idx][j] = negInf
		}
	}
	dpPrev[maskIdx[0][0]][startState] = startVal
	for step := 1; step <= m; step++ {
		dpNext := make([][]int64, len(masksByCnt[step]))
		for idx := range dpNext {
			dpNext[idx] = make([]int64, nStates)
			for j := range dpNext[idx] {
				dpNext[idx][j] = negInf
			}
		}
		for pIdx, mask := range masksByCnt[step-1] {
			arr := dpPrev[pIdx]
			for s := 0; s < nStates; s++ {
				val := arr[s]
				if val == negInf {
					continue
				}
				for l := 0; l < alphabet; l++ {
					if mask&(1<<l) != 0 {
						continue
					}
					newMask := mask | (1 << l)
					nIdx := maskIdx[step][newMask]
					if nIdx == -1 {
						continue
					}
					s1 := trie[s].next[l]
					v1 := val + trie[s1].val
					s2 := segNext[step][s1]
					v2 := v1 + segVal[step][s1]
					if v2 > dpNext[nIdx][s2] {
						dpNext[nIdx][s2] = v2
					}
				}
			}
		}
		dpPrev = dpNext
	}
	ans := negInf
	for _, arr := range dpPrev {
		for s := 0; s < nStates; s++ {
			ans = max64(ans, arr[s])
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(7)
	t := 100
	for idx := 0; idx < t; idx++ {
		k := rand.Intn(3) + 1
		ts := make([]string, k)
		cs := make([]int64, k)
		for i := 0; i < k; i++ {
			length := rand.Intn(3) + 1
			b := make([]byte, length)
			for j := range b {
				b[j] = byte('a' + rand.Intn(alphabet))
			}
			ts[i] = string(b)
			cs[i] = int64(rand.Intn(21) - 10)
		}
		lengthS := rand.Intn(8) + 1
		qm := rand.Intn(min(3, lengthS))
		qmPositions := rand.Perm(lengthS)[:qm]
		posMap := make(map[int]bool)
		for _, p := range qmPositions {
			posMap[p] = true
		}
		b := make([]byte, lengthS)
		for j := 0; j < lengthS; j++ {
			if posMap[j] {
				b[j] = '?'
			} else {
				b[j] = byte('a' + rand.Intn(alphabet))
			}
		}
		S := string(b)
		want := fmt.Sprintf("%d\n", solve(k, ts, cs, S))

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", k))
		for i := 0; i < k; i++ {
			input.WriteString(fmt.Sprintf("%s %d\n", ts[i], cs[i]))
		}
		input.WriteString(fmt.Sprintf("%s\n", S))

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Runtime error on test %d: %v\n%s", idx+1, err, out.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != strings.TrimSpace(want) {
			fmt.Printf("Wrong answer on test %d\nExpected: %s\nGot: %s\n", idx+1, strings.TrimSpace(want), got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
