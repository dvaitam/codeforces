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

type pt struct{ x, y int }

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// expected uses the reference solution 374E.go to compute the answer
func expected(input string) (string, error) {
	return runProg("./374E.go", input)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(3) + 2
		m := rng.Intn(3) + 2
		blue := make([]pt, n)
		red := make([]pt, m)
		for i := 0; i < n; i++ {
			blue[i] = pt{rng.Intn(7) - 3, rng.Intn(7) - 3}
		}
		for i := 0; i < m; i++ {
			red[i] = pt{rng.Intn(7) - 3, rng.Intn(7) - 3}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, p := range blue {
			sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
		}
		for _, p := range red {
			sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
		}
		input := sb.String()
		exp, err := expected(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
