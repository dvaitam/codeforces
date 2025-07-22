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

type TestCase struct {
	a   []int
	b   []int
	ans int
}

func compute(a, b []int) int {
	curr := 0
	maxc := 0
	for i := range a {
		curr -= a[i]
		curr += b[i]
		if curr > maxc {
			maxc = curr
		}
	}
	return maxc
}

func genCase() TestCase {
	n := rand.Intn(8) + 2 // 2..9
	a := make([]int, n)
	b := make([]int, n)
	curr := 0
	for i := 0; i < n-1; i++ {
		if curr > 0 {
			a[i] = rand.Intn(curr + 1)
		}
		curr -= a[i]
		add := rand.Intn(11)
		b[i] = add
		curr += add
	}
	a[n-1] = curr
	b[n-1] = 0
	return TestCase{a, b, compute(a, b)}
}

func genCases(n int) []TestCase {
	rand.Seed(time.Now().UnixNano())
	cs := make([]TestCase, n)
	for i := 0; i < n; i++ {
		cs[i] = genCase()
	}
	return cs
}

func buildInput(tc TestCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(tc.a))
	for i := 0; i < len(tc.a); i++ {
		fmt.Fprintf(&sb, "%d %d\n", tc.a[i], tc.b[i])
	}
	return sb.String()
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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases(100)
	for i, tc := range cases {
		input := buildInput(tc)
		expected := fmt.Sprint(tc.ans)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
