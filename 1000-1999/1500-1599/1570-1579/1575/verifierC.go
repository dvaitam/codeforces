package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	if !strings.Contains(bin, "/") {
		bin = "./" + bin
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func genTest() string {
	n := rand.Intn(5) + 1
	m := rand.Intn(5) + 1
	primes := []int{2, 3, 5, 7}
	k := 1
	if rand.Intn(2) == 0 {
		k = primes[rand.Intn(len(primes))]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rand.Intn(10)))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierC.go /path/to/binary\n")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "refC"
	if err := exec.Command("go", "build", "-o", ref, "1575C.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rand.Seed(1)
	for i := 0; i < 100; i++ {
		input := genTest()
		exp, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed:", err)
			os.Exit(1)
		}
		got, err := runBinary(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed to run: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "test %d failed:\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
