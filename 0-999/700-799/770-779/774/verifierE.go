package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(s string, m int64) int64 {
	n := len(s)
	t := s + s
	pref := make([]int64, 2*n+1)
	for i := 0; i < 2*n; i++ {
		d := int64(t[i] - '0')
		pref[i+1] = (pref[i]*10 + d) % m
	}
	pow10n := int64(1)
	for i := 0; i < n; i++ {
		pow10n = (pow10n * 10) % m
	}
	minRem := int64(-1)
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			continue
		}
		val := (pref[i+n] - (pref[i]*pow10n)%m + m) % m
		if minRem == -1 || val < minRem {
			minRem = val
		}
	}
	return minRem
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(6) + 1
	if rng.Float64() < 0.1 {
		n = rng.Intn(12) + 1
	}
	sb := make([]byte, n)
	for i := 0; i < n; i++ {
		sb[i] = byte('0' + rng.Intn(10))
	}
	s := string(sb)
	m := int64(rng.Intn(1000) + 1)
	input := fmt.Sprintf("%s\n%d\n", s, m)
	return input, expected(s, m)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
