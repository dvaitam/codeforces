package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func genTests() []byte {
	rand.Seed(46)
	t := 110
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rand.Intn(30) + 1
		x := rand.Intn(60)
		y := rand.Intn(60) + 1
		s := rand.Intn(100)
		fmt.Fprintf(&sb, "%d %d %d %d\n", n, x, y, s)
	}
	return []byte(sb.String())
}

func runCmd(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	input := genTests()
	candCmd := exec.Command(os.Args[1])
	candOut, err := runCmd(candCmd, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate error: %v\n", err)
		os.Exit(1)
	}
	refCmd := exec.Command("go", "run", "1928E.go")
	refOut, err := runCmd(refCmd, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference error: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
		fmt.Println("WA")
		fmt.Println("input:\n" + string(input))
		fmt.Println("expected:\n" + refOut)
		fmt.Println("got:\n" + candOut)
		os.Exit(1)
	}
	fmt.Println("OK")
}
