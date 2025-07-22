package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func getSubtree(adj [][]int, root, parent int, res *[]int) {
	*res = append(*res, root)
	for _, v := range adj[root] {
		if v != parent {
			getSubtree(adj, v, root, res)
		}
	}
}

func solve(n, m int, edges [][2]int, ops [][2]int) []int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}
	lists := make([]map[int]struct{}, n+1)
	for i := 1; i <= n; i++ {
		lists[i] = make(map[int]struct{})
	}
	for idx, op := range ops {
		a, b := op[0], op[1]
		nodes := []int{}
		getSubtree(adj, a, 0, &nodes)
		for _, u := range nodes {
			lists[u][idx+1] = struct{}{}
		}
		nodes = nodes[:0]
		getSubtree(adj, b, 0, &nodes)
		for _, u := range nodes {
			lists[u][idx+1] = struct{}{}
		}
	}
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if i == j {
				continue
			}
			flag := false
			for x := range lists[i] {
				if _, ok := lists[j][x]; ok {
					flag = true
					break
				}
			}
			if flag {
				ans[i]++
			}
		}
	}
	return ans[1:]
}

func generateTests() []testCase {
	var tests []testCase
	n := 2
	m := 1
	for len(tests) < 100 {
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			edges[i] = [2]int{i + 1, i + 2}
		}
		ops := make([][2]int, m)
		for i := 0; i < m; i++ {
			ops[i] = [2]int{1, n}
		}
		ans := solve(n, m, edges, ops)
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for _, op := range ops {
			sb.WriteString(fmt.Sprintf("%d %d\n", op[0], op[1]))
		}
		outParts := make([]string, len(ans))
		for i, v := range ans {
			outParts[i] = fmt.Sprint(v)
		}
		outStr := strings.Join(outParts, " ")
		tests = append(tests, testCase{in: sb.String(), out: outStr})
		n++
		if n > 4 {
			n = 2
			m++
			if m > 2 {
				m = 1
			}
		}
	}
	return tests
}

func runTest(bin string, tc testCase) (string, error) {
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(tc.in))
	}()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if strings.HasSuffix(bin, ".go") {
		tmp, err := ioutil.TempFile("", "solbin*")
		if err != nil {
			fmt.Println("cannot create temp file:", err)
			os.Exit(1)
		}
		tmp.Close()
		exec.Command("go", "build", "-o", tmp.Name(), bin).Run()
		bin = tmp.Name()
		defer os.Remove(bin)
	}
	tests := generateTests()
	for i, tc := range tests {
		got, err := runTest(bin, tc)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Printf("wrong answer on test %d\ninput: %sexpected: %s\ngot: %s\n", i+1, tc.in, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
