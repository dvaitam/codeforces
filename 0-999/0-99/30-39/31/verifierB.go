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

type testCaseB struct {
	s     string
	valid bool
}

func generateValidEmail(rng *rand.Rand) string {
	m := rng.Intn(5) + 1
	parts := make([]string, m)
	for i := 0; i < m; i++ {
		aLen := rng.Intn(3) + 1
		bLen := rng.Intn(3) + 1
		a := randomLetters(rng, aLen)
		b := randomLetters(rng, bLen)
		parts[i] = a + "@" + b
	}
	return strings.Join(parts, "")
}

func generateInvalidString(rng *rand.Rand) string {
	// random letters and '@'
	n := rng.Intn(20) + 1
	var sb strings.Builder
	atUsed := false
	for i := 0; i < n; i++ {
		if rng.Intn(5) == 0 {
			sb.WriteByte('@')
			atUsed = true
		} else {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
	}
	if !atUsed {
		sb.WriteByte('@')
	}
	return sb.String()
}

func randomLetters(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCaseB(rng *rand.Rand) testCaseB {
	if rng.Intn(2) == 0 {
		return testCaseB{s: generateValidEmail(rng), valid: true}
	}
	return testCaseB{s: generateInvalidString(rng), valid: false}
}

func solveB(s string) (string, bool) {
	n := len(s)
	var pos []int
	for i, c := range s {
		if c == '@' {
			pos = append(pos, i)
		}
	}
	m := len(pos)
	if m == 0 {
		return "No solution", false
	}
	if pos[0] < 1 || pos[m-1] > n-2 {
		return "No solution", false
	}
	for i := 0; i < m-1; i++ {
		if pos[i+1]-pos[i] < 3 {
			return "No solution", false
		}
	}
	var parts []string
	start := 0
	for i := 0; i < m-1; i++ {
		end := pos[i+1] - 2
		if end < start {
			return "No solution", false
		}
		parts = append(parts, s[start:end+1])
		start = end + 1
	}
	parts = append(parts, s[start:])
	for _, p := range parts {
		idx := strings.Index(p, "@")
		if idx <= 0 || idx >= len(p)-1 || strings.Count(p, "@") != 1 {
			return "No solution", false
		}
	}
	return strings.Join(parts, ","), true
}

func checkOutputB(input, out string) error {
	_, ok := solveB(input)
	out = strings.TrimSpace(out)
	if !ok {
		if out != "No solution" {
			return fmt.Errorf("expected 'No solution', got %q", out)
		}
		return nil
	}
	// verify candidate solution is valid
	if strings.ReplaceAll(out, ",", "") != input {
		return fmt.Errorf("output does not reconstruct original string")
	}
	segments := strings.Split(out, ",")
	for _, p := range segments {
		idx := strings.Index(p, "@")
		if idx <= 0 || idx >= len(p)-1 || strings.Count(p, "@") != 1 {
			return fmt.Errorf("segment %q is invalid", p)
		}
	}
	return nil
}

func runCaseB(bin string, tc testCaseB) error {
	input := tc.s + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return checkOutputB(tc.s, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s\n", i+1, err, tc.s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
