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

func expected(s string) string {
	var cnt [26]int
	var rep []byte
	for i := 0; i < len(s); i++ {
		ch := s[i]
		idx := ch - 'a'
		cnt[idx]++
		if cnt[idx] == 2 {
			rep = append(rep, ch)
		}
	}
	ans := make([]byte, 0, len(s))
	ans = append(ans, rep...)
	ans = append(ans, rep...)
	for ch := byte('a'); ch <= 'z'; ch++ {
		if cnt[ch-'a'] == 1 {
			ans = append(ans, ch)
		}
	}
	return string(ans)
}

func runCase(bin, s string) error {
	input := fmt.Sprintf("1\n%s\n", s)
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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expected(s)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) string {
	length := rng.Intn(52) + 1
	counts := [26]int{}
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		for {
			c := byte(rng.Intn(26))
			if counts[c] < 2 {
				counts[c]++
				b[i] = 'a' + c
				break
			}
		}
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"a", "aa", "ab", "abac", "zz", "abcdef", "abacaba"}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for idx, s := range cases {
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
