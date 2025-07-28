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

func solveCase(n, m uint64) ([]uint64, bool) {
	if m >= n || m == 0 {
		return nil, false
	}
	if n&(n-1) == 0 {
		return nil, false
	}
	pow := uint64(1) << (bits.Len64(n) - 1)
	r := n - pow
	if m >= pow {
		if n^m >= n {
			return nil, false
		}
		return []uint64{n, m}, true
	}
	h := uint64(1) << (bits.Len64(m) - 1)
	if r&h != 0 {
		return []uint64{n, m}, true
	}
	if r > h {
		return []uint64{n, pow + h, m}, true
	}
	return nil, false
}

func genTests() [][2]uint64 {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([][2]uint64, 100)
	for i := range tests {
		n := rng.Uint64()%1_000_000 + 2
		// ensure not power of two, else expectation may be false
		for n&(n-1) == 0 {
			n = rng.Uint64()%1_000_000 + 2
		}
		m := rng.Uint64()%uint64(n-1) + 1
		tests[i] = [2]uint64{n, m}
	}
	// edge case impossible
	tests = append(tests, [2]uint64{7, 3})
	tests = append(tests, [2]uint64{4, 2})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTests()
	for idx, c := range cases {
		input := fmt.Sprintf("1\n%d %d\n", c[0], c[1])
		seq, ok := solveCase(c[0], c[1])
		var expect string
		if !ok {
			expect = "-1"
		} else {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("%d\n", len(seq)-1))
			for i, v := range seq {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", v))
			}
			expect = strings.TrimSpace(sb.String())
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed: n=%d m=%d expected %q got %q\n", idx+1, c[0], c[1], expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
