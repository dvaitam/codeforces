package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "1610H.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTree(n int) []int {
	parents := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parents[i-2] = rand.Intn(i-1) + 1
	}
	return parents
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(5) + 2
		m := rand.Intn(5) + 1
		parents := genTree(n)
		input := fmt.Sprintf("%d %d\n", n, m)
		for _, p := range parents {
			input += fmt.Sprintf("%d ", p)
		}
		input += "\n"
		for i := 0; i < m; i++ {
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			input += fmt.Sprintf("%d %d\n", x, y)
		}
		exp, err := runRef(input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d exec failed: %v\n", t, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed expected %s got %s\n", t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
