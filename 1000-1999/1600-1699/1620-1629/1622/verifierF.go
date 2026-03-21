package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Embedded solver for 1622F.
func solve1622F(input string) string {
	var n int
	fmt.Sscan(strings.TrimSpace(input), &n)
	if n < 2 {
		// n=1: only 1! = {1}, product is 1, already a perfect square
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		out := bufio.NewWriter(&sb)
		for i := 1; i <= n; i++ {
			if i > 1 {
				out.WriteByte(' ')
			}
			fmt.Fprintf(out, "%d", i)
		}
		out.WriteByte('\n')
		out.Flush()
		return sb.String()
	}

	var rngState uint64 = 1337
	nextRand := func() uint64 {
		rngState ^= rngState << 13
		rngState ^= rngState >> 7
		rngState ^= rngState << 17
		return rngState
	}

	type Hash struct {
		h1, h2 uint64
	}
	xorH := func(a, b Hash) Hash {
		return Hash{a.h1 ^ b.h1, a.h2 ^ b.h2}
	}

	minPrime := make([]int, n+1)
	primes := make([]int, 0, n/10)
	for i := 2; i <= n; i++ {
		if minPrime[i] == 0 {
			minPrime[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p*i > n {
				break
			}
			minPrime[p*i] = p
			if i%p == 0 {
				break
			}
		}
	}

	primeHash := make([]Hash, n+1)
	for _, p := range primes {
		primeHash[p] = Hash{nextRand(), nextRand()}
	}

	valHash := make([]Hash, n+1)
	for i := 2; i <= n; i++ {
		p := minPrime[i]
		valHash[i] = xorH(valHash[i/p], primeHash[p])
	}

	factHash := make([]Hash, n+1)
	var P Hash
	for i := 1; i <= n; i++ {
		factHash[i] = xorH(factHash[i-1], valHash[i])
		P = xorH(P, factHash[i])
	}

	printRes := func(removed []int) string {
		remMap := make(map[int]bool)
		for _, x := range removed {
			remMap[x] = true
		}
		var sb strings.Builder
		fmt.Fprintln(&sb, n-len(removed))
		out := bufio.NewWriter(&sb)
		for i := 1; i <= n; i++ {
			if !remMap[i] {
				if n == 0 {
					break
				}
				var buf [20]byte
				bi := 19
				v := i
				if v == 0 {
					out.WriteByte('0')
				} else {
					for v > 0 {
						buf[bi] = byte(v%10 + '0')
						bi--
						v /= 10
					}
					out.Write(buf[bi+1:])
				}
				out.WriteByte(' ')
			}
		}
		out.WriteByte('\n')
		out.Flush()
		return sb.String()
	}

	if P.h1 == 0 && P.h2 == 0 {
		return printRes(nil)
	}

	type Element struct {
		h  Hash
		id int
	}
	elements := make([]Element, n)
	for i := 1; i <= n; i++ {
		elements[i-1] = Element{factHash[i], i}
	}

	sort.Slice(elements, func(i, j int) bool {
		if elements[i].h.h1 != elements[j].h.h1 {
			return elements[i].h.h1 < elements[j].h.h1
		}
		return elements[i].h.h2 < elements[j].h.h2
	})

	find := func(h Hash) int {
		l, r := 0, n-1
		for l <= r {
			m := (l + r) / 2
			if elements[m].h == h {
				return elements[m].id
			}
			if elements[m].h.h1 < h.h1 || (elements[m].h.h1 == h.h1 && elements[m].h.h2 < h.h2) {
				l = m + 1
			} else {
				r = m - 1
			}
		}
		return -1
	}

	id := find(P)
	if id != -1 {
		return printRes([]int{id})
	}

	for i := 1; i <= n; i++ {
		target := xorH(P, factHash[i])
		id2 := find(target)
		if id2 != -1 && id2 != i {
			return printRes([]int{i, id2})
		}
	}

	k := n / 2
	return printRes([]int{2, k, n})
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(42))
	var tests []string
	edges := []int{1, 2, 3, 4, 5}
	for _, n := range edges {
		tests = append(tests, fmt.Sprintf("%d\n", n))
	}
	for len(tests) < 100 {
		n := rng.Intn(50) + 1
		tests = append(tests, fmt.Sprintf("%d\n", n))
	}
	return tests
}

func runBin(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, input := range tests {
		expect := strings.TrimSpace(solve1622F(input))
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
