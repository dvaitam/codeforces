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

func solve(n int) string {
	if n == 3 || n == 5 {
		return "NO"
	}
	ara := make([]int, n)
	if n%2 == 0 {
		p := []int{-1, 1, -1, -2, 2, 1}
		for i := 0; i < n; i++ {
			ara[i] = p[i%6]
		}
	} else {
		p := []int{2, -2, -3, 3, 1, -1}
		s := []int{3, 3, -2, -1, 1, -1, 2}
		for i := 0; i < n; i++ {
			if i < len(s) {
				ara[i] = s[i]
			} else {
				ara[i] = p[i%6]
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i, v := range ara {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solve(n)
	if n == 3 || n == 5 {
		if strings.ToUpper(got) != "NO" {
			return fmt.Errorf("expected NO got %s", got)
		}
		return nil
	}
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		var n int
		if i == 0 {
			n = 4
		} else {
			n = rng.Intn(50) + 2
			if n == 3 || n == 5 {
				n = 6
			}
		}
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
