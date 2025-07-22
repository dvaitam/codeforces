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

var dir = []byte{'<', '^', '>'}

func randBoard(rng *rand.Rand, n, m int) []string {
	board := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			b[j] = dir[rng.Intn(3)]
		}
		board[i] = string(b)
	}
	return board
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		q := rng.Intn(5) + 1
		board := randBoard(rng, n, m)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for i := 0; i < n; i++ {
			sb.WriteString(board[i])
			sb.WriteByte('\n')
		}
		for i := 0; i < q; i++ {
			if rng.Intn(2) == 0 {
				x := rng.Intn(n) + 1
				y := rng.Intn(m) + 1
				sb.WriteString(fmt.Sprintf("A %d %d\n", x, y))
			} else {
				x := rng.Intn(n) + 1
				y := rng.Intn(m) + 1
				c := string(dir[rng.Intn(3)])
				sb.WriteString(fmt.Sprintf("C %d %d %s\n", x, y, c))
			}
		}
		input := sb.String()
		expected, err := run("487D.go", input)
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
