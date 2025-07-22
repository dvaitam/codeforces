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

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

func genGraph(rng *rand.Rand, n, m int) [][2]int {
	edges := make([][2]int, 0, m)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		exist := false
		for _, e := range edges {
			if (e[0] == a && e[1] == b) || (e[0] == b && e[1] == a) {
				exist = true
				break
			}
		}
		if !exist {
			edges = append(edges, [2]int{a, b})
		}
	}
	return edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 2
		m := rng.Intn(3) + n - 1
		q := rng.Intn(5) + 1
		weights := make([]int, n)
		for i := 0; i < n; i++ {
			weights[i] = rng.Intn(100) + 1
		}
		edges := genGraph(rng, n, m)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", weights[i]))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for i := 0; i < q; i++ {
			if rng.Intn(2) == 0 {
				a := rng.Intn(n) + 1
				b := rng.Intn(n) + 1
				sb.WriteString(fmt.Sprintf("A %d %d\n", a, b))
			} else {
				a := rng.Intn(n) + 1
				w := rng.Intn(100) + 1
				sb.WriteString(fmt.Sprintf("C %d %d\n", a, w))
			}
		}
		input := sb.String()
		expected, err := run("487E.go", input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal reference failed on case %d: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\ninput:\n%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
