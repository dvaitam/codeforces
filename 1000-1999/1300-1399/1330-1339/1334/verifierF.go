package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Embedded correct solver for 1334/F
func solveF(input string) string {
	idx := 0
	data := []byte(input)
	nextInt := func() int {
		n := len(data)
		for idx < n {
			c := data[idx]
			if (c >= '0' && c <= '9') || c == '-' {
				break
			}
			idx++
		}
		sign := 1
		if data[idx] == '-' {
			sign = -1
			idx++
		}
		val := 0
		for idx < n {
			c := data[idx]
			if c < '0' || c > '9' {
				break
			}
			val = val*10 + int(c-'0')
			idx++
		}
		return sign * val
	}

	n := nextInt()
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = nextInt()
	}
	p := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		x := int64(nextInt())
		p[i] = x
		total += x
	}
	m := nextInt()
	b := make([]int, m+1)
	for i := 1; i <= m; i++ {
		b[i] = nextInt()
	}

	exact := make([]int, n+1)
	for i := 1; i <= m; i++ {
		exact[b[i]] = i
	}

	stage := make([]int, n+1)
	j := 1
	for v := 1; v <= n; v++ {
		for j <= m && b[j] < v {
			j++
		}
		stage[v] = j
	}

	// BIT
	bitN := m
	bit := make([]int64, bitN+2)
	bitAdd := func(i int, delta int64) {
		for i <= bitN {
			bit[i] += delta
			i += i & -i
		}
	}
	bitSum := func(i int) int64 {
		var s int64
		for i > 0 {
			s += bit[i]
			i -= i & -i
		}
		return s
	}

	val := make([]int64, m+1)
	reach := make([]bool, m+1)
	reach[0] = true

	for i := 0; i < n; i++ {
		v := a[i]
		cost := p[i]

		t := stage[v]
		if t <= m && cost > 0 {
			bitAdd(t, cost)
		}

		jj := exact[v]
		if jj != 0 {
			var prev int64
			if jj == 1 {
				prev = 0
			} else {
				if !reach[jj-1] {
					continue
				}
				prev = val[jj-1] + bitSum(jj-1)
			}
			cand := prev + cost - bitSum(jj)
			if !reach[jj] || cand > val[jj] {
				val[jj] = cand
				reach[jj] = true
			}
		}
	}

	if !reach[m] {
		return "NO"
	}
	keep := val[m] + bitSum(m)
	ans := total - keep
	return fmt.Sprintf("YES\n%d", ans)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(8) + 1
		a := make([]int, n)
		pArr := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(n) + 1
			pArr[i] = rng.Intn(11) - 5
		}
		m := rng.Intn(n) + 1
		pool := make([]int, n)
		for i := 0; i < n; i++ {
			pool[i] = i + 1
		}
		rng.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		b := pool[:m]
		sort.Ints(b)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteString("\n")
		for i, v := range pArr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteString("\n")
		input := sb.String()

		exp := solveF(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
