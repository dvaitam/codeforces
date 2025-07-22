package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solveB(n int, t, parent []int) []int {
	depth := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if t[i] != 1 || depth[i] != 0 {
			continue
		}
		path := make([]int, 0)
		u := i
		for u != 0 && depth[u] == 0 {
			path = append(path, u)
			u = parent[u]
		}
		base := 0
		if u != 0 {
			base = depth[u]
		}
		for j := len(path) - 1; j >= 0; j-- {
			base++
			depth[path[j]] = base
		}
	}
	best := 1
	maxd := 0
	for i := 1; i <= n; i++ {
		if t[i] == 1 && depth[i] > maxd {
			maxd = depth[i]
			best = i
		}
	}
	res := make([]int, 0, maxd)
	u := best
	for u != 0 {
		res = append(res, u)
		u = parent[u]
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		t := make([]int, n+1)
		hasHotel := false
		for i := 1; i <= n; i++ {
			if rng.Intn(2) == 1 {
				t[i] = 1
				hasHotel = true
			}
		}
		if !hasHotel {
			t[1] = 1
		}
		parent := make([]int, n+1)
		for i := 1; i <= n; i++ {
			parent[i] = rng.Intn(i) // 0..i-1 ensures no cycles
		}
		path := solveB(n, t, parent)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(t[i]))
		}
		sb.WriteByte('\n')
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(parent[i]))
		}
		sb.WriteByte('\n')
		var out strings.Builder
		fmt.Fprintf(&out, "%d\n", len(path))
		for i, v := range path {
			if i > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(strconv.Itoa(v))
		}
		out.WriteByte('\n')
		tests = append(tests, test{sb.String(), out.String()})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
