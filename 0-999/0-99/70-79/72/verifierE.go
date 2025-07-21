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
	bestCount := 0
	bestLen := 0
	bestS := ""
	n := len(s)
	for i := 0; i < n; i++ {
		for l := 1; i+l <= n; l++ {
			sub := s[i : i+l]
			cnt := 0
			for j := 0; j+l <= n; j++ {
				if s[j:j+l] == sub {
					cnt++
				}
			}
			if cnt > bestCount || (cnt == bestCount && (l > bestLen || (l == bestLen && sub > bestS))) {
				bestCount = cnt
				bestLen = l
				bestS = sub
			}
		}
	}
	return bestS
}

func randomString(rng *rand.Rand) string {
	l := rng.Intn(8) + 2
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	s := randomString(rng)
	return fmt.Sprintf("%s\n", s), expected(s)
}

func runCase(bin, input, expected string) error {
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected '%s' got '%s'", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
