package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pair struct{ v, p int }

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func genTree(n int) []int {
	parents := make([]int, n+1)
	level := []int{1}
	next := []int{}
	idx := 2
	for len(level) > 0 && idx <= n {
		for _, p := range level {
			child := rand.Intn(3) + 1
			for c := 0; c < child && idx <= n; c++ {
				parents[idx] = p
				next = append(next, idx)
				idx++
			}
			if idx > n {
				break
			}
		}
		level = next
		next = []int{}
	}
	for idx <= n {
		parents[idx] = 1
		idx++
	}
	return parents
}

func genCase() string {
	n := rand.Intn(10) + 1
	q := rand.Intn(n) + 1
	A := make([]int, n+1)
	for i := 1; i <= n; i++ {
		A[i] = rand.Intn(5) + 1
	}
	parents := genTree(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", A[i]))
	}
	sb.WriteByte('\n')
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", parents[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		x := rand.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d\n", x))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	refBin := "./1425I_ref"
	if err := exec.Command("go", "build", "-o", refBin, "1425I.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	rand.Seed(42)
	for i := 0; i < 100; i++ {
		input := genCase()
		expect, err := run(refBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed on test", i+1, ":", err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "program failed on test", i+1, ":", err)
			os.Exit(1)
		}
		if expect != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
