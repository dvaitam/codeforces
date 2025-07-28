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
		expVals := solveD(tc.n, tc.edges)
		expStr := strings.TrimSpace(strings.Join(func(arr []int) []string {
			s := make([]string, len(arr))
			for i, v := range arr {
				s[i] = fmt.Sprintf("%d", v)
			}
			return s
		}(expVals), " "))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expStr {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, expStr, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
