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

const refMOD = 998244353

func refPower(base, exp int64) int64 {
	res := int64(1)
	base %= refMOD
	for exp > 0 {
		if exp%2 == 1 {
			res = res * base % refMOD
		}
		base = base * base % refMOD
		exp /= 2
	}
	return res
}

func solve(input string) string {
	r := strings.NewReader(input)
	var m int
	fmt.Fscan(r, &m)

	var cnt [100005]int64
	maxA := 0
	for i := 0; i < m; i++ {
		var a int
		var f int64
		fmt.Fscan(r, &a, &f)
		cnt[a] += f
		if a > maxA {
			maxA = a
		}
	}

	mu := make([]int64, maxA+1)
	primes := make([]int, 0)
	isPrime := make([]bool, maxA+1)
	for i := 2; i <= maxA; i++ {
		isPrime[i] = true
	}

	mu[1] = 1
	for i := 2; i <= maxA; i++ {
		if isPrime[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if i*p > maxA {
				break
			}
			isPrime[i*p] = false
			if i%p == 0 {
				mu[i*p] = 0
				break
			} else {
				mu[i*p] = -mu[i]
			}
		}
	}

	ans := int64(0)
	inv8 := int64(873463809)
	P := int64(refMOD)
	P1 := P - 1

	for d := 1; d <= maxA; d++ {
		if mu[d] == 0 {
			continue
		}
		Nd := int64(0)
		S1 := int64(0)
		S2 := int64(0)

		for i := d; i <= maxA; i += d {
			if cnt[i] > 0 {
				Nd += cnt[i]
				val := int64(i)
				cF := cnt[i] % P
				S1 = (S1 + val*cF) % P
				S2 = (S2 + (val*val%P)*cF) % P
			}
		}

		if Nd == 0 {
			continue
		}

		nModP := Nd % P
		nModP1 := Nd % P1

		pow2 := refPower(2, nModP1)
		term1 := (nModP - 2 + P) % P * S2 % P
		term2 := nModP * S1 % P * S1 % P
		Ed := pow2 * inv8 % P * ((term1 + term2) % P) % P

		if mu[d] == 1 {
			ans = (ans + Ed) % P
		} else if mu[d] == -1 {
			ans = (ans - Ed + P) % P
		}
	}

	return fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	m := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", m)
	used := make(map[int]bool)
	for i := 0; i < m; i++ {
		a := rng.Intn(100) + 1
		for used[a] {
			a = rng.Intn(100) + 1
		}
		used[a] = true
		freq := rng.Intn(10) + 1
		fmt.Fprintf(&sb, "%d %d\n", a, freq)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		in := genCase(rng)
		want := solve(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
