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

func runCandidate(bin, input string) (string, error) {
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

func feasible(n, m int) bool {
	if n > m+1 || n*2+2 < m {
		return false
	}
	return true
}

func runCase(bin string, n, m int) error {
	input := fmt.Sprintf("%d %d\n", n, m)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if !feasible(n, m) {
		if strings.TrimSpace(out) != "-1" {
			return fmt.Errorf("expected -1 got %s", out)
		}
		return nil
	}
	s := strings.TrimSpace(out)
	if len(s) != n+m {
		return fmt.Errorf("wrong length %d", len(s))
	}
	zeros := 0
	ones := 0
	consecOnes := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '0' {
			zeros++
			consecOnes = 0
			if i > 0 && s[i-1] == '0' {
				return fmt.Errorf("has 00 substring")
			}
		} else if c == '1' {
			ones++
			consecOnes++
			if consecOnes > 2 {
				return fmt.Errorf("more than two consecutive ones")
			}
		} else {
			return fmt.Errorf("invalid char %c", c)
		}
	}
	if zeros != n || ones != m {
		return fmt.Errorf("expected %d zeros %d ones got %d zeros %d ones", n, m, zeros, ones)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10)
		m := rng.Intn(10)
		if err := runCase(bin, n, m); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
