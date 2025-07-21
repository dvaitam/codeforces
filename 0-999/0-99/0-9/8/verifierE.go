package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	var k uint64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return ""
	}
	prefix := make([]int8, n)
	for i := range prefix {
		prefix[i] = -1
	}
	numPairs := n / 2
	var dpMemo [51][2][2]uint64
	var dpVis [51][2][2]bool
	var countDP func(i int, revState, crState int) uint64
	countDP = func(i int, revState, crState int) uint64 {
		if prefix[0] == 1 {
			return 0
		}
		if i == numPairs {
			if n%2 == 1 {
				mid := numPairs
				bit := prefix[mid]
				if bit >= 0 {
					if crState == 0 && bit == 1 {
						return 0
					}
					return 1
				}
				if crState == 0 {
					return 1
				}
				return 2
			}
			return 1
		}
		if dpVis[i][revState][crState] {
			return dpMemo[i][revState][crState]
		}
		var res uint64
		j := n - 1 - i
		ai := prefix[i]
		bj := prefix[j]
		if ai >= 0 && bj >= 0 {
			a := ai
			b := bj
			ns := revState
			if ns == 0 {
				if a < b {
					ns = 1
				} else if a > b {
					dpVis[i][revState][crState] = true
					return 0
				}
			}
			nc := crState
			if nc == 0 {
				if a < 1-b {
					nc = 1
				} else if a > 1-b {
					dpVis[i][revState][crState] = true
					return 0
				}
			}
			res = countDP(i+1, ns, nc)
		} else if ai >= 0 {
			a := ai
			for b := int8(0); b <= 1; b++ {
				ns := revState
				if ns == 0 {
					if a < b {
						ns = 1
					} else if a > b {
						continue
					}
				}
				nc := crState
				if nc == 0 {
					if a < 1-b {
						nc = 1
					} else if a > 1-b {
						continue
					}
				}
				prefix[j] = b
				res += countDP(i+1, ns, nc)
				prefix[j] = -1
			}
		} else {
			for a := int8(0); a <= 1; a++ {
				if i == 0 && a != 0 {
					continue
				}
				for b := int8(0); b <= 1; b++ {
					ns := revState
					if ns == 0 {
						if a < b {
							ns = 1
						} else if a > b {
							continue
						}
					}
					nc := crState
					if nc == 0 {
						if a < 1-b {
							nc = 1
						} else if a > 1-b {
							continue
						}
					}
					prefix[i], prefix[j] = a, b
					res += countDP(i+1, ns, nc)
					prefix[i], prefix[j] = -1, -1
				}
			}
		}
		dpVis[i][revState][crState] = true
		dpMemo[i][revState][crState] = res
		return res
	}

	total := countDP(0, 0, 0)
	if k+1 > total {
		return "-1"
	}
	K := k + 1
	prefix[0] = 0
	for pos := 1; pos < n; pos++ {
		prefix[pos] = 0
		for i := 0; i <= numPairs; i++ {
			for rv := 0; rv < 2; rv++ {
				for cr := 0; cr < 2; cr++ {
					dpVis[i][rv][cr] = false
				}
			}
		}
		cnt0 := countDP(0, 0, 0)
		if K <= cnt0 {
			continue
		}
		K -= cnt0
		prefix[pos] = 1
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		if prefix[i] == 0 {
			out[i] = '0'
		} else {
			out[i] = '1'
		}
	}
	return string(out)
}

type test struct{ input, expected string }

func generateTests() []test {
	rand.Seed(555)
	var tests []test
	for len(tests) < 100 {
		n := rand.Intn(8) + 2
		var limit uint64 = 1
		if n < 20 {
			limit = 1 << n
		} else {
			limit = 1000
		}
		k := rand.Uint64()%limit + 1
		inp := fmt.Sprintf("%d %d\n", n, k)
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.expected {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
