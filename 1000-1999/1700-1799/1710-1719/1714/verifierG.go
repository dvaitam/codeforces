package main

import (
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

type Edge struct {
	to   int
	a, b int64
}

func run(bin, input string) (string, error) {
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

func solve(n int, edges [][]Edge) []int {
	res := make([]int, n+1)
	prefix := []int64{0}
	type Node struct {
		id   int
		aSum int64
		idx  int
	}
	stack := []Node{{id: 1, aSum: 0, idx: 0}}
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.idx == len(edges[top.id]) {
			if top.id != 1 {
				prefix = prefix[:len(prefix)-1]
			}
			stack = stack[:len(stack)-1]
			continue
		}
		e := edges[top.id][top.idx]
		top.idx++
		newASum := top.aSum + e.a
		prefix = append(prefix, prefix[len(prefix)-1]+e.b)
		pos := sort.Search(len(prefix), func(i int) bool { return prefix[i] > newASum }) - 1
		res[e.to] = pos
		stack = append(stack, Node{id: e.to, aSum: newASum, idx: 0})
	}
	return res
}

func genCase(rng *rand.Rand) (string, int, [][]Edge, []int) {
	n := rng.Intn(8) + 2
	edges := make([][]Edge, n+1)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		a := rng.Int63n(100) + 1
		b := rng.Int63n(100) + 1
		edges[p] = append(edges[p], Edge{to: i, a: a, b: b})
		sb.WriteString(fmt.Sprintf("%d %d %d\n", p, a, b))
	}
	ans := solve(n, edges)
	return sb.String(), n, edges, ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, edges, ans := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != n-1 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", i+1, n-1, len(fields))
			os.Exit(1)
		}
		for j, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil || v != ans[j+2] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\ninput:\n%s", i+1, ans[2:], fields, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
