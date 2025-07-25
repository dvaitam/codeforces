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

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := "587C.go"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(8) + 2
		m := rng.Intn(n) + 1
		q := rng.Intn(5) + 1
		edges := make([][2]int, n-1)
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges[i-2] = [2]int{p, i}
		}
		people := make([]int, m)
		for i := 0; i < m; i++ {
			people[i] = rng.Intn(n) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for i, c := range people {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", c))
		}
		sb.WriteString("\n")
		for i := 0; i < q; i++ {
			v := rng.Intn(n) + 1
			u := rng.Intn(n) + 1
			a := rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d\n", v, u, a))
		}
		input := sb.String()
		expected, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error running reference: %v\n", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
