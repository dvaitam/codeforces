package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const MOD = 1000000007

var fact []int64
var invFact []int64

func power(base, exp int64) int64 {
	var res int64 = 1
	base %= MOD
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % MOD
		}
		base = (base * base) % MOD
		exp /= 2
	}
	return res
}

func initComb(n int) {
	fact = make([]int64, n+1)
	invFact = make([]int64, n+1)
	fact[0] = 1
	invFact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = (fact[i-1] * int64(i)) % MOD
	}
	invFact[n] = power(fact[n], MOD-2)
	for i := n - 1; i >= 1; i-- {
		invFact[i] = (invFact[i+1] * int64(i+1)) % MOD
	}
}

func C(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	num := fact[n]
	den := (invFact[k] * invFact[n-k]) % MOD
	return (num * den) % MOD
}

type Fixed struct {
	U int
	P int
}

func solveD(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}

	initComb(n)

	var fixed []Fixed
	for i := 1; i <= n; i++ {
		var u int
		fmt.Fscan(reader, &u)
		if u > 0 {
			fixed = append(fixed, Fixed{U: u, P: i})
		}
	}

	sort.Slice(fixed, func(i, j int) bool {
		return fixed[i].U > fixed[j].U
	})

	states := make(map[int]int64)
	states[1] = 1
	v_prev := n

	for _, f := range fixed {
		U := f.U
		P := f.P
		delta := v_prev - U

		if delta < 0 {
			return "0"
		}

		newStates := make(map[int]int64)

		for L_prev, ways := range states {
			L_new1 := P
			new_L1 := P + 1
			k1 := L_new1 - L_prev
			if k1 >= 0 && k1 <= delta {
				ways_add := (ways * C(delta, k1)) % MOD
				newStates[new_L1] = (newStates[new_L1] + ways_add) % MOD
			}

			if U > 1 {
				L_new2 := P - U + 1
				new_L2 := P - U + 1
				k2 := L_new2 - L_prev
				if k2 >= 0 && k2 <= delta {
					ways_add := (ways * C(delta, k2)) % MOD
					newStates[new_L2] = (newStates[new_L2] + ways_add) % MOD
				}
			}
		}

		states = newStates
		v_prev = U - 1
	}

	var ans int64 = 0
	for _, ways := range states {
		ans = (ans + ways) % MOD
	}

	if v_prev > 0 {
		ans = (ans * power(2, int64(v_prev-1))) % MOD
	}

	return fmt.Sprintf("%d", ans)
}

func genCaseD(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	perm := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			val := perm[i] + 1
			sb.WriteString(fmt.Sprint(val))
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCaseD(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(out.String()))
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
		in := genCaseD(rng)
		expect := solveD(in)
		if err := runCaseD(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
