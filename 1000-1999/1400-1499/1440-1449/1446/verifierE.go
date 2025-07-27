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

type testCase struct {
	pts [][2]int
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(t testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(t.pts)))
	for _, p := range t.pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return runProg("1446E.go", sb.String())
}

func generateTests() []testCase {
	tests := []testCase{
		{pts: [][2]int{{0, 0}, {1, 0}}},
		{pts: [][2]int{{0, 0}, {1, 0}, {0, 1}}},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := r.Intn(5) + 1
		m := make(map[[2]int]bool)
		pts := make([][2]int, 0, n)
		for len(pts) < n {
			x := r.Intn(20)
			y := r.Intn(20)
			p := [2]int{x, y}
			if !m[p] {
				m[p] = true
				pts = append(pts, p)
			}
		}
		tests = append(tests, testCase{pts: pts})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(t.pts)))
		for _, p := range t.pts {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		want, err := expected(t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
