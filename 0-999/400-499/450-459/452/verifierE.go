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

const mod = 1000000007

func solve(s1, s2, s3 string) string {
	minL := len(s1)
	if len(s2) < minL {
		minL = len(s2)
	}
	if len(s3) < minL {
		minL = len(s3)
	}
	res := make([]int64, minL+1)
	for l := 1; l <= minL; l++ {
		var cnt int64
		for i := 0; i+l <= len(s1); i++ {
			sub := s1[i : i+l]
			for j := 0; j+l <= len(s2); j++ {
				if s2[j:j+l] != sub {
					continue
				}
				for k := 0; k+l <= len(s3); k++ {
					if s3[k:k+l] == sub {
						cnt++
					}
				}
			}
		}
		res[l] = cnt % mod
	}
	var sb strings.Builder
	for l := 1; l <= minL; l++ {
		if l > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(res[l], 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n1 := rng.Intn(5) + 1
	n2 := rng.Intn(5) + 1
	n3 := rng.Intn(5) + 1
	s1 := randString(rng, n1)
	s2 := randString(rng, n2)
	s3 := randString(rng, n3)
	input := fmt.Sprintf("%s\n%s\n%s\n", s1, s2, s3)
	expected := solve(s1, s2, s3)
	return input, expected
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
