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

func solveE(n, m, k int, special []int) string {
	vis := make([]int, n+1)
	ss := special[0]
	for i, x := range special {
		vis[x] = 1
		if i == 0 {
			vis[x] = 2
		}
	}
	a := make([]int, 0, n)
	b := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if vis[i] != 2 {
			a = append(a, i)
		}
		if vis[i] == 0 {
			b = append(b, i)
		}
	}
	tot := (n - k) + ((n-1)*(n-2))/2
	if m > tot || k == n {
		return "-1"
	}
	var out strings.Builder
	edgesLeft := m
	for i := 0; i < len(a) && edgesLeft > 1; i++ {
		for j := i + 1; j < len(a) && edgesLeft > 1; j++ {
			fmt.Fprintf(&out, "%d %d\n", a[i], a[j])
			edgesLeft--
		}
	}
	for i := 0; i < edgesLeft; i++ {
		fmt.Fprintf(&out, "%d %d\n", ss, b[i])
	}
	return strings.TrimSpace(out.String())
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(8) + 3
		k := rng.Intn(n-1) + 2
		special := make([]int, 0, k)
		used := make(map[int]bool)
		for len(special) < k {
			x := rng.Intn(n) + 1
			if !used[x] {
				used[x] = true
				special = append(special, x)
			}
		}
		tot := (n - k) + ((n-1)*(n-2))/2
		m := rng.Intn(tot + 1) // may trigger -1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		for i, v := range special {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		expected := solveE(n, m, k, special)
		tests = append(tests, test{sb.String(), expected})
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
