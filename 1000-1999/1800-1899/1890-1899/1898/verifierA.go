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

func expectedA(n, k int, s string) string {
	cntB := 0
	for i := 0; i < n; i++ {
		if s[i] == 'B' {
			cntB++
		}
	}
	if cntB == k {
		return "0"
	}
	if cntB < k {
		diff := k - cntB
		pref := 0
		for i := 0; i < n; i++ {
			if s[i] == 'A' {
				pref++
			}
			if pref == diff {
				return fmt.Sprintf("1\n%d B", i+1)
			}
		}
	} else {
		diff := cntB - k
		pref := 0
		for i := 0; i < n; i++ {
			if s[i] == 'B' {
				pref++
			}
			if pref == diff {
				return fmt.Sprintf("1\n%d A", i+1)
			}
		}
	}
	return ""
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(98) + 3
	k := rng.Intn(n + 1)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('A')
		} else {
			sb.WriteByte('B')
		}
	}
	s := sb.String()
	input := fmt.Sprintf("1\n%d %d\n%s\n", n, k, s)
	expected := expectedA(n, k, s)
	return input, expected
}

func runCase(bin, input, expect string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	m, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid first token: %v", err)
	}
	var got string
	if m == 0 {
		if len(fields) != 1 {
			return fmt.Errorf("expected single integer")
		}
		got = "0"
	} else if m == 1 {
		if len(fields) != 3 {
			return fmt.Errorf("expected three tokens for one operation")
		}
		i, err := strconv.Atoi(fields[1])
		if err != nil {
			return fmt.Errorf("bad index: %v", err)
		}
		c := strings.ToUpper(fields[2])
		if c != "A" && c != "B" {
			return fmt.Errorf("bad character %s", c)
		}
		got = fmt.Sprintf("1\n%d %s", i, c)
	} else {
		return fmt.Errorf("invalid m %d", m)
	}
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
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
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
