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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solve(n int, k int64, a []int64) int64 {
	g := int64(0)
	allZero := true
	b := make([]int64, n)
	for i, v := range a {
		b[i] = v - k
		if b[i] != 0 {
			allZero = false
		}
		if b[i] < 0 {
			b[i] = -b[i]
		}
		g = gcd(g, b[i])
	}
	if allZero {
		return 0
	}
	divisors := []int64{}
	for d := int64(1); d*d <= g; d++ {
		if g%d == 0 {
			divisors = append(divisors, d)
			if d != g/d {
				divisors = append(divisors, g/d)
			}
		}
	}
	ans := int64(-1)
	for _, d := range divisors {
		for _, div := range []int64{d, -d} {
			Tval := k + div
			if Tval <= 0 {
				continue
			}
			ok := true
			sumQ := int64(0)
			for i := 0; i < n; i++ {
				if (a[i]-k)%div != 0 {
					ok = false
					break
				}
				q := (a[i] - k) / div
				if q <= 0 {
					ok = false
					break
				}
				sumQ += q
			}
			if ok {
				ops := sumQ - int64(n)
				if ans == -1 || ops < ans {
					ans = ops
				}
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	k := int64(rng.Intn(20) + 1)
	a := make([]int64, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(30) + 1)
		sb.WriteString(fmt.Sprintf("%d ", a[i]))
	}
	sb.WriteString("\n")
	expect := fmt.Sprintf("%d", solve(n, k, append([]int64(nil), a...)))
	return sb.String(), expect
}

func runCase(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
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
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
