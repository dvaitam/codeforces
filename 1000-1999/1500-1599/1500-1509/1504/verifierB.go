package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
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

func canTransform(a, b string) bool {
	n := len(a)
	balanced := make([]bool, n)
	diff := 0
	for i := 0; i < n; i++ {
		if a[i] == '1' {
			diff++
		} else {
			diff--
		}
		if diff == 0 {
			balanced[i] = true
		}
	}
	flip := false
	for i := n - 1; i >= 0; i-- {
		ch := a[i]
		if flip {
			if ch == '1' {
				ch = '0'
			} else {
				ch = '1'
			}
		}
		if ch == b[i] {
			continue
		}
		if !balanced[i] {
			return false
		}
		flip = !flip
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	type pair struct{ a, b string }
	cases := []pair{
		{"0", "0"},
		{"0", "1"},
		{"01", "10"},
		{"10", "10"},
		{"1100", "0011"},
		{"01", "01"},
		{"11", "00"},
		{"1010", "0101"},
		{"000", "111"},
	}
	rng := rand.New(rand.NewSource(42))
	for len(cases) < 100 {
		n := rng.Intn(20) + 1
		var sa, sb strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sa.WriteByte('0')
			} else {
				sa.WriteByte('1')
			}
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		cases = append(cases, pair{sa.String(), sb.String()})
	}
	for idx, p := range cases {
		wantYes := canTransform(p.a, p.b)
		expected := "NO"
		if wantYes {
			expected = "YES"
		}
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", len(p.a), p.a, p.b)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: a=%q b=%q expected %q got %q\n", idx+1, p.a, p.b, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
