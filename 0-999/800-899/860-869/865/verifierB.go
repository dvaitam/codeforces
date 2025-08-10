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

type person struct {
	s int
	a int
	b int
}

func solve(n, S int, arr []person) int64 {
	type item struct{ diff, s int }
	items := make([]item, 0, n+1)
	total := 0
	var base int64
	for _, p := range arr {
		items = append(items, item{diff: p.b - p.a, s: p.s})
		total += p.s
		base += int64(p.s) * int64(p.a)
	}
	pizzas := (total + S - 1) / S
	if rem := pizzas*S - total; rem > 0 {
		items = append(items, item{diff: 0, s: rem})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].diff > items[j].diff })
	ps := make([]int, len(items)+1)
	pg := make([]int64, len(items)+1)
	for i := 0; i < len(items); i++ {
		ps[i+1] = ps[i] + items[i].s
		pg[i+1] = pg[i] + int64(items[i].s)*int64(items[i].diff)
	}
	ans := base
	for x := 0; x <= pizzas; x++ {
		t := x * S
		idx := sort.Search(len(items), func(i int) bool { return ps[i+1] >= t })
		gain := pg[idx]
		if idx < len(items) && t > ps[idx] {
			gain += int64(t-ps[idx]) * int64(items[idx].diff)
		}
		cur := base + gain
		if cur > ans {
			ans = cur
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		S := rng.Intn(10) + 1
		arr := make([]person, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, S)
		for i := 0; i < n; i++ {
			arr[i].s = rng.Intn(20) + 1
			arr[i].a = rng.Intn(20) + 1
			arr[i].b = rng.Intn(20) + 1
			fmt.Fprintf(&sb, "%d %d %d\n", arr[i].s, arr[i].a, arr[i].b)
		}
		expected := fmt.Sprintf("%d", solve(n, S, arr))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:\n%s\n---\ngot:\n%s\n", t+1, sb.String(), expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
