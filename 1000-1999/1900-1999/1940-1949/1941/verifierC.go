package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseC struct {
	s   string
	exp string
}

func solveC(s string) string {
	n := len(s)
	ans := 0
	for i := 0; i < n; {
		if i+4 < n && s[i:i+5] == "mapie" {
			ans++
			i += 5
		} else if i+2 < n && (s[i:i+3] == "pie" || s[i:i+3] == "map") {
			ans++
			i += 3
		} else {
			i++
		}
	}
	return fmt.Sprint(ans)
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCaseC {
	rng := rand.New(rand.NewSource(3))
	cases := make([]testCaseC, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := range cases {
		n := rng.Intn(20) + 1
		sb := strings.Builder{}
		for j := 0; j < n; j++ {
			sb.WriteRune(letters[rng.Intn(len(letters))])
		}
		s := sb.String()
		cases[i] = testCaseC{s: s, exp: solveC(s)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		input := fmt.Sprintf("1\n%d\n%s\n", len(tc.s), tc.s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
