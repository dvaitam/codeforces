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

const MOD = 1000000007

type Cow struct {
	posL, posR int
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(input string) string {
	in := strings.NewReader(input)
	var n, m int
	fmt.Fscan(in, &n, &m)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}
	occ := make([][]int, n+1)
	for i, v := range s {
		occ[v] = append(occ[v], i+1)
	}
	cowsByF := make([][]Cow, n+1)
	for i := 0; i < m; i++ {
		var f, h int
		fmt.Fscan(in, &f, &h)
		list := occ[f]
		sz := len(list)
		var posL, posR int
		if h <= sz {
			posL = list[h-1]
			posR = list[sz-h]
		} else {
			posL = -1
			posR = -1
		}
		if posL < 0 && posR < 0 {
			continue
		}
		cowsByF[f] = append(cowsByF[f], Cow{posL, posR})
	}
	pow2 := make([]int, m+1)
	pow2[0] = 1
	for i := 1; i <= m; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}
	bestCnt := 0
	bestWays := 0
	for k := 0; k <= n; k++ {
		total := 0
		ways := 1
		for f := 1; f <= n; f++ {
			list := cowsByF[f]
			if len(list) == 0 {
				continue
			}
			ca, cb, inter := 0, 0, 0
			for _, c := range list {
				inA := c.posL > 0 && c.posL <= k
				inB := c.posR > 0 && c.posR > k
				if inA {
					ca++
				}
				if inB {
					cb++
				}
				if inA && inB {
					inter++
				}
			}
			uni := ca + cb - inter
			if uni == 0 {
				continue
			}
			total += uni
			if inter > 0 {
				ways = ways * pow2[inter] % MOD
			}
		}
		if total > bestCnt {
			bestCnt = total
			bestWays = ways
		} else if total == bestCnt {
			bestWays = (bestWays + ways) % MOD
		}
	}
	return fmt.Sprintf("%d %d", bestCnt, bestWays)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	s := make([]int, n)
	for i := range s {
		s[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, v := range s {
		if i+1 == n {
			fmt.Fprintf(&sb, "%d\n", v)
		} else {
			fmt.Fprintf(&sb, "%d ", v)
		}
	}
	for i := 0; i < m; i++ {
		f := rng.Intn(n) + 1
		h := rng.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d\n", f, h)
	}
	in := sb.String()
	return in, solve(in)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
