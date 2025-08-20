package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func solve(x, d int64) bool {
	k := 0
	tmp := x
	for tmp%d == 0 {
		tmp /= d
		k++
	}
	r := tmp
	if k <= 1 {
		return false
	}
	if r != 1 && !isPrime(r) {
		return true
	}
	if k <= 2 {
		return false
	}
	if isPrime(d) {
		return false
	}
	p := int64(math.Sqrt(float64(d)))
	if p*p == d && isPrime(p) && r == p && k == 3 {
		return false
	}
	return true
}

func generateCase(rng *rand.Rand) (string, string) {
	d := int64(rng.Intn(1000) + 2)
	k := rng.Intn(5) + 1
	tmp := int64(1)
	for i := 0; i < k; i++ {
		if tmp > 1e9/d {
			break
		}
		tmp *= d
	}
	r := int64(rng.Intn(50) + 1)
	if r%d == 0 {
		r++
	}
	x := tmp * r
	if x > 1e9 {
		x = tmp
	}
	input := fmt.Sprintf("1\n%d %d\n", x, d)
	expected := "NO\n"
	if solve(x, d) {
		expected = "YES\n"
	}
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
    outStr := strings.TrimSpace(out.String())
    exp := strings.TrimSpace(expected)
    // accept case-insensitive YES/NO
    if !strings.EqualFold(outStr, exp) {
        return fmt.Errorf("expected %q got %q", exp, outStr)
    }
    return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
