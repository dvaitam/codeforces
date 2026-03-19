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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// refSolve is the correct embedded reference solver for 1383C.
func refSolve(input []byte) string {
	tokens := strings.Fields(string(input))
	idx := 0
	nextInt := func() int {
		v := 0
		s := tokens[idx]
		idx++
		for _, c := range s {
			v = v*10 + int(c-'0')
		}
		return v
	}

	t := nextInt()
	var results []string
	for ; t > 0; t-- {
		n := nextInt()
		a := make([]byte, n)
		b := make([]byte, n)
		aStr := tokens[idx]
		idx++
		bStr := tokens[idx]
		idx++
		for i := 0; i < n; i++ {
			a[i] = aStr[i]
			b[i] = bStr[i]
		}

		var dir [20][20]bool
		var used [20]bool
		var undir [20][20]bool

		for i := 0; i < n; i++ {
			u := int(a[i] - 'a')
			v := int(b[i] - 'a')
			if u == v {
				continue
			}
			dir[u][v] = true
			used[u] = true
			used[v] = true
			undir[u][v] = true
			undir[v][u] = true
		}

		ans := 0
		var vis [20]bool

		for s := 0; s < 20; s++ {
			if !used[s] || vis[s] {
				continue
			}

			queue := make([]int, 0, 20)
			comp := make([]int, 0, 20)
			queue = append(queue, s)
			vis[s] = true

			for head := 0; head < len(queue); head++ {
				u := queue[head]
				comp = append(comp, u)
				for v := 0; v < 20; v++ {
					if undir[u][v] && !vis[v] {
						vis[v] = true
						queue = append(queue, v)
					}
				}
			}

			m := len(comp)
			if m == 1 {
				continue
			}

			pos := make([]int, 20)
			for i := 0; i < 20; i++ {
				pos[i] = -1
			}
			for i, v := range comp {
				pos[v] = i
			}

			inMask := make([]uint32, m)
			outMask := make([]uint32, m)

			for i, u := range comp {
				for v := 0; v < 20; v++ {
					if dir[u][v] {
						j := pos[v]
						if j >= 0 {
							outMask[i] |= 1 << j
							inMask[j] |= 1 << i
						}
					}
				}
			}

			limit := 1 << m
			acyclic := make([]uint8, limit)
			acyclic[0] = 1
			best := 0

			for mask := 1; mask < limit; mask++ {
				smask := uint32(mask)
				x := smask
				ok := false
				for x != 0 {
					lsb := x & -x
					i := bits.TrailingZeros32(lsb)
					prev := mask ^ int(lsb)
					if acyclic[prev] == 1 {
						if (outMask[i]&smask) == 0 || (inMask[i]&smask) == 0 {
							ok = true
							break
						}
					}
					x -= lsb
				}
				if ok {
					acyclic[mask] = 1
					pc := bits.OnesCount32(smask)
					if pc > best {
						best = pc
					}
				}
			}

			fvs := m - best
			ans += (m - 1) + fvs
		}

		results = append(results, fmt.Sprintf("%d", ans))
	}
	return strings.Join(results, "\n")
}

func genCase(rng *rand.Rand) []byte {
	n := rng.Intn(20) + 1
	a := make([]byte, n)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		a[i] = byte('a' + rng.Intn(20))
		b[i] = byte('a' + rng.Intn(20))
	}
	return []byte(fmt.Sprintf("1\n%d\n%s\n%s\n", n, string(a), string(b)))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		in := genCase(rng)
		exp := refSolve(in)
		got, err := run(cand, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, string(in), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
