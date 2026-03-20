package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSourceF = `package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var buffer []byte
var pos int

func nextInt() int {
	for pos < len(buffer) && buffer[pos] <= ' ' {
		pos++
	}
	if pos >= len(buffer) {
		return 0
	}
	sign := 1
	if buffer[pos] == '-' {
		sign = -1
		pos++
	}
	res := 0
	for pos < len(buffer) && buffer[pos] > ' ' {
		res = res*10 + int(buffer[pos]-'0')
		pos++
	}
	return res * sign
}

func lit(x int) int {
	if x > 0 {
		return 2 * (x - 1)
	}
	return 2*(-x-1) + 1
}

func neg(u int) int {
	return u ^ 1
}

func solve2SAT(n int, edges [][]int, force []int) (bool, []int) {
	adj := make([][]int, 2*n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[neg(u)] = append(adj[neg(u)], v)
		adj[neg(v)] = append(adj[neg(v)], u)
	}
	for _, u := range force {
		adj[neg(u)] = append(adj[neg(u)], u)
	}

	dfn := make([]int, 2*n)
	for i := range dfn {
		dfn[i] = -1
	}
	low := make([]int, 2*n)
	comp := make([]int, 2*n)
	inStk := make([]bool, 2*n)
	stk := make([]int, 0, 2*n)
	timer := 0
	sccCnt := 0

	var dfs func(int)
	dfs = func(u int) {
		dfn[u] = timer
		low[u] = timer
		timer++
		stk = append(stk, u)
		inStk[u] = true
		for _, v := range adj[u] {
			if dfn[v] == -1 {
				dfs(v)
				if low[v] < low[u] {
					low[u] = low[v]
				}
			} else if inStk[v] {
				if dfn[v] < low[u] {
					low[u] = dfn[v]
				}
			}
		}
		if low[u] == dfn[u] {
			for {
				v := stk[len(stk)-1]
				stk = stk[:len(stk)-1]
				inStk[v] = false
				comp[v] = sccCnt
				if v == u {
					break
				}
			}
			sccCnt++
		}
	}

	for i := 0; i < 2*n; i++ {
		if dfn[i] == -1 {
			dfs(i)
		}
	}

	for i := 0; i < n; i++ {
		if comp[2*i] == comp[2*i+1] {
			return false, nil
		}
	}

	ass := make([]int, n)
	for i := 0; i < n; i++ {
		if comp[2*i] < comp[2*i+1] {
			ass[i] = 1
		} else {
			ass[i] = 0
		}
	}
	return true, ass
}

func buildReach(n int, edges [][]int) ([][]uint64, []int) {
	adj := make([][]int, 2*n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[neg(u)] = append(adj[neg(u)], v)
		adj[neg(v)] = append(adj[neg(v)], u)
	}

	dfn := make([]int, 2*n)
	for i := range dfn {
		dfn[i] = -1
	}
	low := make([]int, 2*n)
	comp := make([]int, 2*n)
	inStk := make([]bool, 2*n)
	stk := make([]int, 0, 2*n)
	timer := 0
	sccCnt := 0

	var dfs func(int)
	dfs = func(u int) {
		dfn[u] = timer
		low[u] = timer
		timer++
		stk = append(stk, u)
		inStk[u] = true
		for _, v := range adj[u] {
			if dfn[v] == -1 {
				dfs(v)
				if low[v] < low[u] {
					low[u] = low[v]
				}
			} else if inStk[v] {
				if dfn[v] < low[u] {
					low[u] = dfn[v]
				}
			}
		}
		if low[u] == dfn[u] {
			for {
				v := stk[len(stk)-1]
				stk = stk[:len(stk)-1]
				inStk[v] = false
				comp[v] = sccCnt
				if v == u {
					break
				}
			}
			sccCnt++
		}
	}

	for i := 0; i < 2*n; i++ {
		if dfn[i] == -1 {
			dfs(i)
		}
	}

	reachSCC := make([][]uint64, sccCnt)
	for i := 0; i < sccCnt; i++ {
		reachSCC[i] = make([]uint64, (sccCnt+63)/64)
		reachSCC[i][i/64] |= 1 << (i % 64)
	}

	sccNodes := make([][]int, sccCnt)
	for u := 0; u < 2*n; u++ {
		sccNodes[comp[u]] = append(sccNodes[comp[u]], u)
	}

	dag := make([][]int, sccCnt)
	seen := make([]int, sccCnt)
	for i := range seen {
		seen[i] = -1
	}

	for c := 0; c < sccCnt; c++ {
		seen[c] = c
		for _, u := range sccNodes[c] {
			for _, v := range adj[u] {
				cv := comp[v]
				if seen[cv] != c {
					seen[cv] = c
					dag[c] = append(dag[c], cv)
				}
			}
		}
	}

	for c := 0; c < sccCnt; c++ {
		for _, nxt := range dag[c] {
			for k := 0; k < len(reachSCC[c]); k++ {
				reachSCC[c][k] |= reachSCC[nxt][k]
			}
		}
	}

	return reachSCC, comp
}

func canReach(u, v int, reachSCC [][]uint64, comp []int) bool {
	cu, cv := comp[u], comp[v]
	return (reachSCC[cu][cv/64] & (1 << (cv % 64))) != 0
}

func printAss(ass []int) {
	strs := make([]string, len(ass))
	for i, v := range ass {
		strs[i] = string('0' + byte(v))
	}
	fmt.Println(strings.Join(strs, " "))
}

func main() {
	buffer, _ = ioutil.ReadAll(os.Stdin)
	n := nextInt()
	if n == 0 {
		return
	}
	m1 := nextInt()
	m2 := nextInt()

	edgesF := make([][]int, m1)
	for i := 0; i < m1; i++ {
		u, v := nextInt(), nextInt()
		edgesF[i] = []int{lit(u), lit(v)}
	}

	edgesG := make([][]int, m2)
	for i := 0; i < m2; i++ {
		u, v := nextInt(), nextInt()
		edgesG[i] = []int{lit(u), lit(v)}
	}

	satF, assF := solve2SAT(n, edgesF, nil)
	satG, assG := solve2SAT(n, edgesG, nil)

	if !satF && !satG {
		fmt.Println("SIMILAR")
		return
	}
	if !satF && satG {
		printAss(assG)
		return
	}
	if satF && !satG {
		printAss(assF)
		return
	}

	reachF, compF := buildReach(n, edgesF)
	for _, clause := range edgesG {
		A, B := clause[0], clause[1]
		if canReach(neg(A), B, reachF, compF) || canReach(neg(A), A, reachF, compF) || canReach(neg(B), B, reachF, compF) {
			continue
		}
		_, ass := solve2SAT(n, edgesF, []int{neg(A), neg(B)})
		printAss(ass)
		return
	}

	reachG, compG := buildReach(n, edgesG)
	for _, clause := range edgesF {
		A, B := clause[0], clause[1]
		if canReach(neg(A), B, reachG, compG) || canReach(neg(A), A, reachG, compG) || canReach(neg(B), B, reachG, compG) {
			continue
		}
		_, ass := solve2SAT(n, edgesG, []int{neg(A), neg(B)})
		printAss(ass)
		return
	}

	fmt.Println("SIMILAR")
}
`

