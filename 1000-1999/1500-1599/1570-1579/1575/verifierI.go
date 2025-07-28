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
	n := rand.Intn(4) + 2
	q := rand.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rand.Intn(5)-2))
	}
	sb.WriteByte('\n')
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", p, i))
	}
	hasQuery := false
	for i := 0; i < q; i++ {
		if rand.Intn(2) == 0 {
			u := rand.Intn(n) + 1
			c := rand.Intn(5) - 2
			sb.WriteString(fmt.Sprintf("1 %d %d\n", u, c))
		} else {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("2 %d %d\n", u, v))
			hasQuery = true
		}
	}
	if !hasQuery {
		sb.WriteString(fmt.Sprintf("2 %d %d\n", 1, n))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierI.go /path/to/binary\n")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "refI"
	if err := exec.Command("go", "build", "-o", ref, "1575I.go").Run(); err != nil {
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
