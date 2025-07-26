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
	rng := rand.New(rand.NewSource(43))
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		h := rng.Intn(10) + 1
		a := make([]int, m)
		b := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, h))
		for j := 0; j < m; j++ {
			a[j] = rng.Intn(h + 1)
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(a[j]))
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			b[j] = rng.Intn(h + 1)
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(b[j]))
		}
		sb.WriteByte('\n')
		t := make([][]int, n)
		for i2 := 0; i2 < n; i2++ {
			t[i2] = make([]int, m)
			for j := 0; j < m; j++ {
				t[i2][j] = rng.Intn(2)
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(t[i2][j]))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		var out strings.Builder
		for i2 := 0; i2 < n; i2++ {
			for j := 0; j < m; j++ {
				val := 0
				if t[i2][j] == 1 {
					if a[j] < b[i2] {
						val = a[j]
					} else {
						val = b[i2]
					}
				}
				if j > 0 {
					out.WriteByte(' ')
				}
				out.WriteString(strconv.Itoa(val))
			}
			if i2 < n-1 {
				out.WriteByte('\n')
			}
		}
		tests = append(tests, Test{input: input, expected: out.String()})
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
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, tc.input, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
