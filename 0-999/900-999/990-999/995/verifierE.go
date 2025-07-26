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

func modPow(base, exp, mod int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

type caseE struct {
	u, v, p int64
}

var primes = []int64{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}

func genCaseE(rng *rand.Rand) caseE {
	p := primes[rng.Intn(len(primes))]
	u := rng.Int63n(p)
	v := rng.Int63n(p)
	return caseE{u: u, v: v, p: p}
}

func runCaseE(bin string, tc caseE) error {
	input := fmt.Sprintf("%d %d %d\n", tc.u, tc.v, tc.p)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	l, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid length: %v", err)
	}
	if l > 200 || l < 0 {
		return fmt.Errorf("invalid length %d", l)
	}
	if len(fields)-1 != l {
		return fmt.Errorf("expected %d numbers got %d", l, len(fields)-1)
	}
	cur := tc.u
	for i := 0; i < l; i++ {
		c := fields[1+i]
		switch c {
		case "1":
			cur = (cur + 1) % tc.p
		case "2":
			cur = (cur + tc.p - 1) % tc.p
		case "3":
			cur = modPow(cur, tc.p-2, tc.p)
		default:
			return fmt.Errorf("invalid command %q", c)
		}
	}
	if cur != tc.v {
		return fmt.Errorf("final value %d expected %d", cur, tc.v)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseE(rng)
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
