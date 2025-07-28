package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func mulMod(a, b, mod uint64) uint64 {
	return new(big.Int).Mod(new(big.Int).Mul(new(big.Int).SetUint64(a), new(big.Int).SetUint64(b)), new(big.Int).SetUint64(mod)).Uint64()
}

func powMod(a, d, mod uint64) uint64 {
	res := uint64(1)
	a %= mod
	for d > 0 {
		if d&1 == 1 {
			res = mulMod(res, a, mod)
		}
		a = mulMod(a, a, mod)
		d >>= 1
	}
	return res
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	small := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	for _, p := range small {
		if n%p == 0 {
			return n == p
		}
	}
	d := n - 1
	s := 0
	for d&1 == 0 {
		d >>= 1
		s++
	}
	bases := []uint64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
	for _, a := range bases {
		if a%n == 0 {
			continue
		}
		x := powMod(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		composite := true
		for r := 1; r < s; r++ {
			x = mulMod(x, x, n)
			if x == n-1 {
				composite = false
				break
			}
		}
		if composite {
			return false
		}
	}
	return true
}

func pollardsRho(n uint64) uint64 {
	if n%2 == 0 {
		return 2
	}
	if n%3 == 0 {
		return 3
	}
	for {
		c := uint64(rand.Int63n(int64(n-1))) + 1
		x := uint64(rand.Int63n(int64(n)))
		y := x
		d := uint64(1)
		for d == 1 {
			x = (mulMod(x, x, n) + c) % n
			y = (mulMod(y, y, n) + c) % n
			y = (mulMod(y, y, n) + c) % n
			if x > y {
				d = gcd(x-y, n)
			} else {
				d = gcd(y-x, n)
			}
			if d == n {
				break
			}
		}
		if d > 1 && d < n {
			return d
		}
	}
}

func factor(n uint64, res *[]uint64) {
	if n == 1 {
		return
	}
	if isPrime(n) {
		*res = append(*res, n)
		return
	}
	d := pollardsRho(n)
	factor(d, res)
	factor(n/d, res)
}

func divisors(n uint64) []uint64 {
	var pf []uint64
	factor(n, &pf)
	mp := make(map[uint64]int)
	for _, p := range pf {
		mp[p]++
	}
	divs := []uint64{1}
	for p, e := range mp {
		sz := len(divs)
		mul := uint64(1)
		for i := 0; i < e; i++ {
			mul *= p
			for j := 0; j < sz; j++ {
				divs = append(divs, divs[j]*mul)
			}
		}
	}
	sort.Slice(divs, func(i, j int) bool { return divs[i] < divs[j] })
	return divs
}

type kData struct {
	present map[uint64]struct{}
	mex     uint64
}

func solve(reader *bufio.Reader) string {
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return ""
	}
	numbers := make(map[uint64]struct{})
	numbers[0] = struct{}{}
	divisorCache := make(map[uint64][]uint64)
	kMap := make(map[uint64]*kData)
	var sb strings.Builder
	for ; q > 0; q-- {
		var op string
		fmt.Fscan(reader, &op)
		if op == "+" {
			var x uint64
			fmt.Fscan(reader, &x)
			if _, ok := numbers[x]; ok {
				continue
			}
			numbers[x] = struct{}{}
			divs, ok := divisorCache[x]
			if !ok {
				divs = divisors(x)
				divisorCache[x] = divs
			}
			for _, d := range divs {
				if kd, ok := kMap[d]; ok {
					idx := x / d
					kd.present[idx] = struct{}{}
					if idx == kd.mex {
						for {
							if _, ex := kd.present[kd.mex]; ex {
								kd.mex++
							} else {
								break
							}
						}
					}
				}
			}
		} else if op == "-" {
			var x uint64
			fmt.Fscan(reader, &x)
			delete(numbers, x)
			divs := divisorCache[x]
			for _, d := range divs {
				if kd, ok := kMap[d]; ok {
					idx := x / d
					delete(kd.present, idx)
					if idx < kd.mex {
						kd.mex = idx
					}
				}
			}
		} else if op == "?" {
			var k uint64
			fmt.Fscan(reader, &k)
			kd, ok := kMap[k]
			if !ok {
				kd = &kData{present: make(map[uint64]struct{})}
				for num := range numbers {
					if num%k == 0 {
						kd.present[num/k] = struct{}{}
					}
				}
				kd.mex = 0
				for {
					if _, ex := kd.present[kd.mex]; ex {
						kd.mex++
					} else {
						break
					}
				}
				kMap[k] = kd
			} else {
				for {
					if _, ex := kd.present[kd.mex]; ex {
						kd.mex++
					} else {
						break
					}
				}
			}
			fmt.Fprintf(&sb, "%d\n", kd.mex*k)
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	q := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	set := map[uint64]bool{0: true}
	ask := false
	for i := 0; i < q; i++ {
		typ := rng.Intn(3) // 0 add,1 remove,2 query
		if len(set) <= 1 && typ == 1 {
			typ = 0
		}
		if i == q-1 && !ask {
			typ = 2
		}
		if typ == 0 {
			var x uint64
			for {
				x = uint64(rng.Intn(1000000) + 1)
				if !set[x] {
					break
				}
			}
			set[x] = true
			fmt.Fprintf(&sb, "+ %d\n", x)
		} else if typ == 1 {
			var idx int
			if len(set) == 0 {
				i--
				continue
			}
			idx = rng.Intn(len(set))
			var val uint64
			j := 0
			for k := range set {
				if j == idx {
					val = k
					break
				}
				j++
			}
			delete(set, val)
			fmt.Fprintf(&sb, "- %d\n", val)
		} else {
			k := uint64(rng.Intn(10) + 1)
			fmt.Fprintf(&sb, "? %d\n", k)
			ask = true
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
