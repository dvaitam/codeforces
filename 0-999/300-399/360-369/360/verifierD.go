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

// --- Correct solver from accepted solution, adapted as function ---

type Factor struct {
	p int64
	e int64
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func modPow(a, e, mod int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func factorize(x int64) []Factor {
	res := make([]Factor, 0)
	if x%2 == 0 {
		cnt := int64(0)
		for x%2 == 0 {
			x /= 2
			cnt++
		}
		res = append(res, Factor{2, cnt})
	}
	for d := int64(3); d*d <= x; d += 2 {
		if x%d == 0 {
			cnt := int64(0)
			for x%d == 0 {
				x /= d
				cnt++
			}
			res = append(res, Factor{d, cnt})
		}
	}
	if x > 1 {
		res = append(res, Factor{x, 1})
	}
	return res
}

func multiplicativeOrder(a, q, p int64, primes []int64) int64 {
	ord := q
	for _, r := range primes {
		for ord%r == 0 && modPow(a, ord/r, p) == 1 {
			ord /= r
		}
	}
	return ord
}

func phi(x int64, primes []int64) int64 {
	res := x
	for _, r := range primes {
		if x%r == 0 {
			res = res / r * (r - 1)
		}
	}
	return res
}

func solve(n, m int, p int64, a []int64, bArr []int64) int64 {
	var g int64
	for i := 0; i < m; i++ {
		if i == 0 {
			g = bArr[i]
		} else {
			g = gcd(g, bArr[i])
		}
	}

	q := p - 1
	factors := factorize(q)
	primes := make([]int64, len(factors))
	for i, f := range factors {
		primes[i] = f.p
	}

	divisors := []int64{1}
	for _, f := range factors {
		cur := len(divisors)
		mul := int64(1)
		for e := int64(1); e <= f.e; e++ {
			mul *= f.p
			for i := 0; i < cur; i++ {
				divisors = append(divisors, divisors[i]*mul)
			}
		}
	}

	seen := make(map[int64]struct{})
	for _, x := range a {
		ord := multiplicativeOrder(x, q, p, primes)
		s := ord / gcd(ord, g)
		seen[s] = struct{}{}
	}

	orders := make([]int64, 0, len(seen))
	for s := range seen {
		orders = append(orders, s)
	}

	covered := make([]bool, len(divisors))
	for _, s := range orders {
		for i, d := range divisors {
			if s%d == 0 {
				covered[i] = true
			}
		}
	}

	ans := int64(0)
	for i, d := range divisors {
		if covered[i] {
			ans += phi(d, primes)
		}
	}

	return ans
}

// --- end correct solver ---

var primes = []int{3, 5, 7, 11, 13, 17, 19}

func generateCase(r *rand.Rand) string {
	n := r.Intn(3) + 1
	m := r.Intn(3) + 1
	p := primes[r.Intn(len(primes))]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, p))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d", r.Intn(p-1)+1))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d", r.Intn(p)))
		if i+1 < m {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInput(input string) (int, int, int64, []int64, []int64) {
	r := strings.NewReader(input)
	var n, m int
	var p int64
	fmt.Fscan(r, &n, &m, &p)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	b := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &b[i])
	}
	return n, m, p, a, b
}

func run(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		n, m, p, a, b := parseInput(in)
		expect := fmt.Sprintf("%d", solve(n, m, p, a, b))
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
