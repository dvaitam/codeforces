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

type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(i, v int) {
	for i <= f.n {
		f.bit[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	if i > f.n {
		i = f.n
	}
	s := 0
	for i > 0 {
		s += f.bit[i]
		i -= i & -i
	}
	return s
}

func solveE2(arr []int) string {
	vals := append([]int(nil), arr...)
	sort.Ints(vals)
	vals = unique(vals)
	m := len(vals)
	comp := make(map[int]int, m)
	for i, v := range vals {
		comp[v] = i + 1
	}
	ft := NewFenwick(m)
	total := 0
	ans := 0
	for _, v := range arr {
		idx := comp[v]
		less := ft.Sum(idx - 1)
		greater := total - ft.Sum(idx)
		if less < greater {
			ans += less
		} else {
			ans += greater
		}
		ft.Add(idx, 1)
		total++
	}
	return fmt.Sprintf("%d\n", ans)
}

func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	res := []int{a[0]}
	for _, v := range a[1:] {
		if v != res[len(res)-1] {
			res = append(res, v)
		}
	}
	return res
}

func genCaseE2(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(21) - 10
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := solveE2(arr)
	return input, expect
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := genCaseE2(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
