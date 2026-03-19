package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const refINF int64 = 1 << 60

func refPointUpdate(sum, mn []int64, size, pos int, delta int64) {
	i := size + pos - 1
	sum[i] += delta
	mn[i] = sum[i]
	for i >>= 1; i > 0; i >>= 1 {
		ls := sum[i<<1]
		sum[i] = ls + sum[i<<1|1]
		v := mn[i<<1]
		t := ls + mn[i<<1|1]
		if t < v {
			v = t
		}
		mn[i] = v
	}
}

func refPrefixMin(sum, mn []int64, size, r int) int64 {
	l := size
	rr := size + r
	lsum, lmn := int64(0), refINF
	rmn := refINF
	for l < rr {
		if l&1 == 1 {
			t := lsum + mn[l]
			if t < lmn {
				lmn = t
			}
			lsum += sum[l]
			l++
		}
		if rr&1 == 1 {
			rr--
			t := sum[rr] + rmn
			if mn[rr] < t {
				rmn = mn[rr]
			} else {
				rmn = t
			}
		}
		l >>= 1
		rr >>= 1
	}
	t := lsum + rmn
	if t < lmn {
		return t
	}
	return lmn
}

// refSolve is the correct embedded reference solver for 1603D.
func refSolve(input string) string {
	tokens := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v := 0
		s := tokens[idx]
		idx++
		for i := 0; i < len(s); i++ {
			v = v*10 + int(s[i]-'0')
		}
		return v
	}

	t := nextInt()

	ns := make([]int, t)
	ks := make([]int, t)
	maxN, maxK := 0, 0
	for i := 0; i < t; i++ {
		n := nextInt()
		k := nextInt()
		ns[i] = n
		ks[i] = k
		if n > maxN {
			maxN = n
		}
		if k > maxK {
			maxK = k
		}
	}

	globalThreshold := int(bits.Len(uint(maxN)))
	preK := globalThreshold - 1
	if preK > maxK {
		preK = maxK
	}

	var ans [][]int64
	if preK >= 1 {
		ans = make([][]int64, preK+1)

		prev := make([]int64, maxN+1)
		a1 := make([]int64, maxN+1)
		for n := 1; n <= maxN; n++ {
			v := int64(n) * int64(n+1) / 2
			prev[n] = v
			a1[n] = v
		}
		ans[1] = a1

		if preK >= 2 {
			phi := make([]int, maxN+1)
			for i := 0; i <= maxN; i++ {
				phi[i] = i
			}
			if maxN >= 1 {
				phi[1] = 1
			}
			for i := 2; i <= maxN; i++ {
				if phi[i] == i {
					for j := i; j <= maxN; j += i {
						phi[j] -= phi[j] / i
					}
				}
			}

			cnt := make([]int, maxN+1)
			for d := 1; d*2 <= maxN; d++ {
				for e := d * 2; e <= maxN; e += d {
					cnt[e]++
				}
			}
			divs := make([][]int, maxN+1)
			for e := 2; e <= maxN; e++ {
				if cnt[e] > 0 {
					divs[e] = make([]int, 0, cnt[e])
				}
			}
			for d := 1; d*2 <= maxN; d++ {
				for e := d * 2; e <= maxN; e += d {
					divs[e] = append(divs[e], d)
				}
			}

			size := 1
			for size < maxN {
				size <<= 1
			}
			sum := make([]int64, size<<1)
			mn := make([]int64, size<<1)

			for k := 2; k <= preK; k++ {
				for i := 1; i < size<<1; i++ {
					sum[i] = 0
					mn[i] = refINF
				}

				prevB := int64(0)
				for l := 1; l <= maxN; l++ {
					var b int64
					if l < k {
						b = refINF
					} else {
						b = prev[l-1] - int64(l) + 1
					}
					d := b - prevB
					p := size + l - 1
					sum[p] = d
					mn[p] = d
					prevB = b
				}

				for i := size - 1; i >= 1; i-- {
					ls := sum[i<<1]
					sum[i] = ls + sum[i<<1|1]
					v := mn[i<<1]
					t2 := ls + mn[i<<1|1]
					if t2 < v {
						v = t2
					}
					mn[i] = v
				}

				cur := make([]int64, maxN+1)
				for r := 1; r <= maxN; r++ {
					if r > 1 {
						refPointUpdate(sum, mn, size, 1, int64(r-1))
					}
					for _, d := range divs[r] {
						w := int64(phi[r/d])
						refPointUpdate(sum, mn, size, d+1, -w)
					}
					minv := refPrefixMin(sum, mn, size, r)
					if minv >= refINF/2 {
						cur[r] = refINF
					} else {
						cur[r] = int64(r) + minv
					}
				}

				ans[k] = cur
				prev = cur
			}
		}
	}

	var results []string
	for i := 0; i < t; i++ {
		n, k := ns[i], ks[i]
		var v int64
		if k >= int(bits.Len(uint(n))) {
			v = int64(n)
		} else {
			v = ans[k][n]
		}
		results = append(results, fmt.Sprintf("%d", v))
	}
	return strings.Join(results, "\n")
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(100) + 1
	k := rng.Intn(n) + 1
	return fmt.Sprintf("%d %d\n", n, k)
}

func runCmd(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := "1\n" + generateCase(rng)
		expected := strings.TrimSpace(refSolve(in))
		got, err := runCmd(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
