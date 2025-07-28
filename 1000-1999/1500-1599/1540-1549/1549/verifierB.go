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

type testB struct {
	n     int
	enemy string
	greg  string
}

func runProg(bin, input string) (string, error) {
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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, enemy, greg string) int {
	e := []byte(enemy)
	g := []byte(greg)
	ans := 0
	for i := 0; i < n; i++ {
		if g[i] == '1' {
			if e[i] == '0' {
				ans++
				e[i] = '2'
			} else if i > 0 && e[i-1] == '1' {
				ans++
				e[i-1] = '2'
			} else if i+1 < n && e[i+1] == '1' {
				ans++
				e[i+1] = '2'
			}
		}
	}
	return ans
}

func genTests() []testB {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []testB{
		{2, "00", "00"},
		{2, "11", "11"},
		{3, "101", "010"},
		{4, "1111", "1111"},
		{5, "00000", "11111"},
	}
	for len(tests) < 100 {
		n := rng.Intn(100) + 2
		var sb1, sb2 strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 1 {
				sb1.WriteByte('1')
			} else {
				sb1.WriteByte('0')
			}
			if rng.Intn(2) == 1 {
				sb2.WriteByte('1')
			} else {
				sb2.WriteByte('0')
			}
		}
		tests = append(tests, testB{n, sb1.String(), sb2.String()})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", tc.n, tc.enemy, tc.greg)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		var got int
		if n, _ := fmt.Sscan(out, &got); n != 1 {
			fmt.Fprintf(os.Stderr, "case %d: expected single integer, got %q\n", i+1, out)
			os.Exit(1)
		}
		exp := expected(tc.n, tc.enemy, tc.greg)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
