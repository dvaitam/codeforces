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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func expected(a []int) string {
	if len(a) <= 1 {
		return "-1"
	}
	allEqual := true
	for i := 1; i < len(a); i++ {
		if a[i] != a[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return "-1"
	}
	g := 0
	for i := 1; i < len(a); i++ {
		diff := a[i] - a[0]
		if diff < 0 {
			diff = -diff
		}
		g = gcd(g, diff)
	}
	return fmt.Sprintf("%d", g)
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(7) + 2
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(41) - 20
	}
	return a
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		arr := generateCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d\n", len(arr))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expectedOutput := expected(arr)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expectedOutput {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, expectedOutput, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
