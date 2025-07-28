package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const maxA = 10000000

var spf []int

func sieve(n int) {
	spf = make([]int, n+1)
	primes := make([]int, 0, 664000)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || p*i > n {
				break
			}
			spf[p*i] = p
		}
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveD(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	var out strings.Builder
	first := true
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if !first {
			out.WriteByte(' ')
		}
		first = false
		if gcd(x, y) != 1 {
			fmt.Fprint(&out, 0)
			continue
		}
		d := y - x
		if d == 1 {
			fmt.Fprint(&out, -1)
			continue
		}
		ans := int(1<<31 - 1)
		tmp := d
		for tmp > 1 {
			p := spf[tmp]
			if p == 0 {
				p = tmp
			}
			if v := p - x%p; v < ans {
				ans = v
			}
			for tmp%p == 0 {
				tmp /= p
			}
		}
		fmt.Fprint(&out, ans)
	}
	return out.String()
}

func genTestD(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d\n", n)
	for i := 0; i < n; i++ {
		x := rng.Intn(10000000-2) + 1
		y := x + rng.Intn(100) + 1
		fmt.Fprintf(&buf, "%d %d\n", x, y)
	}
	return buf.String()
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	sieve(maxA)
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestD(rng)
		expect := solveD(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