func buildOracle() (string, error) {
	tmp, err := os.CreateTemp("", "oracleF_*.go")
	if err != nil {
		return "", err
	}
	if _, err := tmp.WriteString(refSourceF); err != nil {
		tmp.Close()
		return "", err
	}
	tmp.Close()
	defer os.Remove(tmp.Name())
	oracle := filepath.Join(os.TempDir(), "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, tmp.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesFRaw = `3 2 2 -3 -2 -3 -3 2 -3 -1 2
1 2 0 1 -1 1 1
2 1 1 -1 1 -2 -2
2 0 2 1 1 -2 2
4 1 3 -2 4 4 -4 2 -1 3 1
4 1 3 2 -3 -3 -4 -1 1 3 -3
4 0 3 -3 -3 1 -1 -1 2
3 0 0
2 0 1 -1 2
4 1 0 -4 -2
4 3 0 -2 -1 4 4 -2 -4
4 1 0 3 4
4 2 1 1 2 -3 4 3 3
2 0 2 -2 -1 -2 -2
4 3 1 3 -4 4 -1 1 -3 4 -2
2 0 3 2 2 1 -2 -1 -1
2 2 3 1 -2 1 -1 -1 -1 2 -2 2 -1
2 0 3 -2 -2 -2 2 1 2
2 1 0 -1 2
3 1 0 2 -1
1 3 0 1 1 -1 -1 -1 -1
1 1 0 1 1
3 2 2 -1 3 -3 -3 -1 2 1 3
3 2 0 -3 3 -3 3
4 0 0
4 1 3 2 1 4 -3 1 -3 -1 4
2 0 1 -2 -2
3 2 3 -1 1 1 -2 1 -2 -3 -3 2 1
4 2 1 1 -2 4 3 -2 1
2 1 2 -1 -2 -1 1 2 -1
3 2 0 -3 2 -2 3
1 3 2 -1 1 -1 1 -1 1 1 1 1 -1
1 3 0 1 -1 1 1 -1 1
3 2 3 -3 3 1 2 3 -2 3 1 2 -2
4 2 0 -3 2 -1 1
4 2 1 1 3 -4 -1 2 -1
4 3 1 1 -1 -4 -1 -3 -4 -2 -2
1 3 3 -1 -1 -1 1 -1 1 -1 1 -1 1 -1 1
1 0 2 1 -1 1 -1
2 0 3 1 -2 -1 -1 1 -1
3 2 3 -1 -2 2 -3 -1 2 2 -1 2 2
1 1 3 -1 -1 1 -1 -1 1 -1 -1
4 1 3 1 2 3 1 1 4 -3 2
4 3 2 -2 1 -2 -2 3 -1 1 -2 3 4
3 2 2 -2 -2 3 1 3 -2 -1 3
3 1 0 3 -2
1 3 1 1 -1 -1 -1 1 -1 1 1
1 2 1 -1 -1 1 -1 -1 1
4 0 1 -1 -4
4 1 3 3 -1 -2 2 2 -2 1 4
3 3 2 -3 -1 2 2 -2 -1 3 2 -1 -3
3 1 1 -2 -2 3 -2
1 2 3 1 -1 1 1 1 1 -1 1 -1 -1
4 3 1 -4 4 -1 2 -1 -3 -2 2
3 2 1 1 -1 -2 1 2 -1
3 2 1 1 2 2 2 1 1
4 3 3 1 4 1 4 -4 1 3 3 1 4 1 -3
3 0 3 3 -1 -2 -3 -3 3
1 0 2 1 1 -1 1
4 2 1 -3 -2 -4 2 -3 -3
1 3 1 -1 -1 1 -1 -1 -1 -1 1
3 3 3 -2 3 -3 3 2 -1 3 3 3 -2 -3 -1
1 2 0 1 1 1 -1
4 0 3 -4 -1 3 -3 4 2
4 2 0 1 3 2 -1
4 1 3 -3 -2 -2 -4 1 3 -3 1
4 0 3 -3 3 1 1 -1 -3
3 3 1 2 -3 2 3 1 2 2 -2
2 3 0 -1 -2 -1 -1 -1 -2
1 1 3 -1 1 1 -1 -1 -1 -1 -1
3 3 1 2 -3 2 1 1 -3 -1 3
4 3 0 -2 -3 1 -1 2 -1
1 3 3 1 1 1 1 -1 -1 1 1 1 -1 1 -1
4 2 1 -2 -3 -1 -1 -3 4
2 1 3 2 -1 -1 -2 1 2 -1 -1
2 2 2 -1 1 -2 -1 -1 -1 -2 2
1 2 3 -1 -1 1 -1 1 1 1 -1 -1 -1
1 2 3 1 -1 -1 -1 1 1 -1 1 -1 1
1 1 0 1 -1
1 0 1 1 1
1 3 3 -1 1 1 1 1 -1 -1 -1 1 -1 -1 1
2 1 2 2 1 1 -2 2 -2
3 2 2 2 -2 -3 -2 -3 -3 2 3
4 3 1 4 -1 4 -1 -1 2 -3 1
3 3 3 1 3 2 -1 -3 -2 -1 -3 -2 2 2 1
4 1 3 -2 1 4 1 -3 -3 2 -3
1 3 0 1 1 -1 -1 -1 1
2 2 1 1 2 2 -1 -1 -1
2 2 1 1 1 -2 -2 -1 -1
3 3 2 2 1 -1 2 -1 2 -2 2 -2 3
4 3 3 -2 -3 -4 -4 -1 -3 -1 3 1 -4 4 -3
2 1 3 2 -2 1 -2 1 1 2 1
4 0 0
2 1 2 -2 -1 2 -1 -1 2
2 3 3 -2 2 -2 1 2 -2 -1 1 -1 -2 -2 -1
2 1 0 2 2
3 2 2 1 3 1 3 3 -2 -3 -3
2 2 1 -1 1 2 2 -2 1
3 3 2 -1 -1 -2 -3 2 -3 3 -3 -1 3
4 1 0 3 3`

	scanner := bufio.NewScanner(strings.NewReader(testcasesFRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		n := atoi(fields[0])
		m1 := atoi(fields[1])
		m2 := atoi(fields[2])
		if len(fields) != 3+2*(m1+m2) {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		i := 3
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, m1, m2))
		for j := 0; j < m1; j++ {
			input.WriteString(fields[i])
			input.WriteByte(' ')
			input.WriteString(fields[i+1])
			input.WriteByte('\n')
			i += 2
		}
		for j := 0; j < m2; j++ {
			input.WriteString(fields[i])
			input.WriteByte(' ')
			input.WriteString(fields[i+1])
			input.WriteByte('\n')
			i += 2
		}
		inputStr := input.String()

		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(inputStr)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(inputStr)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
