package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const refMOD int64 = 998244353

// ---------- embedded reference solver for 1698E ----------

func refSolve(input string) string {
	data := []byte(input)
	idx := 0
	nextInt := func() int {
		for idx < len(data) {
			c := data[idx]
			if c != ' ' && c != '\n' && c != '\r' && c != '\t' {
				break
			}
			idx++
		}
		sign := 1
		if idx < len(data) && data[idx] == '-' {
			sign = -1
			idx++
		}
		val := 0
		for idx < len(data) {
			c := data[idx]
			if c < '0' || c > '9' {
				break
			}
			val = val*10 + int(c-'0')
			idx++
		}
		return sign * val
	}

	out := &strings.Builder{}
	t := nextInt()
	for ; t > 0; t-- {
		n := nextInt()
		s := nextInt()

		a := make([]int, n+1)
		posA := make([]int, n+1)
		for i := 1; i <= n; i++ {
			a[i] = nextInt()
			posA[a[i]] = i
		}

		b := make([]int, n+1)
		for i := 1; i <= n; i++ {
			b[i] = nextInt()
		}

		used := make([]bool, n+1)
		missPos := make([]int, 0, n)
		ok := true

		for i := 1; i <= n; i++ {
			x := b[posA[i]]
			if x == -1 {
				missPos = append(missPos, i)
			} else {
				if used[x] {
					ok = false
				}
				used[x] = true
				if i > x+s {
					ok = false
				}
			}
		}

		if !ok {
			out.WriteString("0\n")
			continue
		}

		suf := make([]int, n+2)
		for v := n; v >= 1; v-- {
			suf[v] = suf[v+1]
			if !used[v] {
				suf[v]++
			}
		}

		var ans int64 = 1
		assigned := 0
		for i := len(missPos) - 1; i >= 0; i-- {
			pos := missPos[i]
			thr := pos - s
			if thr < 1 {
				thr = 1
			}
			choices := suf[thr] - assigned
			if choices <= 0 {
				ans = 0
				break
			}
			ans = ans * int64(choices) % refMOD
			assigned++
		}

		out.WriteString(strconv.FormatInt(ans, 10))
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

// ---------- verifier harness ----------

type Test struct {
	n int
	s int
	a []int
	b []int
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.s))
	for i, v := range t.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	for i, v := range t.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildCase(n int, s int, a []int) Test {
	b := make([]int, n)
	for i := range b {
		b[i] = -1
	}
	return Test{n: n, s: s, a: a, b: b}
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 0, 102)
	tests = append(tests, buildCase(3, 1, []int{2, 1, 3}))
	tests = append(tests, buildCase(3, 2, []int{2, 1, 3}))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		s := rng.Intn(n) + 1
		perm := rng.Perm(n)
		a := make([]int, n)
		for j := range a {
			a[j] = perm[j] + 1
		}
		b := make([]int, n)
		for j := range b {
			b[j] = -1
		}
		numbers := rng.Perm(n)
		k := rng.Intn(n + 1)
		usedIdx := rng.Perm(n)[:k]
		for idx, pos := range usedIdx {
			b[pos] = numbers[idx] + 1
		}
		tests = append(tests, Test{n: n, s: s, a: a, b: b})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = bufio.NewReader(os.Stdin)
	_ = io.Discard

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp := refSolve(input)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:%sExpected:%s\nGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
