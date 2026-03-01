package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n, a, b int
}

func (t Test) Input() string {
	return fmt.Sprintf("%d %d %d\n", t.n, t.a, t.b)
}

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(6) + 1
		a := rand.Intn(n) + 1
		b := rand.Intn(n) + 1
		tests = append(tests, Test{n: n, a: a, b: b})
	}
	return tests
}

func isImpossible(n, a, b int) bool {
	if a > 1 && b > 1 {
		return true
	}
	if a == 1 && b == 1 && (n == 2 || n == 3) {
		return true
	}
	return false
}

func countComponents(g [][]bool) int {
	n := len(g)
	vis := make([]bool, n)
	q := make([]int, 0, n)
	comps := 0

	for s := 0; s < n; s++ {
		if vis[s] {
			continue
		}
		comps++
		vis[s] = true
		q = q[:0]
		q = append(q, s)
		for head := 0; head < len(q); head++ {
			v := q[head]
			for to := 0; to < n; to++ {
				if g[v][to] && !vis[to] {
					vis[to] = true
					q = append(q, to)
				}
			}
		}
	}
	return comps
}

func validateOutput(out string, t Test) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}

	head := strings.TrimSpace(lines[0])
	if head != "YES" && head != "NO" {
		return fmt.Errorf("first line must be YES or NO, got %q", head)
	}

	impossible := isImpossible(t.n, t.a, t.b)
	if head == "NO" {
		if !impossible {
			return fmt.Errorf("printed NO but a solution exists")
		}
		if len(lines) != 1 {
			return fmt.Errorf("NO output should not contain matrix lines")
		}
		return nil
	}

	if impossible {
		return fmt.Errorf("printed YES for impossible case")
	}

	if len(lines) != t.n+1 {
		return fmt.Errorf("expected %d matrix lines, got %d", t.n, len(lines)-1)
	}

	g := make([][]bool, t.n)
	for i := 0; i < t.n; i++ {
		row := strings.TrimSpace(lines[i+1])
		if len(row) != t.n {
			return fmt.Errorf("row %d has length %d, expected %d", i+1, len(row), t.n)
		}
		g[i] = make([]bool, t.n)
		for j := 0; j < t.n; j++ {
			c := row[j]
			if c != '0' && c != '1' {
				return fmt.Errorf("row %d contains non-binary character %q", i+1, c)
			}
			if i == j && c != '0' {
				return fmt.Errorf("diagonal element (%d,%d) is not zero", i+1, j+1)
			}
			g[i][j] = c == '1'
		}
	}

	for i := 0; i < t.n; i++ {
		for j := i + 1; j < t.n; j++ {
			if g[i][j] != g[j][i] {
				return fmt.Errorf("matrix is not symmetric at (%d,%d)", i+1, j+1)
			}
		}
	}

	ca := countComponents(g)
	if ca != t.a {
		return fmt.Errorf("graph has %d components, expected %d", ca, t.a)
	}

	comp := make([][]bool, t.n)
	for i := 0; i < t.n; i++ {
		comp[i] = make([]bool, t.n)
		for j := 0; j < t.n; j++ {
			if i != j {
				comp[i][j] = !g[i][j]
			}
		}
	}
	cb := countComponents(comp)
	if cb != t.b {
		return fmt.Errorf("complement has %d components, expected %d", cb, t.b)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validateOutput(got, tc); err != nil {
			fmt.Printf("Test %d failed\nInput:%sOutput:%sReason:%v\n", i+1, input, got, err)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
