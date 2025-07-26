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
	n, q int
	a    []int64
	k    []int64
}

func solveCase(tc testCase) string {
	pref := make([]int64, tc.n)
	var sum int64
	for i := 0; i < tc.n; i++ {
		sum += tc.a[i]
		pref[i] = sum
	}
	var dmg int64
	var sb strings.Builder
	for _, x := range tc.k {
		dmg += x
		if dmg >= pref[tc.n-1] {
			dmg = 0
			sb.WriteString(fmt.Sprintf("%d\n", tc.n))
			continue
		}
		idx := sort.Search(len(pref), func(i int) bool { return pref[i] > dmg })
		sb.WriteString(fmt.Sprintf("%d\n", tc.n-idx))
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	tc := testCase{}
	tc.n = rng.Intn(10) + 1
	tc.q = rng.Intn(10) + 1
	tc.a = make([]int64, tc.n)
	for i := range tc.a {
		tc.a[i] = int64(rng.Intn(10) + 1)
	}
	tc.k = make([]int64, tc.q)
	for i := range tc.k {
		tc.k[i] = int64(rng.Intn(20) + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for _, v := range tc.k {
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	return sb.String(), solveCase(tc)
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
