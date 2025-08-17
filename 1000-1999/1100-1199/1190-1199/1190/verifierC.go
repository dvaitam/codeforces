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

func run(bin, input string) (string, error) {
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

func isUniform(b []byte) bool {
	for i := 1; i < len(b); i++ {
		if b[i] != b[0] {
			return false
		}
	}
	return true
}

func moves(s string, k int) []string {
	b := []byte(s)
	res := make([]string, 0, (len(b)-k+1)*2)
	for l := 0; l+k <= len(b); l++ {
		for _, ch := range []byte{'0', '1'} {
			t := append([]byte(nil), b...)
			for i := 0; i < k; i++ {
				t[l+i] = ch
			}
			res = append(res, string(t))
		}
	}
	return res
}

func solve(n, k int, s string) string {
	for _, t := range moves(s, k) {
		if isUniform([]byte(t)) {
			return "tokitsukaze"
		}
	}
	for _, t := range moves(s, k) {
		secWin := false
		for _, u := range moves(t, k) {
			if isUniform([]byte(u)) {
				secWin = true
				break
			}
		}
		if !secWin {
			return "once again"
		}
	}
	return "quailty"
}

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	k := rng.Intn(n) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	s := string(b)
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	return input, solve(n, k, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
