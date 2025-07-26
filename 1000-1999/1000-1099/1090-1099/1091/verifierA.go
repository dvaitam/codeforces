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

func expected(y, b, r int) string {
	m1 := 3*r - 3
	m2 := 3 * b
	m3 := 3*y + 3
	res := m1
	if m2 < res {
		res = m2
	}
	if m3 < res {
		res = m3
	}
	return fmt.Sprintf("%d", res)
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		y := rand.Intn(100) + 1
		b := rand.Intn(100) + 1
		r := rand.Intn(100) + 1
		input := fmt.Sprintf("%d %d %d\n", y, b, r)
		expect := expected(y, b, r)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:", input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("wrong answer on test %d\ninput: %sexpected: %s\n got: %s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
