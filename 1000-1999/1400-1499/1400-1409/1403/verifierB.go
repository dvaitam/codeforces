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
	"time"
)

// ---------- Embedded reference solver for 1403B ----------

func refSolve(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 0, 10*1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	read := func() int {
		scanner.Scan()
		num, _ := strconv.Atoi(scanner.Text())
		return num
	}

	N := read()
	Q := read()

	adj := make([][]int, N+1)
	origDegree := make([]int, N+1)
	for i := 0; i < N-1; i++ {
		u, v := read(), read()
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		origDegree[u]++
		origDegree[v]++
	}

	LTotal := 0
	for i := 1; i <= N; i++ {
		if origDegree[i] == 1 {
			LTotal++
		}
	}

	tin := make([]int, N+1)
	tout := make([]int, N+1)
	up := make([][20]int, N+1)
	depth := make([]int, N+1)
	lOrig := make([]int, N+1)
	timer := 0

	var dfs1 func(u, p int)
	dfs1 = func(u, p int) {
		timer++
		tin[u] = timer
		up[u][0] = p
		for i := 1; i < 20; i++ {
			up[u][i] = up[up[u][i-1]][i-1]
		}
		if origDegree[u] == 1 {
			lOrig[u] = 1
		}
		for _, v := range adj[u] {
			if v != p {
				depth[v] = depth[u] + 1
				dfs1(v, u)
				lOrig[u] += lOrig[v]
			}
		}
		tout[u] = timer
	}

	depth[1] = 0
	dfs1(1, 1)

	DistEven := make([]int, N+1)
	DistOdd := make([]int, N+1)
	TotalOriginalEvenCount := 0

	var dfs2 func(u, p int)
	dfs2 = func(u, p int) {
		for _, v := range adj[u] {
			if v != p {
				if lOrig[v]%2 == 0 {
					DistEven[v] = DistEven[u] + 1
					DistOdd[v] = DistOdd[u]
					TotalOriginalEvenCount++
				} else {
					DistEven[v] = DistEven[u]
					DistOdd[v] = DistOdd[u] + 1
				}
				dfs2(v, u)
			}
		}
	}
	dfs2(1, 1)

	getLca := func(u, v int) int {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		diff := depth[u] - depth[v]
		for i := 0; i < 20; i++ {
			if (diff & (1 << i)) != 0 {
				u = up[u][i]
			}
		}
		if u == v {
			return u
		}
		for i := 19; i >= 0; i-- {
			if up[u][i] != up[v][i] {
				u = up[u][i]
				v = up[v][i]
			}
		}
		return up[u][0]
	}

	count := make([]int, N+1)
	C := make([]int, 0, N)
	virtAdj := make([][]int, N+1)
	isF := make([]bool, N+1)

	var out strings.Builder

	for q := 0; q < Q; q++ {
		Di := read()
		C = C[:0]
		for j := 0; j < Di; j++ {
			u := read()
			if count[u] == 0 {
				C = append(C, u)
			}
			count[u]++
		}

		newLeaves := LTotal + Di
		for _, u := range C {
			if origDegree[u] == 1 {
				newLeaves--
			}
		}

		if newLeaves%2 != 0 {
			out.WriteString("-1\n")
			for _, u := range C {
				count[u] = 0
			}
			continue
		}

		SInitial := make([]int, 0, len(C))
		for _, u := range C {
			cu := count[u]
			fu := 0
			if origDegree[u] == 1 {
				if cu%2 == 0 {
					fu = 1
				}
			} else {
				if cu%2 != 0 {
					fu = 1
				}
			}
			if fu == 1 {
				SInitial = append(SInitial, u)
			}
		}

		S := make([]int, len(SInitial))
		copy(S, SInitial)
		hasRoot := false
		for _, u := range S {
			if u == 1 {
				hasRoot = true
				break
			}
		}
		if !hasRoot {
			S = append(S, 1)
		}

		for _, u := range SInitial {
			isF[u] = true
		}

		sort.Slice(S, func(i, j int) bool { return tin[S[i]] < tin[S[j]] })
		m := len(S)
		for i := 0; i < m-1; i++ {
			lca := getLca(S[i], S[i+1])
			S = append(S, lca)
		}

		sort.Slice(S, func(i, j int) bool { return tin[S[i]] < tin[S[j]] })
		uniqueS := []int{S[0]}
		for i := 1; i < len(S); i++ {
			if S[i] != S[i-1] {
				uniqueS = append(uniqueS, S[i])
			}
		}
		S = uniqueS

		for _, u := range S {
			virtAdj[u] = virtAdj[u][:0]
		}

		stack := []int{S[0]}
		for i := 1; i < len(S); i++ {
			v := S[i]
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if tin[top] <= tin[v] && tout[v] <= tout[top] {
					break
				}
				stack = stack[:len(stack)-1]
			}
			p := stack[len(stack)-1]
			virtAdj[p] = append(virtAdj[p], v)
			stack = append(stack, v)
		}

		flipCount := 0
		var dp func(u int) int
		dp = func(u int) int {
			w := 0
			if isF[u] {
				w = 1
			}
			for _, v := range virtAdj[u] {
				childW := dp(v)
				w = (w + childW) % 2
				if childW == 1 {
					eEven := DistEven[v] - DistEven[u]
					eOdd := DistOdd[v] - DistOdd[u]
					flipCount += eOdd - eEven
				}
			}
			return w
		}

		dp(S[0])

		ans := (N - 1 + Di) + TotalOriginalEvenCount + flipCount
		out.WriteString(strconv.Itoa(ans))
		out.WriteByte('\n')

		for _, u := range SInitial {
			isF[u] = false
		}
		for _, u := range C {
			count[u] = 0
		}
	}

	return strings.TrimSpace(out.String())
}

// ---------- test harness ----------

func generateCase(rng *rand.Rand) string {
	N := rng.Intn(5) + 2
	Q := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", N, Q))
	for i := 2; i <= N; i++ {
		v := rng.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", i, v))
	}
	for i := 0; i < Q; i++ {
		Di := rng.Intn(N) + 1
		sb.WriteString(fmt.Sprintf("%d", Di))
		for j := 0; j < Di; j++ {
			sb.WriteString(fmt.Sprintf(" %d", rng.Intn(N)+1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp := refSolve(input)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed:\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
