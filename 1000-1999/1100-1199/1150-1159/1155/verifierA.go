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

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// findPair returns whether there exists l<r so that reversing makes string smaller
func findPair(s string) (bool, int, int) {
	n := len(s)
	minChar := make([]byte, n)
	pos := make([]int, n)
	minChar[n-1] = s[n-1]
	pos[n-1] = n - 1
	for i := n - 2; i >= 0; i-- {
		if s[i+1] < minChar[i+1] {
			minChar[i] = s[i+1]
			pos[i] = i + 1
		} else {
			minChar[i] = minChar[i+1]
			pos[i] = pos[i+1]
		}
	}
	for i := 0; i < n-1; i++ {
		if minChar[i] < s[i] {
			return true, i + 1, pos[i] + 1
		}
	}
	return false, 0, 0
}

func checkCase(bin string, n int, s string) error {
	input := fmt.Sprintf("%d\n%s\n", n, s)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	has, _, _ := findPair(s)
	if !has {
		if len(fields) != 1 || strings.ToUpper(fields[0]) != "NO" {
			return fmt.Errorf("expected NO, got %s", out)
		}
		return nil
	}
	if len(fields) < 3 || strings.ToUpper(fields[0]) != "YES" {
		return fmt.Errorf("expected YES l r, got %s", out)
	}
	l, err1 := strconv.Atoi(fields[1])
	r, err2 := strconv.Atoi(fields[2])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("invalid indices in output: %s", out)
	}
	if l < 1 || l >= r || r > n {
		return fmt.Errorf("indices out of range: %s", out)
	}
	b := []byte(s)
	for i, j := l-1, r-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	if !(string(b) < s) {
		return fmt.Errorf("reversing [%d,%d] does not make string smaller", l, r)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{"ab", "aa", "abc", "cba", "abacaba"}
	for len(tests) < 100 {
		n := rng.Intn(8) + 2
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		tests = append(tests, sb.String())
	}
	for i, s := range tests {
		if err := checkCase(bin, len(s), s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
