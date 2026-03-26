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

// ─── embedded correct solver ───

func solveD(input string) string {
	data := []byte(input)
	pos := 0
	nextInt := func() int {
		n := len(data)
		for pos < n && (data[pos] < '0' || data[pos] > '9') {
			pos++
		}
		x := 0
		for pos < n && data[pos] >= '0' && data[pos] <= '9' {
			x = x*10 + int(data[pos]-'0')
			pos++
		}
		return x
	}

	n := nextInt()
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		u := nextInt()
		v := nextInt()
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := make([]int, 1, n)
	stack[0] = 1
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

	ids := make(map[uint64]int, n)
	idNext := 1

	getID := func(cnt, child int) int {
		if cnt == 0 {
			return 1
		}
		key := (uint64(cnt) << 32) | uint64(child)
		if v, ok := ids[key]; ok {
			return v
		}
		idNext++
		ids[key] = idNext
		return idNext
	}

	down := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		u := order[i]
		cnt := 0
		first := 0
		ok := true
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			cnt++
			s := down[v]
			if s == -1 {
				ok = false
				break
			}
			if first == 0 {
				first = s
			} else if s != first {
				ok = false
				break
			}
		}
		if !ok {
			down[u] = -1
		} else if cnt == 0 {
			down[u] = 1
		} else {
			down[u] = getID(cnt, first)
		}
	}

	up := make([]int, n+1)
	ans := -1

	for _, u := range order {
		total := len(adj[u])
		invalid := 0
		d := 0
		sig1, c1 := 0, 0
		sig2, c2 := 0, 0

		if parent[u] != 0 {
			s := up[u]
			if s == -1 {
				invalid++
			} else if d == 0 {
				d = 1
				sig1 = s
				c1 = 1
			} else if s == sig1 {
				c1++
			} else if d == 1 {
				d = 2
				sig2 = s
				c2 = 1
			} else if d == 2 {
				if s == sig2 {
					c2++
				} else {
					d = 3
				}
			}
		}

		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			s := down[v]
			if s == -1 {
				invalid++
			} else if d == 0 {
				d = 1
				sig1 = s
				c1 = 1
			} else if s == sig1 {
				c1++
			} else if d == 1 {
				d = 2
				sig2 = s
				c2 = 1
			} else if d == 2 {
				if s == sig2 {
					c2++
				} else {
					d = 3
				}
			}
		}

		if total == 0 || (invalid == 0 && d == 1) {
			ans = u
			break
		}

		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			t := down[v]
			rem := total - 1
			out := -1
			if rem == 0 {
				out = 1
			} else if invalid > 0 {
				if t == -1 && invalid == 1 && d == 1 {
					out = getID(rem, sig1)
				}
			} else {
				if d == 1 {
					out = getID(rem, sig1)
				} else if d == 2 {
					if t == sig1 && c1 == 1 {
						out = getID(rem, sig2)
					} else if t == sig2 && c2 == 1 {
						out = getID(rem, sig1)
					}
				}
			}
			up[v] = out
		}
	}

	return strconv.Itoa(ans)
}

// ─── verifier ───

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() []byte {
	n := rand.Intn(20) + 1
	if n == 1 {
		return []byte("1\n")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for v := 2; v <= n; v++ {
		p := rand.Intn(v-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", p, v))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		want := solveD(string(input))
		gotRaw, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, string(input))
			os.Exit(1)
		}
		got := strings.TrimSpace(gotRaw)
		if strings.TrimSpace(want) != got {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
