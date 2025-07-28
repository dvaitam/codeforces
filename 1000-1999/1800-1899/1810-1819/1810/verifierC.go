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

type testCase struct {
	c   int64
	d   int64
	arr []int
}

func expected(tc testCase) string {
	n := len(tc.arr)
	arr := append([]int(nil), tc.arr...)
	sort.Ints(arr)
	uniq := make([]int, 0, n)
	last := -1
	for _, v := range arr {
		if v != last {
			uniq = append(uniq, v)
			last = v
		}
	}
	cost := func(m int) int64 {
		p := sort.Search(len(uniq), func(i int) bool { return uniq[i] > m })
		del := n - p
		ins := m - p
		return int64(del)*tc.c + int64(ins)*tc.d
	}
	best := cost(1)
	for _, v := range uniq {
		if c := cost(v); c < best {
			best = c
		}
	}
	return fmt.Sprint(best)
}

func runCase(exe string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", len(tc.arr), tc.c, tc.d))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expected(tc)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		c := int64(rng.Intn(5) + 1)
		d := int64(rng.Intn(5) + 1)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(10) + 1
		}
		tc := testCase{c: c, d: d, arr: arr}
		if err := runCase(exe, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
