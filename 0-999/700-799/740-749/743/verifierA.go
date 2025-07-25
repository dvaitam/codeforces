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

func runCandidate(bin, input string) (string, error) {
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

func solveCaseA(n, a, b int, s string) int {
	if s[a-1] == s[b-1] {
		return 0
	}
	pos := [2][]int{}
	for i, ch := range s {
		c := int(ch - '0')
		if c == 0 || c == 1 {
			pos[c] = append(pos[c], i+1)
		}
	}
	src := pos[int(s[a-1]-'0')]
	dst := pos[int(s[b-1]-'0')]
	i, j := 0, 0
	ans := n
	for i < len(src) && j < len(dst) {
		diff := src[i] - dst[j]
		if diff < 0 {
			diff = -diff
		}
		if diff < ans {
			ans = diff
		}
		if src[i] < dst[j] {
			i++
		} else {
			j++
		}
	}
	return ans
}

func generateCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	a := rng.Intn(n) + 1
	b := rng.Intn(n) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	s := sb.String()
	input := fmt.Sprintf("%d %d %d\n%s\n", n, a, b, s)
	expect := fmt.Sprintf("%d", solveCaseA(n, a, b, s))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseA(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
