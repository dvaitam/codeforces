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

func solveD(n, x int) string {
	var sb strings.Builder
	if n == 5 {
		sb.WriteString(">...v\n")
		sb.WriteString("v.<..\n")
		sb.WriteString("..^..\n")
		sb.WriteString(">....\n")
		sb.WriteString("..^.<\n")
		sb.WriteString("1 1")
		return sb.String()
	}
	if n == 3 {
		sb.WriteString(">vv\n")
		sb.WriteString("^<.\n")
		sb.WriteString("^.<\n")
		sb.WriteString("1 3")
		return sb.String()
	}
	mp := make([][]rune, n)
	for i := 0; i < n; i++ {
		mp[i] = make([]rune, n)
		for j := 0; j < n; j++ {
			mp[i][j] = '.'
		}
	}
	for i := 0; i < n; i++ {
		mp[i][0] = '^'
	}
	mp[0][0] = '>'
	for i := 0; i < n; i += 2 {
		for j := 1; j < n-1; j++ {
			if j < n/2 || j%2 == 1 {
				mp[i][j] = '>'
			}
		}
		mp[i][n-1] = 'v'
	}
	for i := 1; i < n; i += 2 {
		for j := n - 1; j > 0; j-- {
			if (n-j) < n/2 || j%2 == 1 {
				mp[i][j] = '<'
			}
		}
		mp[i][1] = 'v'
	}
	mp[n-1][1] = '<'
	for i := 0; i < n; i++ {
		sb.WriteString(string(mp[i]))
		sb.WriteByte('\n')
	}
	sb.WriteString("1 1")
	return sb.String()
}

func generateCaseD(rng *rand.Rand) (string, string) {
	cases := [][2]int{{5, 5}, {3, 2}, {100, 105}}
	pick := cases[rng.Intn(len(cases))]
	n, x := pick[0], pick[1]
	in := fmt.Sprintf("%d %d\n", n, x)
	out := solveD(n, x)
	return in, out
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
