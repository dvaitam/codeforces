package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseB struct {
	n int
	a int
	b int
	c int
	s string
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func computeWin(tc testCaseB) (bool, string) {
	ans := make([]byte, tc.n)
	wins := 0
	a, b, c := tc.a, tc.b, tc.c
	for i, ch := range tc.s {
		switch ch {
		case 'R':
			if b > 0 {
				ans[i] = 'P'
				b--
				wins++
			}
		case 'P':
			if c > 0 {
				ans[i] = 'S'
				c--
				wins++
			}
		case 'S':
			if a > 0 {
				ans[i] = 'R'
				a--
				wins++
			}
		}
	}
	if wins < (tc.n+1)/2 {
		return false, ""
	}
	for i := 0; i < tc.n; i++ {
		if ans[i] == 0 {
			if a > 0 {
				ans[i] = 'R'
				a--
			} else if b > 0 {
				ans[i] = 'P'
				b--
			} else {
				ans[i] = 'S'
				c--
			}
		}
	}
	return true, string(ans)
}

func validateOutput(tc testCaseB, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" || len(lines) < 2 {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	first := strings.ToUpper(lines[0])
	canWin, _ := computeWin(tc)
	if first == "NO" {
		if canWin {
			return fmt.Errorf("should be YES")
		}
		if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
			return fmt.Errorf("extra output after NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("first line must be YES or NO")
	}
	if !canWin {
		return fmt.Errorf("should be NO")
	}
	if len(lines) < 2 {
		return fmt.Errorf("missing answer string")
	}
	ans := lines[1]
	if len(ans) != tc.n {
		return fmt.Errorf("answer length %d expected %d", len(ans), tc.n)
	}
	rc, pc, sc := 0, 0, 0
	wins := 0
	for i := 0; i < tc.n; i++ {
		ch := ans[i]
		switch ch {
		case 'R':
			rc++
			if tc.s[i] == 'S' {
				wins++
			}
		case 'P':
			pc++
			if tc.s[i] == 'R' {
				wins++
			}
		case 'S':
			sc++
			if tc.s[i] == 'P' {
				wins++
			}
		default:
			return fmt.Errorf("invalid char %c", ch)
		}
	}
	if rc != tc.a || pc != tc.b || sc != tc.c {
		return fmt.Errorf("counts mismatch")
	}
	if wins < (tc.n+1)/2 {
		return fmt.Errorf("not enough wins")
	}
	if len(lines) > 2 && strings.TrimSpace(lines[2]) != "" {
		return fmt.Errorf("extra output lines")
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, testCaseB) {
	n := rng.Intn(100) + 1
	a := rng.Intn(n + 1)
	b := rng.Intn(n - a + 1)
	c := n - a - b
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	moves := []byte{'R', 'P', 'S'}
	sbytes := make([]byte, n)
	for i := 0; i < n; i++ {
		sbytes[i] = moves[rng.Intn(3)]
	}
	s := string(sbytes)
	sb.WriteString(s)
	sb.WriteByte('\n')
	return sb.String(), testCaseB{n: n, a: a, b: b, c: c, s: s}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validateOutput(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
