package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseE struct {
	input  string
	expect string
}

const modE = 1013

func solveE(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var f int
	if _, err := fmt.Fscan(reader, &f); err != nil {
		return ""
	}
	if f == 0 {
		return "0"
	}
	if f == 1 {
		return "1"
	}
	pos := make([]int, modE)
	for i := range pos {
		pos[i] = -1
	}
	pos[0] = 0
	pos[1] = 1
	prev, curr := 0, 1
	for i := 2; ; i++ {
		next := (prev + curr) % modE
		if pos[next] == -1 {
			pos[next] = i
		}
		if next == f {
			return fmt.Sprintf("%d", pos[next])
		}
		prev, curr = curr, next
		if prev == 0 && curr == 1 {
			break
		}
	}
	return "-1"
}

func genTests() []testCaseE {
	rand.Seed(42)
	values := []int{0, 1, 2, 3, 5, 8, 13, 100, 512, 1012}
	tests := make([]testCaseE, 0, 1013)
	for _, v := range values {
		input := fmt.Sprintf("%d\n", v)
		tests = append(tests, testCaseE{
			input:  input,
			expect: solveE(input),
		})
	}
	for i := 0; i < 1000; i++ {
		f := rand.Intn(modE)
		input := fmt.Sprintf("%d\n", f)
		tests = append(tests, testCaseE{
			input:  input,
			expect: solveE(input),
		})
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
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, t.input)
			os.Exit(1)
		}
		if got != t.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", i+1, t.input, t.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
