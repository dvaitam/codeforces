package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const modE = 1000000007

func addE(a, b int64) int64 {
	a += b
	if a >= modE {
		a -= modE
	}
	return a
}
func mulE(a, b int64) int64 { return (a * b) % modE }

func solveE(n int, a []int64, k int, bad []int64) int64 {
	fac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % modE
	}
	total := fac[n]
	if k == 0 {
		return total
	}
	n1 := n / 2
	n2 := n - n1
	aL := a[:n1]
	aR := a[n1:]
	mapL := make(map[int64]map[int]int)
	for mask := 0; mask < (1 << n1); mask++ {
		var sum int64
		sz := bits.OnesCount(uint(mask))
		for j := 0; j < n1; j++ {
			if mask&(1<<j) != 0 {
				sum += aL[j]
			}
		}
		msz, ok := mapL[sum]
		if !ok {
			msz = make(map[int]int)
			mapL[sum] = msz
		}
		msz[sz]++
	}
	mapR := make(map[int64]map[int]int)
	for mask := 0; mask < (1 << n2); mask++ {
		var sum int64
		sz := bits.OnesCount(uint(mask))
		for j := 0; j < n2; j++ {
			if mask&(1<<j) != 0 {
				sum += aR[j]
			}
		}
		msz, ok := mapR[sum]
		if !ok {
			msz = make(map[int]int)
			mapR[sum] = msz
		}
		msz[sz]++
	}
	computeF := func(x int64) int64 {
		var res int64
		for sumL, cntL := range mapL {
			sumR := x - sumL
			cntR, ok := mapR[sumR]
			if !ok {
				continue
			}
			for sL, cL := range cntL {
				for sR, cR := range cntR {
					s := sL + sR
					ways := fac[s] * fac[n-s] % modE
					res = (res + ways*int64(cL)*int64(cR)) % modE
				}
			}
		}
		return res
	}
	f := make([]int64, k)
	for i := 0; i < k; i++ {
		f[i] = computeF(bad[i])
	}
	var fxy int64
	if k == 2 {
		x, y := bad[0], bad[1]
		if x > y {
			x, y = y, x
		}
		cntL := 1 << n1
		sumFullL := make([]int64, cntL)
		pcL := make([]int, cntL)
		for mask := 1; mask < cntL; mask++ {
			lsb := mask & -mask
			j := bits.TrailingZeros(uint(lsb))
			prev := mask ^ lsb
			sumFullL[mask] = sumFullL[prev] + aL[j]
			pcL[mask] = pcL[prev] + 1
		}
		cntR := 1 << n2
		sumFullR := make([]int64, cntR)
		pcR := make([]int, cntR)
		for mask := 1; mask < cntR; mask++ {
			lsb := mask & -mask
			j := bits.TrailingZeros(uint(lsb))
			prev := mask ^ lsb
			sumFullR[mask] = sumFullR[prev] + aR[j]
			pcR[mask] = pcR[prev] + 1
		}
		dpL := make([]map[int64]map[int]int, cntL)
		dpL[0] = map[int64]map[int]int{0: {0: 1}}
		for mask := 1; mask < cntL; mask++ {
			lsb := mask & -mask
			j := bits.TrailingZeros(uint(lsb))
			prev := mask ^ lsb
			prevMap := dpL[prev]
			cur := make(map[int64]map[int]int, len(prevMap)+1)
			for sum, msz := range prevMap {
				nm := make(map[int]int, len(msz))
				for s, c := range msz {
					nm[s] = c
				}
				cur[sum] = nm
			}
			ai := aL[j]
			for sum, msz := range prevMap {
				nsum := sum + ai
				inn, ok := cur[nsum]
				if !ok {
					inn = make(map[int]int)
					cur[nsum] = inn
				}
				for s, c := range msz {
					inn[s+1] += c
				}
			}
			dpL[mask] = cur
		}
		dpR := make([]map[int64]map[int]int, cntR)
		dpR[0] = map[int64]map[int]int{0: {0: 1}}
		for mask := 1; mask < cntR; mask++ {
			lsb := mask & -mask
			j := bits.TrailingZeros(uint(lsb))
			prev := mask ^ lsb
			prevMap := dpR[prev]
			cur := make(map[int64]map[int]int, len(prevMap)+1)
			for sum, msz := range prevMap {
				nm := make(map[int]int, len(msz))
				for s, c := range msz {
					nm[s] = c
				}
				cur[sum] = nm
			}
			ai := aR[j]
			for sum, msz := range prevMap {
				nsum := sum + ai
				inn, ok := cur[nsum]
				if !ok {
					inn = make(map[int]int)
					cur[nsum] = inn
				}
				for s, c := range msz {
					inn[s+1] += c
				}
			}
			dpR[mask] = cur
		}
		masksR := make(map[int64][]int)
		for mask := 0; mask < cntR; mask++ {
			s := sumFullR[mask]
			if s > y {
				continue
			}
			masksR[s] = append(masksR[s], mask)
		}
		for maskL, sumL := range sumFullL {
			if sumL > y {
				continue
			}
			sumRNeeded := y - sumL
			listR, ok := masksR[sumRNeeded]
			if !ok {
				continue
			}
			tL := pcL[maskL]
			dpLM := dpL[maskL]
			for _, maskR := range listR {
				t := tL + pcR[maskR]
				dpRM := dpR[maskR]
				for sumSL, mszL := range dpLM {
					sumSR := x - sumSL
					mszR, ok2 := dpRM[sumSR]
					if !ok2 {
						continue
					}
					for szL, cL := range mszL {
						for szR, cR := range mszR {
							s := szL + szR
							ways := fac[s] * fac[t-s] % modE * fac[n-t] % modE
							fxy = (fxy + ways*int64(cL)*int64(cR)) % modE
						}
					}
				}
			}
		}
	}
	ans := total
	for i := 0; i < k; i++ {
		ans = (ans - f[i] + modE) % modE
	}
	if k == 2 {
		ans = (ans + fxy) % modE
	}
	return ans
}

func generateCase(rng *rand.Rand) (int, []int64, int, []int64) {
	n := rng.Intn(6) + 1
	a := make([]int64, n)
	sum := int64(0)
	for i := range a {
		a[i] = int64(rng.Intn(10) + 1)
		sum += a[i]
	}
	k := rng.Intn(3) // 0..2
	bad := make([]int64, k)
	for i := 0; i < k; i++ {
		bad[i] = int64(rng.Intn(int(sum)) + 1)
	}
	return n, a, k, bad
}

func runCase(bin string, n int, a []int64, k int, bad []int64) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", n)
	for i, v := range a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')
	fmt.Fprintf(&input, "%d\n", k)
	if k > 0 {
		for i, v := range bad {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprint(v))
		}
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output %q", outStr)
	}
	expect := solveE(n, a, k, bad)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, a, k, bad := generateCase(rng)
		if err := runCase(bin, n, a, k, bad); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
