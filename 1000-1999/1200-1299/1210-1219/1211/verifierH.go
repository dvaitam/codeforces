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

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genInput() string {
	n := rand.Intn(5) + 2
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]

	ref := "refH_bin"
	if err := exec.Command("go", "build", "-o", ref, "1211H.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer os.Remove(ref)

	for t := 0; t < 100; t++ {
		input := genInput()
		expected, _ := runProgram(ref, []byte(input))
		out, err := runProgram(bin, []byte(input))
		if err != nil || strings.TrimSpace(out) != strings.TrimSpace(expected) {
			fmt.Println("Test", t+1, "failed")
			fmt.Println("Input:\n", input)
			fmt.Println("Expected:\n", expected)
			fmt.Println("Output:\n", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("All tests passed")
}
