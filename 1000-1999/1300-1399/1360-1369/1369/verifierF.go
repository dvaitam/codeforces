package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

var memo map[[2]int64]bool

func win(s, e int64) bool {
	if s > e {
		return false
	}
	key := [2]int64{s, e}
	if v, ok := memo[key]; ok {
		return v
	}
	var res bool
	if 2*s > e {
		res = (e-s)%2 == 1
	} else {
		res = !(win(s+1, e) && win(2*s, e))
	}
	memo[key] = res
	return res
}

func solveF(s, e int64) string {
	memo = make(map[[2]int64]bool)
	if win(s, e) {
		return "1 0"
	}
	return "0 1"
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(6)
	const tests = 100
	for t := 0; t < tests; t++ {
		s := int64(rand.Intn(20) + 1)
		e := s + int64(rand.Intn(20))
		input := fmt.Sprintf("1\n%d %d\n", s, e)
		expected := solveF(s, e) + "\n"
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", t+1, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
