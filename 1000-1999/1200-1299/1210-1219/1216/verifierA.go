package main

import (
	"bytes"
	"fmt"
	"io"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func minimalOps(s string) int {
	cnt := 0
	for i := 0; i < len(s); i += 2 {
		if s[i] == s[i+1] {
			cnt++
		}
	}
	return cnt
}

func checkCase(orig string, outStr string) error {
	var cnt int
	var modified string
	reader := strings.NewReader(outStr)
	if _, err := fmt.Fscan(reader, &cnt, &modified); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err != io.EOF {
		return fmt.Errorf("extra output detected")
	}
	if len(modified) != len(orig) {
		return fmt.Errorf("expected string length %d got %d", len(orig), len(modified))
	}
	// validate characters and pairs
	for i := 0; i < len(modified); i++ {
		if modified[i] != 'a' && modified[i] != 'b' {
			return fmt.Errorf("invalid character %c", modified[i])
		}
	}
	for i := 0; i < len(modified); i += 2 {
		if modified[i] == modified[i+1] {
			return fmt.Errorf("pair %d still equal", i/2)
		}
	}
	diff := 0
	for i := 0; i < len(modified); i++ {
		if modified[i] != orig[i] {
			diff++
		}
	}
	if diff != cnt {
		return fmt.Errorf("reported %d changes but changed %d", cnt, diff)
	}
	need := minimalOps(orig)
	if cnt != need {
		return fmt.Errorf("expected %d operations got %d", need, cnt)
	}
	return nil
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50)*2 + 2 // even length between 2 and 100
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 'a'
		} else {
			b[i] = 'b'
		}
	}
	s := string(b)
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return input, s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, s := genCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkCase(s, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
