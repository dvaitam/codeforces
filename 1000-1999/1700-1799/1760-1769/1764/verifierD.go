package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func modFact(n int, mod int64) int64 {
	res := int64(1)
	for i := 2; i <= n; i++ {
		res = res * int64(i) % mod
	}
	return res
}

func expectedD(n int, p int64) int64 {
	factN := modFact(n, p)
	factA := modFact(n/2, p)
	factB := modFact((n+1)/2, p)
	ans := (factN - 2*(factA*factB%p)) % p
	if ans < 0 {
		ans += p
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]int64, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(15) + 1
		p := int64(1000000007)
		fmt.Fprintf(&input, "%d %d\n", n, p)
		expected[i] = expectedD(n, p)
	}

	cmd := exec.Command(cand)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("candidate run error:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(outBytes))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("missing output for test", i+1)
			os.Exit(1)
		}
		var ans int64
		fmt.Sscan(scanner.Text(), &ans)
		if ans != expected[i] {
			fmt.Printf("wrong answer on test %d: expected %d got %d\n", i+1, expected[i], ans)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
