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

const mod = 12345

// pair of crime type and multiplicity
type crime struct {
	t byte
	m int
}

func solveD(n int64, crimes []crime) int {
	m := len(crimes)
	crat := make([]int, m)
	typ := make([]int, m)
	need := make([]bool, 26)
	for i, c := range crimes {
		typ[i] = int(c.t - 'A')
		crat[i] = c.m
		need[typ[i]] = true
	}
	if m > 0 {
		leaf := make([]bool, m)
		for i := 0; i < m; i++ {
			if !leaf[i] {
				for j := i + 1; j < m; j++ {
					if typ[j] == typ[i] && crat[j] == crat[i] {
						leaf[j] = true
					}
				}
			}
		}
		newT := make([]int, 0, m)
		newC := make([]int, 0, m)
		newT = append(newT, typ[0])
		newC = append(newC, crat[0])
		for i := 1; i < m; i++ {
			if !leaf[i] {
				newT = append(newT, typ[i])
				newC = append(newC, crat[i])
			}
		}
		typ = newT
		crat = newC
		m = len(typ)
	}
	// generate states
	var states [][]int
	var gen func([]int, int)
	gen = func(v []int, t int) {
		if t == m {
			tmp := make([]int, m)
			copy(tmp, v)
			states = append(states, tmp)
			return
		}
		for i := 0; i < crat[t]; i++ {
			gen(append(v, i), t+1)
		}
	}
	gen(nil, 0)
	tot := len(states)
	mn := make([]int, m+1)
	mn[m] = 1
	if m > 0 {
		mn[m-1] = crat[m-1]
		for i := m - 2; i >= 0; i-- {
			mn[i] = mn[i+1] * crat[i]
		}
	}
	g := make([][]int, tot)
	for i := 0; i < tot; i++ {
		g[i] = make([]int, tot)
		for j := 0; j < tot; j++ {
			num := 0
			for t := 0; t < m; t++ {
				num += ((states[i][t] + states[j][t]) % crat[t]) * mn[t+1]
			}
			g[i][j] = num
		}
	}
	nx := make([][]int, tot)
	for i := 0; i < tot; i++ {
		nx[i] = make([]int, 26)
		for j := 0; j < 26; j++ {
			if !need[j] {
				continue
			}
			num := 0
			for t := 0; t < m; t++ {
				add := 0
				if typ[t] == j {
					add = 1
				}
				num += ((states[i][t] + add) % crat[t]) * mn[t+1]
			}
			nx[i][j] = num
		}
	}
	var rec func(int64) []int
	rec = func(n int64) []int {
		if n == 0 {
			ret := make([]int, tot)
			ret[0] = 1
			return ret
		}
		half := rec(n >> 1)
		ret := make([]int, tot)
		for i := 0; i < tot; i++ {
			vi := half[i]
			if vi == 0 {
				continue
			}
			for j := 0; j < tot; j++ {
				ret[g[i][j]] = (ret[g[i][j]] + vi*half[j]) % mod
			}
		}
		if n&1 == 1 {
			tmp := make([]int, tot)
			for j := 0; j < 26; j++ {
				if need[j] {
					tmp[nx[0][j]]++
				}
			}
			ret1 := make([]int, tot)
			for i := 0; i < tot; i++ {
				vi := ret[i]
				if vi == 0 {
					continue
				}
				for j := 0; j < tot; j++ {
					if tmp[j] == 0 {
						continue
					}
					ret1[g[i][j]] = (ret1[g[i][j]] + vi*tmp[j]) % mod
				}
			}
			return ret1
		}
		return ret
	}
	res := rec(n)
	ans := 0
	for i := 0; i < tot; i++ {
		ok := make([]bool, 26)
		for t := 0; t < m; t++ {
			if states[i][t] == 0 {
				ok[typ[t]] = true
			}
		}
		valid := true
		for j := 0; j < 26; j++ {
			if need[j] && !ok[j] {
				valid = false
				break
			}
		}
		if valid {
			ans = (ans + res[i]) % mod
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := rng.Intn(3) + 1
		crimes := make([]crime, c)
		for j := 0; j < c; j++ {
			crimes[j] = crime{t: byte('A' + rng.Intn(3)), m: rng.Intn(3) + 1}
		}
		n := int64(rng.Intn(20))
		expect := solveD(n, crimes)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, c))
		for _, cr := range crimes {
			sb.WriteString(fmt.Sprintf("%c %d\n", cr.t, cr.m))
		}
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != fmt.Sprint(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:%d\ngot:%s\n", i+1, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
