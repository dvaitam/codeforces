package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ---------- embedded solver from accepted solution ----------

func mulMod(a, b, m int64) int64 {
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	_, rem := bits.Div64(hi, lo, uint64(m))
	return int64(rem)
}

func powerMod(base, exp, mod int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp%2 == 1 {
			res = mulMod(res, base, mod)
		}
		base = mulMod(base, base, mod)
		exp /= 2
	}
	return res
}

func factorizeDistinct(n int64) []int64 {
	var factors []int64
	if n%2 == 0 {
		factors = append(factors, 2)
		for n%2 == 0 {
			n /= 2
		}
	}
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			factors = append(factors, i)
			for n%i == 0 {
				n /= i
			}
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

func getPhiFactors(p, e int64) []int64 {
	factors := factorizeDistinct(p - 1)
	if e > 1 {
		found := false
		for _, q := range factors {
			if q == p {
				found = true
				break
			}
		}
		if !found {
			factors = append(factors, p)
		}
	}
	return factors
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	return a / gcd(a, b) * b
}

func solveG(m, x int64) int64 {
	mTemp := m
	var pFactors []int64
	var aFactors []int64
	if mTemp%2 == 0 {
		pFactors = append(pFactors, 2)
		count := int64(0)
		for mTemp%2 == 0 {
			count++
			mTemp /= 2
		}
		aFactors = append(aFactors, count)
	}
	for i := int64(3); i*i <= mTemp; i += 2 {
		if mTemp%i == 0 {
			pFactors = append(pFactors, i)
			count := int64(0)
			for mTemp%i == 0 {
				count++
				mTemp /= i
			}
			aFactors = append(aFactors, count)
		}
	}
	if mTemp > 1 {
		pFactors = append(pFactors, mTemp)
		aFactors = append(aFactors, 1)
	}

	var kMat [][]int64
	var phiMat [][]int64

	for i, p := range pFactors {
		var pArr []int64
		var kArr []int64
		a := aFactors[i]
		phi := int64(1)
		mod := int64(1)
		for e := int64(1); e <= a; e++ {
			if e == 1 {
				phi = p - 1
			} else {
				phi *= p
			}
			mod *= p

			K := phi
			phiFactors := getPhiFactors(p, e)
			for _, q := range phiFactors {
				for K%q == 0 {
					if powerMod(x, K/q, mod) == 1 {
						K /= q
					} else {
						break
					}
				}
			}
			kArr = append(kArr, K)
			pArr = append(pArr, phi)
		}
		kMat = append(kMat, kArr)
		phiMat = append(phiMat, pArr)
	}

	var ans int64 = 0
	var dfs func(idx int, currentPhi int64, currentK int64)
	dfs = func(idx int, currentPhi int64, currentK int64) {
		if idx == len(pFactors) {
			ans += currentPhi / currentK
			return
		}
		dfs(idx+1, currentPhi, currentK)
		for e := int64(1); e <= aFactors[idx]; e++ {
			nextPhi := currentPhi * phiMat[idx][e-1]
			nextK := lcm(currentK, kMat[idx][e-1])
			dfs(idx+1, nextPhi, nextK)
		}
	}
	dfs(0, 1, 1)
	return ans
}

func solveInput(input []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	scanner.Split(bufio.ScanWords)
	var sb strings.Builder
	for {
		if !scanner.Scan() {
			break
		}
		mStr := scanner.Text()
		if !scanner.Scan() {
			break
		}
		xStr := scanner.Text()
		var m, x int64
		fmt.Sscan(mStr, &m)
		fmt.Sscan(xStr, &x)
		fmt.Fprintln(&sb, solveG(m, x))
	}
	return sb.String()
}

// ---------- verifier logic ----------

func runBin(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTest(rng *rand.Rand) []byte {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	for i := 0; i < t; i++ {
		m := rng.Intn(20) + 1
		x := rng.Intn(20) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", m, x))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genTest(rng)
		want := solveInput(input)
		got, err := runBin(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
