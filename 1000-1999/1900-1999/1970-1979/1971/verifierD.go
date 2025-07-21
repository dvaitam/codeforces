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

func expectedPieces(s string) int {
	cnt := 1
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '1' && s[i+1] == '0' {
			cnt++
		}
	}
	return cnt
}

type caseD struct{ s string }

func genCase(rng *rand.Rand) caseD {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return caseD{string(b)}
}

func runCase(bin string, tc caseD) error {
	input := fmt.Sprintf("1\n%s\n", tc.s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expectedPieces(tc.s)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s\n", i+1, err, tc.s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
