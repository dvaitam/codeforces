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

func canUniform(seg string) bool {
	has0 := strings.Contains(seg, "0")
	has1 := strings.Contains(seg, "1")
	return !(has0 && has1)
}

func solveCase(input string) string {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return ""
	}
	n, _ := strconv.Atoi(fields[0])
	s := fields[1]
	res := make([]string, n)
	for x := 1; x <= n; x++ {
		cnt := 0
		i := 0
		for i+x <= n {
			seg := s[i : i+x]
			if canUniform(seg) {
				cnt++
				i += x
			} else {
				i++
			}
		}
		res[x-1] = fmt.Sprint(cnt)
	}
	return strings.Join(res, " ")
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	bytesS := make([]byte, n)
	chars := []byte{'0', '1', '?'}
	for i := 0; i < n; i++ {
		bytesS[i] = chars[rng.Intn(3)]
	}
	return fmt.Sprintf("%d\n%s\n", n, string(bytesS))
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected = strings.TrimSpace(expected)
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := solveCase(in)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
