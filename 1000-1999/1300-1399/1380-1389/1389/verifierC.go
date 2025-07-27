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

func solveCase(s string) int {
	n := len(s)
	best := 0
	freq := make([]int, 10)
	for i := 0; i < n; i++ {
		d := int(s[i] - '0')
		freq[d]++
	}
	for i := 0; i < 10; i++ {
		if freq[i] > best {
			best = freq[i]
		}
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i == j {
				continue
			}
			expect := i
			length := 0
			for k := 0; k < n; k++ {
				d := int(s[k] - '0')
				if d == expect {
					length++
					if expect == i {
						expect = j
					} else {
						expect = i
					}
				}
			}
			if length%2 == 1 {
				length--
			}
			if length > best {
				best = length
			}
		}
	}
	return n - best
}

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
	return strings.TrimSpace(out.String()), nil
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + rng.Intn(10))
	}
	return string(b)
}

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(12) + 1
	s := randString(rng, n)
	input := fmt.Sprintf("1\n%s\n", s)
	exp := fmt.Sprintf("%d", solveCase(s))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, exp := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
