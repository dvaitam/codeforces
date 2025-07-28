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

func genTests() []byte {
	rand.Seed(42)
	var sb strings.Builder
	t := 120
	fmt.Fprintf(&sb, "%d\n", t)
	// few fixed edge cases
	edge := [][2]int64{{1, 1}, {1, 2}, {1, 3}, {2, 2}, {2, 3}, {2, 4}}
	for _, e := range edge {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for i := len(edge); i < t; i++ {
		a := rand.Int63n(1e9) + 1
		b := rand.Int63n(1e9) + 1
		fmt.Fprintf(&sb, "%d %d\n", a, b)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	input := genTests()
	candCmd := exec.Command(os.Args[1])
	candOut, err := runCmd(candCmd, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate error: %v\n", err)
		os.Exit(1)
	}
	refCmd := exec.Command("go", "run", "1928A.go")
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
	_ = time.Now() // avoid import warning
}
