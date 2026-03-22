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

// Embedded solver for 1808E3
func solveE3(n, k, m int64) int64 {
	abs := func(x int64) int64 {
		if x < 0 {
			return -x
		}
		return x
	}
	gcd := func(a, b int64) int64 {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}
	power := func(base, exp, mod int64) int64 {
		base %= mod
		if base < 0 {
			base += mod
		}
		if base == 0 {
			if exp == 0 {
				return 1
			}
			return 0
		}
		res := int64(1)
		for exp > 0 {
			if exp%2 == 1 {
				res = (res * base) % mod
			}
			base = (base * base) % mod
			exp /= 2
		}
		return res
	}

	if k%2 == 1 {
		d := gcd(abs(n-2), k)
		term1 := power(k, n, m)
		term2 := power(k-1, n, m)
		ans := (term1 - term2 + m) % m
		diff := (d - 1) % m
		if n%2 == 0 {
			ans = (ans - diff + m) % m
		} else {
			ans = (ans + diff) % m
		}
		return ans
	}
	d := gcd(abs(n-2), k/2)
	half := (m + 1) / 2
	term1 := (power(k, n, m) * half) % m
	term2 := (power(k-2, n, m) * half) % m
	ans := (term1 - term2 + m) % m
	term3 := (power(2, n-1, m) * ((d - 1) % m)) % m
	if n%2 == 0 {
		ans = (ans - term3 + m) % m
	} else {
		ans = (ans + term3) % m
	}
	return ans
}

var primes = []int64{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73}

func generateCase(rng *rand.Rand) string {
	n := rng.Int63n(50) + 1
	k := rng.Int63n(8) + 2
	m := primes[rng.Intn(len(primes))]
	return fmt.Sprintf("%d %d %d\n", n, k, m)
}

func runCase(bin, input string) error {
	var n, k, m int64
	fmt.Sscanf(input, "%d %d %d", &n, &k, &m)
	exp := fmt.Sprintf("%d", solveE3(n, k, m))

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
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE3.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(bin, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
