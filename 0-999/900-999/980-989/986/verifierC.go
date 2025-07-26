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

func solveCase(n int, arr []int) int {
	m := len(arr)
	visited := make(map[int]bool)
	comps := 0
	for i := 0; i < m; i++ {
		v := arr[i]
		if !visited[v] {
			comps++
			// bfs
			q := []int{v}
			visited[v] = true
			for len(q) > 0 {
				cur := q[0]
				q = q[1:]
				for j := 0; j < m; j++ {
					u := arr[j]
					if !visited[u] && (cur&u) == 0 {
						visited[u] = true
						q = append(q, u)
					}
				}
			}
		}
	}
	return comps
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	size := 1 << uint(n)
	m := rng.Intn(size) + 1
	used := make(map[int]bool)
	arr := make([]int, 0, m)
	for len(arr) < m {
		v := rng.Intn(size)
		if !used[v] {
			used[v] = true
			arr = append(arr, v)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := fmt.Sprintf("%d\n", solveCase(n, arr))
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
