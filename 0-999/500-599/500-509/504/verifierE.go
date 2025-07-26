package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func pathStr(adj [][]int, labels string, u, v int) string {
	n := len(adj)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	queue := []int{u}
	parent[u] = u
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		if x == v {
			break
		}
		for _, y := range adj[x] {
			if parent[y] == -1 {
				parent[y] = x
				queue = append(queue, y)
			}
		}
	}
	path := []int{v}
	for path[len(path)-1] != u {
		path = append(path, parent[path[len(path)-1]])
	}
	for i := 0; i < len(path)/2; i++ {
		path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
	}
	var sb strings.Builder
	for _, x := range path {
		sb.WriteByte(labels[x])
	}
	return sb.String()
}

func lcp(a, b string) int {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return i
}

func genTest(rng *rand.Rand) (string, []int) {
	n := rng.Intn(5) + 1
	adj := make([][]int, n)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		adj[i] = append(adj[i], p)
		adj[p] = append(adj[p], i)
	}
	letters := make([]byte, n)
	for i := 0; i < n; i++ {
		letters[i] = byte('a' + rng.Intn(3))
	}
	m := rng.Intn(5) + 1
	queries := make([][4]int, m)
	for i := 0; i < m; i++ {
		for j := 0; j < 4; j++ {
			queries[i][j] = rng.Intn(n)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	sb.WriteString(string(letters) + "\n")
	for i := 1; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", i+1, adj[i][0]+1))
	}
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		a, b, c, d := queries[i][0]+1, queries[i][1]+1, queries[i][2]+1, queries[i][3]+1
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a, b, c, d))
	}
	answers := make([]int, m)
	for i, q := range queries {
		s1 := pathStr(adj, string(letters), q[0], q[1])
		s2 := pathStr(adj, string(letters), q[2], q[3])
		answers[i] = lcp(s1, s2)
	}
	return sb.String(), answers
}

func check(out string, ans []int) error {
	fields := strings.Fields(out)
	if len(fields) != len(ans) {
		return fmt.Errorf("expected %d numbers", len(ans))
	}
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil || v != ans[i] {
			return fmt.Errorf("wrong answer")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, ans := genTest(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(out, ans); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
