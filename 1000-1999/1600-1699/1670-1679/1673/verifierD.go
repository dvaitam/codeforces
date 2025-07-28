package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MOD int64 = 1_000_000_007

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

type testCase struct {
	b, q, y int64
	c, r, z int64
}

func generateTests() []testCase {
	r := rand.New(rand.NewSource(4))
	tests := []testCase{
		{0, 1, 2, 0, 1, 2},
		{1, 2, 3, 1, 2, 2},
	}
	for len(tests) < 100 {
		b := int64(r.Intn(21) - 10)
		q := int64(r.Intn(10) + 1)
		y := int64(r.Intn(9) + 2)
		c := int64(r.Intn(21) - 10)
		r2 := int64(r.Intn(10) + 1)
		z := int64(r.Intn(9) + 2)
		tests = append(tests, testCase{b, q, y, c, r2, z})
	}
	return tests
}

func expected(t testCase) string {
	b, q, y := t.b, t.q, t.y
	c, rVal, z := t.c, t.r, t.z
	lastB := b + (y-1)*q
	lastC := c + (z-1)*rVal
	if (c-b)%q != 0 || rVal%q != 0 || c < b || lastC > lastB {
		return "0"
	}
	if c-rVal < b || c+z*rVal > lastB {
		return "-1"
	}
	ans := int64(0)
	for d := int64(1); d*d <= rVal; d++ {
		if rVal%d == 0 {
			if lcm(d, q) == rVal {
				x := rVal / d
				ans = (ans + x*x) % MOD
			}
			d2 := rVal / d
			if d2 != d && lcm(d2, q) == rVal {
				x := rVal / d2
				ans = (ans + x*x) % MOD
			}
		}
	}
	return fmt.Sprintf("%d", ans%MOD)
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%d %d %d\n%d %d %d\n", t.b, t.q, t.y, t.c, t.r, t.z)
		want := expected(t)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
