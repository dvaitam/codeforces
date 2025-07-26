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

func solveA(n, t, k, d int) string {
	batches := (n + k - 1) / k
	oneTime := batches * t
	for current := 1; current < oneTime; current++ {
		cakes := (current / t) * k
		if current > d {
			cakes += ((current - d) / t) * k
		}
		if cakes >= n {
			return "YES\n"
		}
	}
	return "NO\n"
}

func genA() (string, string) {
	n := rand.Intn(50) + 1
	t := rand.Intn(10) + 1
	k := rand.Intn(10) + 1
	d := rand.Intn(20) + 1
	input := fmt.Sprintf("%d %d %d %d\n", n, t, k, d)
	out := solveA(n, t, k, d)
	return input, out
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genA()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
