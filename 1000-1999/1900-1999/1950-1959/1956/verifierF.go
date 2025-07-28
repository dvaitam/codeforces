package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func find(parent []int, x int) int {
	for parent[x] != x {
		parent[x] = parent[parent[x]]
		x = parent[x]
	}
	return x
}

func remove(parent []int, x int) {
	parent[x] = find(parent, x+1)
}

func processRange(parent []int, visited []bool, l, r []int, n int, u int, L int, R int, queue *[]int) {
	if L < 1 {
		L = 1
	}
	if R > n {
		R = n
	}
	for L <= R {
		j := find(parent, L)
		if j > R {
			break
		}
		diff := u - j
		if diff < 0 {
			diff = -diff
		}
		if diff >= l[u]+l[j] && diff <= r[u]+r[j] {
			visited[j] = true
			remove(parent, j)
			*queue = append(*queue, j)
			L = j
		} else {
			L = j + 1
		}
	}
}

func solveCase(n int, l, r []int) string {
	parent := make([]int, n+2)
	for i := 1; i <= n+1; i++ {
		parent[i] = i
	}
	visited := make([]bool, n+2)
	components := 0
	queue := make([]int, 0)
	for i := 1; i <= n; i++ {
		if visited[i] {
			continue
		}
		components++
		visited[i] = true
		remove(parent, i)
		queue = append(queue, i)
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			L1 := u - r[u]
			R1 := u - l[u]
			processRange(parent, visited, l, r, n, u, L1, R1, &queue)
			L2 := u + l[u]
			R2 := u + r[u]
			processRange(parent, visited, l, r, n, u, L2, R2, &queue)
		}
	}
	return fmt.Sprint(components)
}

func genCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var in bytes.Buffer
	var out bytes.Buffer
	fmt.Fprintf(&in, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(4) + 1
		fmt.Fprintf(&in, "%d\n", n)
		l := make([]int, n+1)
		rArr := make([]int, n+1)
		for j := 1; j <= n; j++ {
			l[j] = rng.Intn(3)
			rArr[j] = l[j] + rng.Intn(3)
			fmt.Fprintf(&in, "%d %d\n", l[j], rArr[j])
		}
		out.WriteString(solveCase(n, l, rArr))
		if i+1 < t {
			out.WriteByte('\n')
		}
	}
	return in.String(), strings.TrimSpace(out.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
