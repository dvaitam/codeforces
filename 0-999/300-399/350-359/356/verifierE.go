package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func isGray(s []byte) bool {
	n := len(s)
	if n%2 == 0 {
		return false
	}
	if n == 1 {
		return true
	}
	mid := n / 2
	c := s[mid]
	for i, ch := range s {
		if i != mid && ch == c {
			return false
		}
	}
	left := s[:mid]
	right := s[mid+1:]
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return isGray(left)
}

func beauty(s []byte) int64 {
	n := len(s)
	var sum int64
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			sub := s[i : j+1]
			if isGray(sub) {
				l := int64(j - i + 1)
				sum += l * l
			}
		}
	}
	return sum
}

func solveE(in string) string {
	t := strings.TrimSpace(in)
	s := []byte(t)
	n := len(s)
	best := beauty(s)
	for i := 0; i < n; i++ {
		orig := s[i]
		for c := byte('a'); c <= byte('z'); c++ {
			if c == orig {
				continue
			}
			s[i] = c
			v := beauty(s)
			if v > best {
				best = v
			}
			s[i] = orig
		}
	}
	return fmt.Sprint(best)
}

func genTest(r *rand.Rand) string {
	n := r.Intn(10) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + r.Intn(26)))
	}
	return sb.String() + "\n"
}

func runBinary(path, in string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierE <path-to-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(5))
	const tests = 100
	for i := 0; i < tests; i++ {
		in := genTest(r)
		expect := solveE(in)
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
