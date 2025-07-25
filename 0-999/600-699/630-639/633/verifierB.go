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

func countZeros(n int) int {
	count := 0
	for n > 0 {
		n /= 5
		count += n
	}
	return count
}

func solveB(m int) string {
	low, high := 0, 5*m
	for low < high {
		mid := (low + high) / 2
		if countZeros(mid) < m {
			low = mid + 1
		} else {
			high = mid
		}
	}
	if countZeros(low) != m {
		return "0\n"
	}
	var sb strings.Builder
	sb.WriteString("5\n")
	for i := 0; i < 5; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(low + i))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCaseB(rng *rand.Rand) (string, string) {
	m := rng.Intn(100000) + 1
	input := fmt.Sprintf("%d\n", m)
	expect := solveB(m)
	return input, expect
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
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseB(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
