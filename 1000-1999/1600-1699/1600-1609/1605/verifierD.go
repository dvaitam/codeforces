package main

import (
    "bytes"
    "fmt"
    "math/bits"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type testCaseD struct {
	n     int
	edges [][2]int
}

func generateTestsD() []testCaseD {
	r := rand.New(rand.NewSource(1))
	tests := []testCaseD{}
	for len(tests) < 120 {
		n := r.Intn(10) + 1
		edges := make([][2]int, 0, n-1)
		for i := 2; i <= n; i++ {
			u := r.Intn(i-1) + 1
			edges = append(edges, [2]int{u, i})
		}
		tests = append(tests, testCaseD{n, edges})
	}
	return tests
}

func solveD(n int, edges [][2]int) []int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	maxBit := bits.Len(uint(n))
	vBits := make([][]int, maxBit)
	for i := 1; i <= n; i++ {
		b := bits.Len(uint(i)) - 1
		vBits[b] = append(vBits[b], i)
	}
	depth := make([]int, n+1)
	parent := make([]int, n+1)
	parGroups := [2][]int{}
	queue := make([]int, n)
	qi, qj := 0, 0
	queue[qj] = 1
	qj++
	parent[1] = 0
	depth[1] = 0
	for qi < qj {
		x := queue[qi]
		qi++
		parGroups[depth[x]&1] = append(parGroups[depth[x]&1], x)
		for _, y := range adj[x] {
			if y == parent[x] {
				continue
			}
			parent[y] = x
			depth[y] = depth[x] + 1
			queue[qj] = y
			qj++
		}
	}
	if len(parGroups[0]) > len(parGroups[1]) {
		parGroups[0], parGroups[1] = parGroups[1], parGroups[0]
	}
	cnt0 := len(parGroups[0])
	val := make([]int, n+1)
	cnts := [2]int{}
	for b := 0; b < maxBit; b++ {
		var tmp int
		if (cnt0>>b)&1 == 0 {
			tmp = 1
		} else {
			tmp = 0
		}
		for _, x := range vBits[b] {
			idx := cnts[tmp]
			val[parGroups[tmp][idx]] = x
			cnts[tmp]++
		}
	}
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		res[i-1] = val[i]
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parsePermutation(out string, n int) ([]int, error) {
    fields := strings.Fields(out)
    if len(fields) < n {
        return nil, fmt.Errorf("expected at least %d integers, got %d", n, len(fields))
    }
    p := make([]int, n)
    used := make([]bool, n+1)
    for i := 0; i < n; i++ {
        v, err := strconv.Atoi(fields[i])
        if err != nil {
            return nil, fmt.Errorf("invalid integer at position %d: %v", i+1, err)
        }
        if v < 1 || v > n {
            return nil, fmt.Errorf("value out of range: %d", v)
        }
        if used[v] {
            return nil, fmt.Errorf("duplicate value: %d", v)
        }
        used[v] = true
        p[i] = v
    }
    return p, nil
}

func winningCount(n int, edges [][2]int, perm []int) int {
    // Build adjacency
    adj := make([][]int, n+1)
    for _, e := range edges {
        u, v := e[0], e[1]
        adj[u] = append(adj[u], v)
        adj[v] = append(adj[v], u)
    }
    // Allowed move between u and v if p[u]^p[v] <= min(p[u], p[v])
    allowed := func(u, v int) bool {
        a := perm[u-1]
        b := perm[v-1]
        if a < b {
            return (a^b) <= a
        }
        return (a^b) <= b
    }
    // Memoization for state (u,parent)
    memo := make([][]int8, n+1)
    for i := range memo {
        memo[i] = make([]int8, n+1)
    }
    var dfs func(u, p int) bool
    dfs = func(u, p int) bool {
        if memo[u][p] != 0 {
            return memo[u][p] > 0
        }
        win := false
        for _, v := range adj[u] {
            if v == p {
                continue
            }
            if !allowed(u, v) {
                continue
            }
            if !dfs(v, u) {
                win = true
                break
            }
        }
        if win {
            memo[u][p] = 1
        } else {
            memo[u][p] = -1
        }
        return win
    }
    cnt := 0
    for s := 1; s <= n; s++ {
        if dfs(s, 0) {
            cnt++
        }
    }
    return cnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsD()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input := sb.String()
        // Compute expected winning count using reference construction,
        // but accept any permutation achieving the same count.
        expPerm := solveD(tc.n, tc.edges)
        expCount := winningCount(tc.n, tc.edges, expPerm)
        gotOut, err := runCandidate(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        perm, err := parsePermutation(gotOut, tc.n)
        if err != nil {
            fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\noutput:%s\n", i+1, err, gotOut)
            os.Exit(1)
        }
        gotCount := winningCount(tc.n, tc.edges, perm)
        if gotCount != expCount {
            fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected winning count:%d\ngot winning count:%d\noutput:%s\n", i+1, input, expCount, gotCount, gotOut)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}
