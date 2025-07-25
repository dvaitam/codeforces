package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveC(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(reader, &n)
	xs := make([]int64, n)
	hs := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &xs[i], &hs[i])
	}
	if n == 1 {
		return "1"
	}
	ans := 1
	last := xs[0]
	for i := 1; i < n-1; i++ {
		if xs[i]-hs[i] > last {
			ans++
			last = xs[i]
		} else if xs[i]+hs[i] < xs[i+1] {
			ans++
			last = xs[i] + hs[i]
		} else {
			last = xs[i]
		}
	}
	ans++
	return strconv.Itoa(ans)
}

func genTests() []string {
	rand.Seed(3)
	tests := make([]string, 0, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		xs := make([]int64, n)
		hs := make([]int64, n)
		pos := int64(0)
		for i := 0; i < n; i++ {
			pos += int64(rand.Intn(5) + 1)
			xs[i] = pos
			hs[i] = int64(rand.Intn(5) + 1)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintln(n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintln(xs[i], hs[i]))
		}
		tests = append(tests, strings.TrimSpace(sb.String()))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierC <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		expected := solveC(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
