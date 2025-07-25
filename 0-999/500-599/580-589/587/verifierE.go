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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := "587E.go"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(6) + 1
		q := rng.Intn(6) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(32)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString("\n")
		for i := 0; i < q; i++ {
			t := rng.Intn(2) + 1
			if t == 1 {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				k := rng.Intn(32)
				sb.WriteString(fmt.Sprintf("1 %d %d %d\n", l, r, k))
			} else {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				sb.WriteString(fmt.Sprintf("2 %d %d\n", l, r))
			}
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
