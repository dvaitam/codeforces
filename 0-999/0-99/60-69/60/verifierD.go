package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded correct solver for 60D.
func solveD(input string) string {
	data := []byte(input)
	ptr := 0

	nextInt := func() int {
		n := len(data)
		for ptr < n {
			c := data[ptr]
			if c >= '0' && c <= '9' {
				break
			}
			ptr++
		}
		val := 0
		for ptr < n {
			c := data[ptr]
			if c < '0' || c > '9' {
				break
			}
			val = val*10 + int(c-'0')
			ptr++
		}
		return val
	}

	buildSPF := func(n int) []uint32 {
		spf := make([]uint32, n+1)
		if n < 2 {
			return spf
		}
		primes := make([]int, 0, 700000)
		for i := 2; i <= n; i++ {
			if spf[i] == 0 {
				spf[i] = uint32(i)
				primes = append(primes, i)
			}
			si := int(spf[i])
			limit := n / i
			for j := 0; j < len(primes); j++ {
				p := primes[j]
				if p > si || p > limit {
					break
				}
				spf[p*i] = uint32(p)
			}
		}
		return spf
	}

	n := nextInt()
	vals := make([]int, n)
	maxA := 0
	needSieve := false
	for i := 0; i < n; i++ {
		x := nextInt()
		vals[i] = x
		if x > maxA {
			maxA = x
		}
		if x > 1 && ((x&1) == 1 || (x&3) == 0) {
			needSieve = true
		}
	}

	if !needSieve {
		return strconv.Itoa(n)
	}

	idxByVal := make([]uint32, maxA+1)
	for i, x := range vals {
		idxByVal[x] = uint32(i + 1)
	}

	spf := buildSPF(maxA)

	parent := make([]uint32, n)
	size := make([]uint32, n)
	for i := 0; i < n; i++ {
		parent[i] = uint32(i)
		size[i] = 1
	}
	comps := n

	var find func(x int) int
	find = func(x int) int {
		for parent[x] != uint32(x) {
			parent[x] = parent[int(parent[x])]
			x = int(parent[x])
		}
		return x
	}

	union := func(a, b int) {
		ra := find(a)
		rb := find(b)
		if ra == rb {
			return
		}
		if size[ra] < size[rb] {
			ra, rb = rb, ra
		}
		parent[rb] = uint32(ra)
		size[ra] += size[rb]
		comps--
	}

	for i, x := range vals {
		if x > 1 && (x&1) == 1 {
			temp := x
			var blocks [8]int
			cnt := 0
			for temp > 1 {
				p := int(spf[temp])
				pp := 1
				for temp%p == 0 {
					pp *= p
					temp /= p
				}
				blocks[cnt] = pp
				cnt++
			}
			lim := 1 << cnt
			for mask := 0; mask < lim; mask++ {
				s := 1
				for j := 0; j < cnt; j++ {
					if (mask>>j)&1 != 0 {
						s *= blocks[j]
					}
				}
				t := x / s
				if s >= t {
					continue
				}
				ss := int64(s)
				tt := int64(t)

				even := (tt*tt - ss*ss) >> 1
				if even <= int64(maxA) {
					id := idxByVal[int(even)]
					if id != 0 {
						union(i, int(id-1))
					}
				}

				hyp := (tt*tt + ss*ss) >> 1
				if hyp <= int64(maxA) {
					id := idxByVal[int(hyp)]
					if id != 0 {
						union(i, int(id-1))
					}
				}
			}
		} else if (x & 3) == 0 {
			k := x >> 1
			temp := k
			block2 := 1
			for (temp & 1) == 0 {
				block2 <<= 1
				temp >>= 1
			}

			var blocks [8]int
			cnt := 0
			totalOdd := 1
			for temp > 1 {
				p := int(spf[temp])
				pp := 1
				for temp%p == 0 {
					pp *= p
					temp /= p
				}
				blocks[cnt] = pp
				cnt++
				totalOdd *= pp
			}

			lim := 1 << cnt
			for mask := 0; mask < lim; mask++ {
				prod := 1
				for j := 0; j < cnt; j++ {
					if (mask>>j)&1 != 0 {
						prod *= blocks[j]
					}
				}
				evenFactor := block2 * prod
				oddFactor := totalOdd / prod

				mv := evenFactor
				nv := oddFactor
				if mv < nv {
					mv, nv = nv, mv
				}

				mm := int64(mv)
				n64 := int64(nv)

				oddLeg := mm*mm - n64*n64
				if oddLeg <= int64(maxA) {
					id := idxByVal[int(oddLeg)]
					if id != 0 {
						union(i, int(id-1))
					}
				}

				hyp := mm*mm + n64*n64
				if hyp <= int64(maxA) {
					id := idxByVal[int(hyp)]
					if id != 0 {
						union(i, int(id-1))
					}
				}
			}
		}
	}

	return strconv.Itoa(comps)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rng.Intn(40) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		expect := solveD(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
