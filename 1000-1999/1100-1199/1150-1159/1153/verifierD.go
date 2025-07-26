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

type Test struct {
	input    string
	expected string
}

func generateTests() []Test {
	rng := rand.New(rand.NewSource(45))
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		op := make([]int, n+1)
		parent := make([]int, n+1)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for v := 1; v <= n; v++ {
			op[v] = rng.Intn(2)
			if v > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(op[v]))
		}
		sb.WriteByte('\n')
		for v := 2; v <= n; v++ {
			p := rng.Intn(v-1) + 1
			parent[v] = p
			if v > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(p))
		}
		if n > 1 {
			sb.WriteByte('\n')
		}
		input := sb.String()
		children := make([][]int, n+1)
		for v := 2; v <= n; v++ {
			children[parent[v]] = append(children[parent[v]], v)
		}
		f := make([]int, n+1)
		leafCount := 0
		for v := n; v >= 1; v-- {
			if len(children[v]) == 0 {
				f[v] = 1
				leafCount++
			} else if op[v] == 0 {
				sum := 0
				for _, c := range children[v] {
					sum += f[c]
				}
				f[v] = sum
			} else {
				mn := f[children[v][0]]
				for _, c := range children[v] {
					if f[c] < mn {
						mn = f[c]
					}
				}
				f[v] = mn
			}
		}
		res := leafCount - f[1] + 1
		tests = append(tests, Test{input: input, expected: strconv.Itoa(res)})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.input, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
