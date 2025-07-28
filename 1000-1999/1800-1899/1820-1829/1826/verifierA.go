package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(arr []int) int {
	n := len(arr)
	for x := 0; x <= n; x++ {
		cnt := 0
		for _, v := range arr {
			if v > x {
				cnt++
			}
		}
		if cnt == x {
			return x
		}
	}
	return -1
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierA <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(0)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(100) + 1
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rand.Intn(n + 1)
		}
		input := fmt.Sprintf("1\n%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		expected := fmt.Sprint(solve(arr))
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s\n", t, err, output)
			os.Exit(1)
		}
		if output != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", t, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
