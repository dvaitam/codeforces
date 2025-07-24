package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	res, ok := solve(s)
	if !ok {
		fmt.Println("NO")
	} else {
		fmt.Println(res)
	}
}

// solve tries to find the minimal hexadecimal string A such that
// there exists another permutation B of A's digits with |A-B| = D.
// If such A exists it returns the string and true, otherwise "" and false.
func solve(s string) (string, bool) {
	L := len(s)
	// store digits of D from least significant to most significant
	digits := make([]int, L)
	for i := 0; i < L; i++ {
		c := s[L-1-i]
		switch {
		case c >= '0' && c <= '9':
			digits[i] = int(c - '0')
		case c >= 'a' && c <= 'f':
			digits[i] = int(c-'a') + 10
		case c >= 'A' && c <= 'F':
			digits[i] = int(c-'A') + 10
		}
	}

	best := uint64(^uint64(0))
	found := false

	// try both B=A+D and A=B+D
	for _, sign := range []int{1, -1} {
		type key struct {
			pos, carry int
			diff       [16]int
		}
		memo := make(map[key]uint64)
		var dfs func(int, int, [16]int) (uint64, bool)

		dfs = func(pos, carry int, d [16]int) (uint64, bool) {
			k := key{pos, carry, d}
			if v, ok := memo[k]; ok {
				if v == ^uint64(0) {
					return 0, false
				}
				return v, true
			}
			if pos == L {
				if carry == 0 {
					for _, x := range d {
						if x != 0 {
							memo[k] = ^uint64(0)
							return 0, false
						}
					}
					memo[k] = 0
					return 0, true
				}
				memo[k] = ^uint64(0)
				return 0, false
			}
			rem := L - pos
			sumAbs := 0
			for _, x := range d {
				if x > rem || -x > rem {
					memo[k] = ^uint64(0)
					return 0, false
				}
				if x >= 0 {
					sumAbs += x
				} else {
					sumAbs -= x
				}
			}
			if sumAbs > 2*rem {
				memo[k] = ^uint64(0)
				return 0, false
			}

			bestHere := uint64(^uint64(0))
			good := false
			for a := 0; a < 16; a++ {
				t := 0
				nextCarry := 0
				var b int
				if sign == 1 { // B = A + D
					t = a + digits[pos] + carry
					b = t % 16
					nextCarry = t / 16
				} else { // A = B + D -> B = A - D
					t = a - digits[pos] - carry
					if t < 0 {
						t += 16
						nextCarry = 1
					}
					b = t
				}
				d[a]++
				d[b]--
				if val, ok := dfs(pos+1, nextCarry, d); ok {
					cand := val + uint64(a)<<uint(4*pos)
					if cand < bestHere {
						bestHere = cand
					}
					good = true
				}
				d[a]--
				d[b]++
			}
			if good {
				memo[k] = bestHere
				return bestHere, true
			}
			memo[k] = ^uint64(0)
			return 0, false
		}

		val, ok := dfs(0, 0, [16]int{})
		if ok {
			if !found || val < best {
				best = val
				found = true
			}
		}
	}

	if !found {
		return "", false
	}
	out := make([]byte, L)
	for i := 0; i < L; i++ {
		d := (best >> uint(4*i)) & 15
		if d < 10 {
			out[L-1-i] = byte('0' + d)
		} else {
			out[L-1-i] = byte('a' + d - 10)
		}
	}
	return string(out), true
}
