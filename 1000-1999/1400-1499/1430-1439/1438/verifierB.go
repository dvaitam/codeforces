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

func solveCase(b []int) string {
	seen := make(map[int]bool)
	for _, v := range b {
		if seen[v] {
			return "YES"
		}
		seen[v] = true
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(50) + 1
	var in strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	expLines := make([]string, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(50) + 2
		in.WriteString(fmt.Sprintf("%d\n", n))
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if j > 0 {
				in.WriteByte(' ')
			}
			if rng.Intn(5) == 0 && j > 0 {
				arr[j] = arr[rng.Intn(j)]
			} else {
				arr[j] = rng.Intn(1_000_000_000)
			}
			in.WriteString(fmt.Sprint(arr[j]))
		}
		in.WriteByte('\n')
		expLines[i] = solveCase(arr)
	}
	return in.String(), strings.Join(expLines, "\n")
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	outLines := strings.Fields(strings.TrimSpace(out.String()))
	expLines := strings.Fields(strings.TrimSpace(expected))
	if len(outLines) != len(expLines) {
		return fmt.Errorf("expected %d lines got %d", len(expLines), len(outLines))
	}
	for i := range expLines {
		if strings.ToUpper(outLines[i]) != expLines[i] {
			return fmt.Errorf("line %d expected %q got %q", i+1, expLines[i], outLines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
