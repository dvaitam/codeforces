package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	if errBuf.Len() > 0 {
		return "", fmt.Errorf(errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n uint64) string {
	if bits.OnesCount64(n)%2 == 0 {
		p1 := n & -n
		p2 := n - p1
		if p2 == 0 {
			p1, p2 = 1, n-1
		}
		return fmt.Sprintf("first\n%d %d", p1, p2)
	}
	return "second"
}

func genTests() []uint64 {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]uint64, 100)
	for i := range cases {
		cases[i] = rng.Uint64()%1_000_000 + 1
	}
	cases = append(cases, 1, 2, 3, 4)
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTests()
	for idx, n := range cases {
		input := fmt.Sprintf("1\n%d\n", n)
		expect := expected(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed: n=%d expected %q got %q\n", idx+1, n, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
