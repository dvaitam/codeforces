package main

import (
	"bufio"
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

const mod = 1000000007

// oracleSolve computes the expected output for a given (n, d, x) input
// by replicating the logic from the reference solution.
func oracleSolve(n, d int, xIn int64) []int {
	x := xIn
	a := make([]int, n)
	pos := make([]int, n+1)
	for i := 0; i < n; i++ {
		a[i] = i + 1
		pos[i+1] = i
		x = (x*37 + 10007) % int64(mod)
		j := int(x % int64(i+1))
		if i != j {
			vi, vj := a[i], a[j]
			a[i], a[j] = vj, vi
			pos[vi], pos[vj] = j, i
		}
	}

	bArr := make([]byte, n)
	for i := 0; i < d; i++ {
		bArr[i] = 1
	}
	for i := 0; i < n; i++ {
		x = (x*37 + 10007) % int64(mod)
		j := int(x % int64(i+1))
		bArr[i], bArr[j] = bArr[j], bArr[i]
	}

	m := (n + 63) >> 6
	B := make([]uint64, m)
	minOne := n
	for i := 0; i < n; i++ {
		if bArr[i] != 0 {
			B[i>>6] |= uint64(1) << uint(i&63)
			if minOne == n {
				minOne = i
			}
		}
	}

	ans := make([]int, n)
	for i := 0; i < n; i++ {
		ans[i] = -1
	}

	if minOne < n {
		U := make([]uint64, m)
		for i := minOne; i < n; i++ {
			U[i>>6] |= uint64(1) << uint(i&63)
		}
		remaining := n - minOne

		for v := n; v >= 1 && remaining > 0; v-- {
			p := pos[v]
			ws := p >> 6
			bs := uint(p & 63)

			if bs == 0 {
				src := 0
				for w := ws; w < m; w++ {
					inter := U[w] & B[src]
					if inter != 0 {
						U[w] &^= inter
						remaining -= bits.OnesCount64(inter)
						base := w << 6
						for inter != 0 {
							tz := bits.TrailingZeros64(inter)
							ans[base+tz] = v - 1
							inter &= inter - 1
						}
					}
					src++
				}
			} else {
				rsh := 64 - bs
				var prev uint64
				src := 0
				for w := ws; w < m; w++ {
					cur := B[src]
					shifted := (cur << bs) | prev
					prev = cur >> rsh
					inter := U[w] & shifted
					if inter != 0 {
						U[w] &^= inter
						remaining -= bits.OnesCount64(inter)
						base := w << 6
						for inter != 0 {
							tz := bits.TrailingZeros64(inter)
							ans[base+tz] = v - 1
							inter &= inter - 1
						}
					}
					src++
				}
			}
		}
	}

	return ans
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(6) + 1
	d := rng.Intn(n) + 1
	x := rng.Int63n(mod)
	if x == 27777500 {
		x++
	}
	input := fmt.Sprintf("%d %d %d\n", n, d, x)
	return input, oracleSolve(n, d, x)
}

func runCase(bin string, input string, expect []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	got := make([]int, 0, len(expect))
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		got = append(got, v)
	}
	if len(got) != len(expect) {
		return fmt.Errorf("expected %d numbers got %d", len(expect), len(got))
	}
	for i, v := range expect {
		if got[i] != v {
			return fmt.Errorf("pos %d expected %d got %d", i, v, got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
