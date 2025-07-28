package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type Fenwick struct {
	n    int
	tree []int64
}

func newFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int64, n+2)}
}

func (f *Fenwick) add(i int, v int64) {
	for i <= f.n {
		f.tree[i] += v
		i += i & -i
	}
}

func (f *Fenwick) sum(i int) int64 {
	s := int64(0)
	for i > 0 {
		s += f.tree[i]
		i &= i - 1
	}
	return s
}

func solveCase(people [][2]int) string {
	sort.Slice(people, func(i, j int) bool { return people[i][0] < people[j][0] })
	n := len(people)
	bs := make([]int, n)
	for i := 0; i < n; i++ {
		bs[i] = people[i][1]
	}
	sortedBs := append([]int(nil), bs...)
	sort.Ints(sortedBs)
	comp := make(map[int]int, n)
	for i, v := range sortedBs {
		comp[v] = i + 1
	}
	bit := newFenwick(n)
	ans := int64(0)
	for i := 0; i < n; i++ {
		idx := comp[bs[i]]
		ans += int64(i) - bit.sum(idx)
		bit.add(idx, 1)
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	used := map[int]bool{}
	vals := make([]int, 0, 2*n)
	for len(vals) < 2*n {
		x := rng.Intn(200) - 100
		if !used[x] {
			used[x] = true
			vals = append(vals, x)
		}
	}
	people := make([][2]int, n)
	for i := 0; i < n; i++ {
		a := vals[2*i]
		b := vals[2*i+1]
		if a > b {
			a, b = b, a
		}
		people[i] = [2]int{a, b}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", people[i][0], people[i][1])
	}
	input := sb.String()
	expected := solveCase(people)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := generateCase(rng)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
