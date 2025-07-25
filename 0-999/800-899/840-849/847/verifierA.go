package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildSolution() (string, error) {
	sol := "847A.go"
	exe := "./_verifier_solA"
	cmd := exec.Command("go", "build", "-o", exe, sol)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

type testCase struct {
	n     int
	l     []int
	r     []int
	input string
}

func genTest() testCase {
	n := rand.Intn(10) + 1
	nodes := rand.Perm(n)
	l := make([]int, n+1)
	r := make([]int, n+1)
	idx := 0
	for idx < n {
		length := rand.Intn(n-idx) + 1
		for j := 0; j < length; j++ {
			node := nodes[idx+j] + 1
			if j == 0 {
				l[node] = 0
			} else {
				l[node] = nodes[idx+j-1] + 1
			}
			if j == length-1 {
				r[node] = 0
			} else {
				r[node] = nodes[idx+j+1] + 1
			}
		}
		idx += length
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", l[i], r[i]))
	}
	return testCase{n: n, l: l, r: r, input: sb.String()}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())

	solExe, err := buildSolution()
	if err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer os.Remove(solExe)

	for t := 1; t <= 100; t++ {
		tc := genTest()
		exp, err := runBinary(solExe, tc.input)
		if err != nil {
			fmt.Printf("reference solution failed on test %d: %v\n", t, err)
			return
		}
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", t, err)
			return
		}
		if normalize(exp) != normalize(out) {
			fmt.Printf("wrong answer on test %d\nexpected:\n%s\nactual:\n%s\n", t, exp, out)
			return
		}
	}
	fmt.Println("tests passed")
}
