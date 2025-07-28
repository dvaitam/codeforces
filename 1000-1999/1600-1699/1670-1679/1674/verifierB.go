package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(s string) string {
	i := int(s[0] - 'a' + 1)
	j := int(s[1] - 'a' + 1)
	ans := (i-1)*25 + j
	if j > i {
		ans--
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(43))
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for t := 0; t < 100; t++ {
		var a, b byte
		for {
			a = letters[rng.Intn(26)]
			b = letters[rng.Intn(26)]
			if a != b {
				break
			}
		}
		s := string([]byte{a, b})
		input := fmt.Sprintf("1\n%s\n", s)
		exp := expected(s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\n", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
