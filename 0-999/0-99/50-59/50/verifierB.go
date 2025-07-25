package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

func expected(s string) int64 {
	var freq [256]int64
	for i := 0; i < len(s); i++ {
		freq[s[i]]++
	}
	var ans int64
	for _, v := range freq {
		ans += v * v
	}
	return ans
}

func randomString(n int, r *rand.Rand) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(42))
	const cases = 120
	for i := 0; i < cases; i++ {
		l := r.Intn(50) + 1
		s := randomString(l, r)
		input := s + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Fscan(&out, &got); err != nil {
			fmt.Printf("Test %d failed to parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		want := expected(s)
		if got != want {
			fmt.Printf("Test %d failed: expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
