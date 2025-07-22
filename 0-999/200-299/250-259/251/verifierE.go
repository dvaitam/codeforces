package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MOD = 1000000007

func add(a, b int) int {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}

func mul(a, b int) int {
	return int((int64(a) * int64(b)) % MOD)
}

func solveE(n int, edges [][2]int) int {
	m := 2 * n
	adj := make([][]int, m)
	deg := make([]int, m)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}
	if n == 1 {
		return 2
	}
	parent := make([]int, m)
	order := make([]int, 0, m)
	parent[0] = -1
	stack := []int{0}
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			stack = append(stack, v)
		}
	}
	dp0 := make([][8]int, m)
	dp1 := make([][8]int, m)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		var children []int
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			children = append(children, v)
		}
		csz := len(children)
		A := make([][8]int, csz+1)
		B := make([][8]int, csz+1)
		A[0][0] = 1
		for j := 0; j < csz; j++ {
			var tmp [8]int
			for s1 := 0; s1 < 8; s1++ {
				if A[j][s1] == 0 {
					continue
				}
				b1 := s1 >> 2
				k1 := s1 & 3
				for s2 := 0; s2 < 8; s2++ {
					v := dp0[children[j]][s2]
					if v == 0 {
						continue
					}
					b2 := s2 >> 2
					k2 := s2 & 3
					nb := b1 | b2
					nk := k1 + k2
					if nk > 3 {
						nk = 3
					}
					idx := (nb << 2) | nk
					tmp[idx] = (tmp[idx] + int((int64(A[j][s1])*int64(v))%MOD)) % MOD
				}
			}
			A[j+1] = tmp
		}
		B[csz][0] = 1
		for j := csz - 1; j >= 0; j-- {
			var tmp [8]int
			for s1 := 0; s1 < 8; s1++ {
				if B[j+1][s1] == 0 {
					continue
				}
				b1 := s1 >> 2
				k1 := s1 & 3
				for s2 := 0; s2 < 8; s2++ {
					v := dp0[children[j]][s2]
					if v == 0 {
						continue
					}
					b2 := s2 >> 2
					k2 := s2 & 3
					nb := b1 | b2
					nk := k1 + k2
					if nk > 3 {
						nk = 3
					}
					idx := (nb << 2) | nk
					tmp[idx] = (tmp[idx] + int((int64(B[j+1][s1])*int64(v))%MOD)) % MOD
				}
			}
			B[j] = tmp
		}
		dp1[u] = A[csz]
		var res0 [8]int
		for j, v := range children {
			var others [8]int
			for s1 := 0; s1 < 8; s1++ {
				if A[j][s1] == 0 {
					continue
				}
				b1 := s1 >> 2
				k1 := s1 & 3
				for s2 := 0; s2 < 8; s2++ {
					c2 := B[j+1][s2]
					if c2 == 0 {
						continue
					}
					b2 := s2 >> 2
					k2 := s2 & 3
					nb := b1 | b2
					nk := k1 + k2
					if nk > 3 {
						nk = 3
					}
					idx := (nb << 2) | nk
					others[idx] = (others[idx] + int((int64(A[j][s1])*int64(c2))%MOD)) % MOD
				}
			}
			ce := deg[u] + deg[v] - 2
			for s := 0; s < 8; s++ {
				if others[s] == 0 {
					continue
				}
				b0 := s >> 2
				k0 := s & 3
				for s2 := 0; s2 < 8; s2++ {
					cnt2 := dp1[v][s2]
					if cnt2 == 0 {
						continue
					}
					b2 := s2 >> 2
					k2 := s2 & 3
					nb := b0 | b2
					nk := k0 + k2
					if nk > 3 {
						nk = 3
					}
					if ce == 0 {
						nb = 1
					} else if ce == 1 {
						nk++
						if nk > 3 {
							nk = 3
						}
					}
					idx := (nb << 2) | nk
					res0[idx] = (res0[idx] + int((int64(others[s])*int64(cnt2))%MOD)) % MOD
				}
			}
		}
		dp0[u] = res0
	}
	ans := dp0[0][0*4+2]
	ans = int((int64(ans) * 4) % MOD)
	return ans
}

type testCaseE struct {
	n     int
	edges [][2]int
	ans   int
}

func genCaseE() testCaseE {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(4) + 1
	m := 2 * n
	if n == 1 {
		return testCaseE{1, [][2]int{{0, 1}}, 2}
	}
	edges := make([][2]int, m-1)
	for i := 1; i < m; i++ {
		p := rand.Intn(i)
		edges[i-1] = [2]int{i, p}
	}
	ans := solveE(n, edges)
	return testCaseE{n, edges, ans}
}

func buildInputE(cs []testCaseE) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cs))
	for _, c := range cs {
		fmt.Fprintln(&sb, c.n)
		for _, e := range c.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := make([]testCaseE, 100)
	for i := range cases {
		cases[i] = genCaseE()
	}
	input := buildInputE(cases)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outputs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outputs) != len(cases) {
		fmt.Printf("expected %d lines got %d\n", len(cases), len(outputs))
		os.Exit(1)
	}
	for i, s := range outputs {
		var val int
		fmt.Sscan(s, &val)
		if val != cases[i].ans {
			fmt.Printf("mismatch on case %d: expected %d got %s\n", i+1, cases[i].ans, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
