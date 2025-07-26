package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct{ input string }

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1105D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genBoard(r *rand.Rand, n, m, p int) []string {
	board := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if r.Intn(5) == 0 {
				row[j] = '#'
			} else {
				row[j] = '.'
			}
		}
		board[i] = row
	}
	// place at least one castle for each player
	for pid := 0; pid < p; pid++ {
		for {
			x := r.Intn(n)
			y := r.Intn(m)
			if board[x][y] == '.' {
				board[x][y] = byte('1' + pid)
				break
			}
		}
	}
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		lines[i] = string(board[i])
	}
	return lines
}

func genTests() []Test {
	r := rand.New(rand.NewSource(3))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		n := r.Intn(5) + 1
		m := r.Intn(5) + 1
		p := r.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, p)
		for j := 1; j <= p; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", r.Intn(3)+1)
		}
		sb.WriteByte('\n')
		lines := genBoard(r, n, m, p)
		for _, line := range lines {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
		tests = append(tests, Test{sb.String()})
	}
	// some fixed simple cases
	tests = append(tests,
		Test{"1 1 1\n1\n1\n"},
		Test{"2 2 2\n1 1\n1.\n.2\n"},
		Test{"3 3 1\n2\n...\n.#.\n..1\n"},
		Test{"2 3 2\n1 1\n1#2\n...\n"},
		Test{"3 3 2\n3 3\n1..\n.#.\n..2\n"},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
