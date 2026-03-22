package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solveL = 2520

func solve(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var out bytes.Buffer
	writer := bufio.NewWriter(&out)

	var n int
	fmt.Fscan(reader, &n)

	k := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &k[i])
	}

	m := make([]int, n)
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &m[i])
		adj[i] = make([]int, m[i])
		for j := 0; j < m[i]; j++ {
			fmt.Fscan(reader, &adj[i][j])
			adj[i][j]--
		}
	}

	totalStates := n * solveL
	nextState := make([]int, totalStates)

	for u := 0; u < n; u++ {
		effK := k[u] % solveL
		if effK < 0 {
			effK += solveL
		}
		modM := m[u]
		baseState := u * solveL

		for r := 0; r < solveL; r++ {
			val := r + effK
			if val >= solveL {
				val -= solveL
			}
			idx := val % modM
			v := adj[u][idx]
			nextState[baseState+r] = v*solveL + val
		}
	}

	ans := make([]int, totalStates)
	vis := make([]byte, totalStates)
	onStack := make([]int, totalStates)
	for i := range onStack {
		onStack[i] = -1
	}

	path := make([]int, 0, 4096)
	seen := make([]bool, n)
	seenList := make([]int, 0, n)

	for i := 0; i < totalStates; i++ {
		if vis[i] != 0 {
			continue
		}

		curr := i
		path = path[:0]

		for vis[curr] == 0 {
			vis[curr] = 1
			onStack[curr] = len(path)
			path = append(path, curr)
			curr = nextState[curr]
		}

		var result int
		if vis[curr] == 1 {
			startIdx := onStack[curr]
			cnt := 0
			for j := startIdx; j < len(path); j++ {
				node := path[j]
				u := node / solveL
				if !seen[u] {
					seen[u] = true
					seenList = append(seenList, u)
					cnt++
				}
			}
			result = cnt
			for _, u := range seenList {
				seen[u] = false
			}
			seenList = seenList[:0]
		} else {
			result = ans[curr]
		}

		for _, node := range path {
			ans[node] = result
			vis[node] = 2
			onStack[node] = -1
		}
	}

	var q int
	fmt.Fscan(reader, &q)
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		u := x - 1
		r := y % solveL
		if r < 0 {
			r += solveL
		}
		fmt.Fprintln(writer, ans[u*solveL+r])
	}

	writer.Flush()
	return strings.TrimSpace(out.String())
}

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randCase(rng *rand.Rand) Test {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(rng.Intn(11) - 5))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		m := rng.Intn(3) + 1
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(rng.Intn(n) + 1))
		}
		sb.WriteByte('\n')
	}
	q := rng.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(11) - 5
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return Test{sb.String()}
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(5))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		tests = append(tests, randCase(rng))
	}
	tests = append(tests, Test{"1\n0\n1\n1\n1\n1 0\n"})
	tests = append(tests, Test{"2\n1 1\n1\n1\n1\n2\n1\n1 0\n"})
	tests = append(tests, Test{"3\n0 0 0\n1\n1\n1\n2\n1 2\n3\n1 0\n2 -1\n3 2\n"})
	tests = append(tests, Test{"1\n5\n1\n1\n2\n1\n1 1\n"})
	tests = append(tests, Test{"2\n-1 2\n1\n2\n1 2\n1\n1 0\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		exp := solve(tc.input)
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
