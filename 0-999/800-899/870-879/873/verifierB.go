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

type TestB struct {
	n int
	s string
}

func (t TestB) Input() string {
	return fmt.Sprintf("%d\n%s\n", t.n, t.s)
}

func expectedB(t TestB) string {
	diffPos := map[int]int{0: 0}
	diff := 0
	best := 0
	for i := 1; i <= t.n; i++ {
		if t.s[i-1] == '1' {
			diff++
		} else {
			diff--
		}
		if pos, ok := diffPos[diff]; ok {
			if i-pos > best {
				best = i - pos
			}
		} else {
			diffPos[diff] = i
		}
	}
	return strconv.Itoa(best)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func genTests() []TestB {
	rand.Seed(2)
	tests := make([]TestB, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(100) + 1
		b := make([]byte, n)
		for j := range b {
			if rand.Intn(2) == 0 {
				b[j] = '0'
			} else {
				b[j] = '1'
			}
		}
		tests = append(tests, TestB{n: n, s: string(b)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		exp := strings.TrimSpace(expectedB(tc))
		gotRaw, err := run(bin, tc.Input())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, gotRaw)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotRaw)
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
