package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveE(t []int) string {
	n := len(t)
	pocz := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		pocz[i] = total
		total += t[i]
	}
	f := make([]int, total)
	for i := 0; i < total; i++ {
		f[i] = i
	}
	comp := total
	res := make([]int, total)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	sort.Slice(p, func(i, j int) bool { return t[p[i]] > t[p[j]] })
	var find func(int) int
	find = func(a int) int {
		if f[a] != a {
			f[a] = find(f[a])
		}
		return f[a]
	}
	uni := func(a, b int) {
		a = find(a)
		b = find(b)
		if a != b {
			f[a] = b
			comp--
		}
	}
	daj := func(x int) []int {
		arr := make([]int, t[x])
		for i := 0; i < t[x]; i++ {
			arr[i] = pocz[x] + i
		}
		return arr
	}
	ak := list.New()
	for _, v := range daj(p[0]) {
		ak.PushBack(v)
	}
	for idx := 1; idx < n; idx++ {
		i := p[idx]
		nextSize := 3
		if idx < n-1 {
			nextSize = t[p[idx+1]]
		}
		pom := daj(i)
		back := ak.Back().Value.(int)
		uni(pom[len(pom)-1], ak.Back().Value.(int))
		ak.Remove(ak.Back())
		pom = pom[:len(pom)-1]
		uni(pom[len(pom)-1], ak.Back().Value.(int))
		ak.Remove(ak.Back())
		for len(pom) > 1 && len(pom)+ak.Len()-2 >= nextSize {
			pom = pom[:len(pom)-1]
			ak.Remove(ak.Back())
			uni(pom[len(pom)-1], ak.Back().Value.(int))
		}
		pom = pom[:len(pom)-1]
		for j := len(pom) - 1; j >= 0; j-- {
			ak.PushBack(pom[j])
		}
		ak.PushFront(back)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", comp))
	ver := 1
	for i := 0; i < n; i++ {
		var line strings.Builder
		for j := 0; j < t[i]; j++ {
			idx := pocz[i] + j
			root := find(idx)
			if res[root] == 0 {
				res[root] = ver
				ver++
			}
			line.WriteString(fmt.Sprintf("%d ", res[root]))
		}
		sb.WriteString(strings.TrimSpace(line.String()))
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(3) + 1
	t := make([]int, n)
	for i := 0; i < n; i++ {
		t[i] = rng.Intn(3) + 1
	}
	return t
}

func formatInput(t []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(t)))
	for i, v := range t {
		if i+1 == len(t) {
			sb.WriteString(fmt.Sprintf("%d\n", v))
		} else {
			sb.WriteString(fmt.Sprintf("%d ", v))
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases [][]int
	cases = append(cases, []int{1})
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, t := range cases {
		in := formatInput(t)
		exp := solveE(t)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\n got:\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
