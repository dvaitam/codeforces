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

type Test struct {
	n, m int
	grid []string
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.m))
	for _, row := range t.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildOracle() (string, error) {
	ref := "oracleD"
	cmd := exec.Command("go", "build", "-o", ref, "1689D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2 // 2..9
		m := rng.Intn(8) + 2
		grid := make([]string, n)
		hasB := false
		for r := 0; r < n; r++ {
			row := make([]byte, m)
			for c := 0; c < m; c++ {
				if rng.Intn(3) == 0 { // 1/3 chance B
					row[c] = 'B'
					hasB = true
				} else {
					row[c] = 'W'
				}
			}
			grid[r] = string(row)
		}
		if !hasB {
			// ensure at least one B
			grid[0] = "B" + grid[0][1:]
		}
		tests[i] = Test{n, m, grid}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("oracle runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
