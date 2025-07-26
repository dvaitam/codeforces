package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildTree(n int, parents []int) [][]int {
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		children[p] = append(children[p], i)
	}
	for i := 1; i <= n; i++ {
		if len(children[i]) > 1 {
			sort.Ints(children[i])
		}
	}
	return children
}

func dfsOrder(children [][]int) ([]int, []int, []int) {
	n := len(children) - 1
	st := make([]int, n+1)
	ed := make([]int, n+1)
	id := make([]int, n+2)
	type frame struct{ node, idx int }
	stack := []frame{{1, 0}}
	num := 0
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.idx == 0 {
			num++
			st[top.node] = num
			id[num] = top.node
		}
		if top.idx < len(children[top.node]) {
			child := children[top.node][top.idx]
			top.idx++
			stack = append(stack, frame{child, 0})
		} else {
			ed[top.node] = num
			stack = stack[:len(stack)-1]
		}
	}
	return st, ed, id
}

func answerQueries(children [][]int, queries [][2]int) []int {
	st, ed, id := dfsOrder(children)
	res := make([]int, len(queries))
	for i, q := range queries {
		u := q[0]
		k := q[1]
		pos := st[u] + k - 1
		if pos <= ed[u] {
			res[i] = id[pos]
		} else {
			res[i] = -1
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(50) + 1
		q := rng.Intn(50) + 1
		parents := make([]int, n-1)
		for j := 0; j < n-1; j++ {
			parents[j] = rng.Intn(j+1) + 1
		}
		children := buildTree(n, parents)
		queries := make([][2]int, q)
		input := fmt.Sprintf("%d %d\n", n, q)
		for j := 0; j < n-1; j++ {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", parents[j])
		}
		input += "\n"
		for j := 0; j < q; j++ {
			u := rng.Intn(n) + 1
			k := rng.Intn(n) + 1
			queries[j] = [2]int{u, k}
			input += fmt.Sprintf("%d %d\n", u, k)
		}
		expectedVals := answerQueries(children, queries)
		expected := strings.TrimSpace(strings.Join(func() []string {
			tmp := make([]string, len(expectedVals))
			for i, v := range expectedVals {
				tmp[i] = fmt.Sprintf("%d", v)
			}
			return tmp
		}(), "\n"))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\nget\n%s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
