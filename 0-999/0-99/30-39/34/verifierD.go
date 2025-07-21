package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type queue []int

func (q *queue) push(x int)  { *q = append(*q, x) }
func (q *queue) pop() int    { v := (*q)[0]; *q = (*q)[1:]; return v }
func (q *queue) empty() bool { return len(*q) == 0 }

func generateCase() string {
	n := rand.Intn(50) + 2
	r1 := rand.Intn(n) + 1
	r2 := rand.Intn(n-1) + 1
	if r2 >= r1 {
		r2++
	}
	parent := make([]int, n+1)
	processed := []int{r1}
	for i := 1; i <= n; i++ {
		if i == r1 {
			continue
		}
		p := processed[rand.Intn(len(processed))]
		parent[i] = p
		processed = append(processed, i)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, r1, r2))
	first := true
	for i := 1; i <= n; i++ {
		if i == r1 {
			continue
		}
		if !first {
			sb.WriteByte(' ')
		} else {
			first = false
		}
		sb.WriteString(fmt.Sprintf("%d", parent[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(in string) string {
	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(bufio.ScanWords)
	var n, r1, r2 int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n)
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &r1)
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &r2)
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if i == r1 {
			continue
		}
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &parent[i])
	}
	adj := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if i == r1 {
			continue
		}
		p := parent[i]
		adj[i] = append(adj[i], p)
		adj[p] = append(adj[p], i)
	}
	resParent := make([]int, n+1)
	visited := make([]bool, n+1)
	q := queue{r2}
	visited[r2] = true
	for !q.empty() {
		u := q.pop()
		for _, v := range adj[u] {
			if !visited[v] {
				visited[v] = true
				resParent[v] = u
				q.push(v)
			}
		}
	}
	var sb strings.Builder
	first := true
	for i := 1; i <= n; i++ {
		if i == r2 {
			continue
		}
		if !first {
			sb.WriteByte(' ')
		} else {
			first = false
		}
		sb.WriteString(fmt.Sprintf("%d", resParent[i]))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := generateCase()
		exp := expected(tc)
		got, err := run(bin, tc)
		if err != nil {
			fmt.Printf("case %d: error executing binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
