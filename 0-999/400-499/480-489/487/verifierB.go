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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		s := rng.Intn(100)
		l := rng.Intn(n) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(201) - 100
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, s, l))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected, err := run("487B.go", input)
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
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
