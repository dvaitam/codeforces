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

func expectedC(s string) string {
	freq := make([]int, 26)
	for i := 0; i < len(s); i++ {
		freq[s[i]-'a']++
	}
	minIdx := 0
	for minIdx < 26 && freq[minIdx] == 0 {
		minIdx++
	}
	stack := make([]byte, 0, len(s))
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		stack = append(stack, c)
		freq[c-'a']--
		for minIdx < 26 && freq[minIdx] == 0 {
			minIdx++
		}
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if minIdx >= 26 || top <= byte('a'+minIdx) {
				result = append(result, top)
				stack = stack[:len(stack)-1]
			} else {
				break
			}
		}
	}
	for i := len(stack) - 1; i >= 0; i-- {
		result = append(result, stack[i])
	}
	return string(result)
}

func generateCase(rng *rand.Rand) (string, string) {
	l := rng.Intn(20) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	s := string(b)
	var sb strings.Builder
	sb.WriteString(s)
	sb.WriteByte('\n')
	return sb.String(), expectedC(s)
}

func runCase(bin string, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	res := strings.TrimSpace(out.String())
	if res != expect {
		return fmt.Errorf("expected %q got %q", expect, res)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
