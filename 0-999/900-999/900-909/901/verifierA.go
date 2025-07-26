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

func solveA(k int, a []int) string {
	n := k + 1
	ambIdx := -1
	for i := 1; i < n; i++ {
		if a[i-1] > 1 && a[i] > 1 {
			ambIdx = i
			break
		}
	}
	if ambIdx == -1 {
		return "perfect"
	}
	var sb strings.Builder
	sb.WriteString("ambiguous\n")
	sum := 0
	for i := 0; i < n; i++ {
		for j := 0; j < a[i]; j++ {
			fmt.Fprintf(&sb, "%d ", sum)
		}
		sum += a[i]
	}
	sb.WriteByte('\n')
	sum = 0
	for i := 0; i < n; i++ {
		for j := 0; j < a[i]; j++ {
			if i == ambIdx && j == 0 {
				fmt.Fprintf(&sb, "%d ", sum-1)
			} else {
				fmt.Fprintf(&sb, "%d ", sum)
			}
		}
		sum += a[i]
	}
	return strings.TrimSpace(sb.String())
}

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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		k := rng.Intn(5) + 1
		n := k + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(3) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", k)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveA(k, a)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
