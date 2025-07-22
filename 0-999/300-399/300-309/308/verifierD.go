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

func expectedD(n, m int) string {
	N := n + 1
	low := m + 1
	high := n - m
	if low > high {
		return "0"
	}
	var sumP, sumQ, sumR int64
	for a := low; a <= high; a++ {
		for b := low; b <= high; b++ {
			C := int64(2*N-2*a-b)*int64(N-2*a) + int64(3*b*N)
			denom := int64(2*N - 2*a + 5*b)
			if denom <= 0 {
				continue
			}
			threshold := C / denom
			kmin := threshold + 1
			if int(kmin) < low {
				kmin = int64(low)
			}
			if int(kmin) <= high {
				sumP += int64(high) - kmin + 1
			}
		}
	}
	for b := low; b <= high; b++ {
		for c := low; c <= high; c++ {
			C := int64(2*N-2*b-c)*int64(N-2*b) + int64(3*c*N)
			denom := int64(2*N - 2*b + 5*c)
			if denom <= 0 {
				continue
			}
			threshold := C / denom
			amin := threshold + 1
			if int(amin) < low {
				amin = int64(low)
			}
			if int(amin) <= high {
				sumQ += int64(high) - amin + 1
			}
		}
	}
	for c := low; c <= high; c++ {
		for a := low; a <= high; a++ {
			C := int64(2*N-2*c-a)*int64(N-2*c) + int64(3*a*N)
			denom := int64(2*N - 2*c + 5*a)
			if denom <= 0 {
				continue
			}
			threshold := C / denom
			bmin := threshold + 1
			if int(bmin) < low {
				bmin = int64(low)
			}
			if int(bmin) <= high {
				sumR += int64(high) - bmin + 1
			}
		}
	}
	return fmt.Sprintf("%d", sumP+sumQ+sumR)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 3
	m := rng.Intn(n/2 + 1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	exp := expectedD(n, m)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
