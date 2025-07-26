package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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

func verifyOutput(x int64, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid t: %v", err)
	}
	if t < 0 || t > 40 {
		return fmt.Errorf("t out of range")
	}
	opsNeeded := (t + 1) / 2
	if len(fields)-1 < opsNeeded {
		return fmt.Errorf("not enough operations provided")
	}
	idx := 1
	for i := 0; i < t; i++ {
		if i%2 == 0 {
			n, err := strconv.Atoi(fields[idx])
			if err != nil {
				return fmt.Errorf("invalid op: %v", err)
			}
			idx++
			if n < 0 || n > 30 {
				return fmt.Errorf("n out of range")
			}
			mask := int64((1 << uint(n)) - 1)
			x ^= mask
		} else {
			x++
		}
	}
	if idx != 1+opsNeeded {
		return fmt.Errorf("too many numbers in output")
	}
	if bits.OnesCount64(uint64(x+1)) != 1 {
		return fmt.Errorf("result not 2^m-1")
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int64) {
	x := rng.Int63n(1_000_000) + 1
	return fmt.Sprintf("%d\n", x), x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	fixed := []int64{1, 2, 3, 39, 1000000}
	idx := 0
	for ; idx < len(fixed); idx++ {
		inp := fmt.Sprintf("%d\n", fixed[idx])
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if err := verifyOutput(fixed[idx], out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s\n", idx+1, err, inp, out)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		inp, val := generateCase(rng)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if err := verifyOutput(val, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s\n", idx+1, err, inp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
