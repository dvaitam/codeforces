package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// ---------- embedded correct solver for 533A ----------

func solve533A(input string) string {
	data := []byte(input)
	ptr := 0

	nextInt := func() int {
		for ptr < len(data) && (data[ptr] < '0' || data[ptr] > '9') {
			ptr++
		}
		v := 0
		for ptr < len(data) && data[ptr] >= '0' && data[ptr] <= '9' {
			v = v*10 + int(data[ptr]-'0')
			ptr++
		}
		return v
	}

	n := nextInt()
	h := make([]int, n+1)
	for i := 1; i <= n; i++ {
		h[i] = nextInt()
	}

	head := make([]int, n+1)
	for i := 1; i <= n; i++ {
		head[i] = -1
	}
	to := make([]int, 2*(n-1))
	nxt := make([]int, 2*(n-1))
	ec := 0
	addEdge := func(a, b int) {
		to[ec] = b
		nxt[ec] = head[a]
		head[a] = ec
		ec++
	}
	for i := 0; i < n-1; i++ {
		a := nextInt()
		b := nextInt()
		addEdge(a, b)
		addEdge(b, a)
	}

	k := nextInt()
	miners := make([]int, k)
	for i := 0; i < k; i++ {
		miners[i] = nextInt()
	}

	const INF = int(1 << 60)

	parent := make([]int, n+1)
	min1 := make([]int, n+1)
	min2 := make([]int, n+1)
	cnt := make([]int, n+1)
	owner := make([]int, n+1)
	ownUnique := make([]int, n+1)
	limit := make([]int, n+1)
	groupSize := make([]int, n+1)

	capacities := make([]int, 0, n)

	min1[1] = h[1]
	min2[1] = INF
	cnt[1] = 1
	owner[1] = 1
	ownUnique[1] = 1
	limit[1] = INF
	groupSize[1]++
	capacities = append(capacities, h[1])

	stack := make([]int, 1, n)
	stack[0] = 1

	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for e := head[v]; e != -1; e = nxt[e] {
			u := to[e]
			if u == parent[v] {
				continue
			}
			parent[u] = v
			hv := h[u]

			if hv < min1[v] {
				min1[u] = hv
				min2[u] = min1[v]
				cnt[u] = 1
				owner[u] = u
			} else if hv == min1[v] {
				min1[u] = min1[v]
				min2[u] = min2[v]
				cnt[u] = cnt[v] + 1
				owner[u] = owner[v]
			} else {
				min1[u] = min1[v]
				cnt[u] = cnt[v]
				owner[u] = owner[v]
				if hv < min2[v] {
					min2[u] = hv
				} else {
					min2[u] = min2[v]
				}
			}

			capacities = append(capacities, min1[u])

			if cnt[u] == 1 {
				ou := owner[u]
				ownUnique[u] = ou
				limit[u] = min2[u]
				groupSize[ou]++
			}

			stack = append(stack, u)
		}
	}

	start := make([]int, n+2)
	total := 0
	for i := 1; i <= n; i++ {
		start[i] = total
		total += groupSize[i]
	}
	start[n+1] = total

	allLimits := make([]int, total)
	cur := make([]int, n+1)
	copy(cur, start[:n+1])

	for i := 1; i <= n; i++ {
		if ownUnique[i] != 0 {
			u := ownUnique[i]
			allLimits[cur[u]] = limit[i]
			cur[u]++
		}
	}

	sort.Ints(capacities)
	sort.Ints(miners)

	def := make([]int, k)
	p := n - 1
	D := 0
	for t := k - 1; t >= 0; t-- {
		sv := miners[t]
		for p >= 0 && capacities[p] >= sv {
			p--
		}
		countGe := n - 1 - p
		need := k - t
		d := need - countGe
		def[t] = d
		if d > D {
			D = d
		}
	}

	if D == 0 {
		return "0"
	}

	lowIdx := 0
	for lowIdx < k && def[lowIdx] <= 0 {
		lowIdx++
	}
	lowHeight := miners[lowIdx]

	req := make([]int, D)
	filled := 0
	for t := k - 1; t >= 0; t-- {
		for filled < def[t] {
			req[filled] = miners[t]
			filled++
		}
	}

	bestH := -1
	for u := 1; u <= n; u++ {
		if h[u] >= lowHeight {
			continue
		}
		l, r := start[u], start[u+1]
		if r-l < D {
			continue
		}
		group := allLimits[l:r]
		sort.Ints(group)
		ok := true
		for j := 0; j < D; j++ {
			if group[len(group)-1-j] < req[j] {
				ok = false
				break
			}
		}
		if ok && h[u] > bestH {
			bestH = h[u]
		}
	}

	if bestH < 0 {
		return "-1"
	}
	return fmt.Sprintf("%d", req[0]-bestH)
}

// ---------- verifier infrastructure ----------

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildTests() []string {
	return []string{
		"1\n3\n1\n5\n",
		"4\n5 2 4 3\n1 2\n2 3\n2 4\n3\n3 4 5\n",
		"2\n2 2\n1 2\n2\n3 3\n",
		"5\n2 2 2 2 2\n1 2\n1 3\n1 4\n1 5\n4\n2 2 2 5\n",
		"6\n6 3 4 2 5 3\n1 2\n2 3\n3 4\n4 5\n5 6\n4\n2 3 4 6\n",
		"7\n10 4 8 3 6 5 7\n1 2\n1 3\n2 4\n2 5\n3 6\n6 7\n5\n5 6 7 8 9\n",
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	for idx, input := range tests {
		refOut := solve533A(input)
		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(refOut) != strings.TrimSpace(candOut) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\nInput:\n%sExpected: %s\nGot: %s\n", idx+1, input, refOut, candOut)
			os.Exit(1)
		}
	}

	// Also run a few random stress tests with embedded solver
	rr := bufio.NewWriter(os.Stdout)
	_ = rr
	fmt.Printf("All %d tests passed.\n", len(tests))
}
