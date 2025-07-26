package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	in  string
	out string
}

func solveD(n int, s1, s2 string) string {
	parent := make([]int, 27)
	for i := 1; i <= 26; i++ {
		parent[i] = i
	}
	var ops [][2]int
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		pa := find(a)
		pb := find(b)
		if pa != pb {
			parent[pa] = pb
		}
	}
	for i := 0; i < n; i++ {
		a := int(s1[i] - 'a' + 1)
		b := int(s2[i] - 'a' + 1)
		if a != b && find(a) != find(b) {
			ops = append(ops, [2]int{a, b})
			union(a, b)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d", len(ops))
	for _, op := range ops {
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%c %c", rune(op[0]-1+'a'), rune(op[1]-1+'a'))
	}
	return sb.String()
}

func generateTests() []test {
	rand.Seed(4)
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		b := []byte("abcdef")
		s1b := make([]byte, n)
		s2b := make([]byte, n)
		for j := 0; j < n; j++ {
			s1b[j] = b[rand.Intn(len(b))]
			s2b[j] = b[rand.Intn(len(b))]
		}
		s1 := string(s1b)
		s2 := string(s2b)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n%s\n%s\n", n, s1, s2)
		tests = append(tests, test{in: sb.String(), out: solveD(n, s1, s2)})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := run(bin, t.in)
		if err != nil {
			fmt.Printf("Test %d failed to run: %v\n", i+1, err)
			fmt.Print(out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != t.out {
			fmt.Printf("Test %d failed. Expected %q, got %q. Input:\n%s", i+1, t.out, out, t.in)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}
