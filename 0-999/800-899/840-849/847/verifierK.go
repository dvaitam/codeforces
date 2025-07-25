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
	sol := "847K.go"
	exe := "./_verifier_solK"
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
	n := rand.Intn(5) + 1
	a := rand.Intn(5) + 1
	b := rand.Intn(5) + 1
	k := rand.Intn(3) + 1
	f := rand.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", n, a, b, k, f))
	stops := []string{"A", "B", "C", "D", "E", "F"}
	for i := 0; i < n; i++ {
		s := stops[rand.Intn(len(stops))]
		t := stops[rand.Intn(len(stops))]
		sb.WriteString(fmt.Sprintf("%s %s\n", s, t))
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
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
