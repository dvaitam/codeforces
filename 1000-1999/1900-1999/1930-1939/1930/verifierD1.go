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

func minOnes(p string) int {
	n := len(p)
	covered := make([]bool, n)
	ans := 0
	for i := 0; i < n; i++ {
		if p[i] == '1' && !covered[i] {
			pos := i + 1
			if pos >= n {
				pos = i
			}
			ans++
			for j := pos - 1; j <= pos+1; j++ {
				if j >= 0 && j < n {
					covered[j] = true
				}
			}
		}
	}
	return ans
}

func solveD1(n int, s string) int {
	total := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			total += minOnes(s[i:j])
		}
	}
	return total
}

func generateCaseD1(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	bytes := make([]byte, n)
	for i := range bytes {
		if rng.Intn(2) == 0 {
			bytes[i] = '0'
		} else {
			bytes[i] = '1'
		}
	}
	s := string(bytes)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
	input := sb.String()
	exp := fmt.Sprintf("%d\n", solveD1(n, s))
	return input, exp
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD1(rng)
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
