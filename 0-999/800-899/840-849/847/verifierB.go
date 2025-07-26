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
	sol := "847B.go"
	exe := "./_verifier_solB"
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

func normalize(s string) string { return strings.TrimSpace(s) }

func genTest() string {
	n := rand.Intn(10) + 1
	arr := rand.Perm(1000)[:n]
	for i := range arr {
		arr[i]++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		input := genTest()
		exp, err := runBinary(solExe, input)
		if err != nil {
			fmt.Printf("reference solution failed on test %d: %v\n", t, err)
			return
		}
		out, err := runBinary(bin, input)
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
