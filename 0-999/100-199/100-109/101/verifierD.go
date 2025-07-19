package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseD struct {
	input string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []testCaseD{
		{input: "2\n1 2 1\n"},
		{input: "3\n1 2 1\n1 3 1\n"},
	}
	for i, t := range tests {
		expect := solveD(strings.NewReader(t.input))
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

type Edge struct{ to, v int }

type Arr struct {
	size int
	sum  int64
}

func solveD(r io.Reader) string {
	in := bufio.NewReader(r)
	var n int
	fmt.Fscan(in, &n)
	adj := make([][]Edge, n+1)
	for i := 1; i < n; i++ {
		var x, y, v int
		fmt.Fscan(in, &x, &y, &v)
		adj[x] = append(adj[x], Edge{to: y, v: v})
		adj[y] = append(adj[y], Edge{to: x, v: v})
	}
	sizeArr := make([]int, n+1)
	sumArr := make([]int64, n+1)
	timeArr := make([]int64, n+1)
	var dfs func(int, int)
	dfs = func(x, fa int) {
		sizeArr[x] = 1
		for _, e := range adj[x] {
			y, v := e.to, e.v
			if y == fa {
				continue
			}
			sumArr[y] = int64(v) * 2
			dfs(y, x)
			sumArr[x] += sumArr[y]
			sizeArr[x] += sizeArr[y]
		}
		var q []Arr
		for _, e := range adj[x] {
			y, v := e.to, e.v
			if y == fa {
				continue
			}
			timeArr[x] += timeArr[y] + int64(v)*int64(sizeArr[y])
			q = append(q, Arr{sizeArr[y], sumArr[y]})
		}
		sort.Slice(q, func(i, j int) bool { return q[i].sum*int64(q[j].size) < q[j].sum*int64(q[i].size) })
		var s int64
		for _, a := range q {
			timeArr[x] += int64(a.size) * s
			s += a.sum
		}
	}
	dfs(1, 0)
	res := float64(timeArr[1]) / float64(n-1)
	return fmt.Sprintf("%.10f\n", res)
}
