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

type Test struct {
	n int
	k int64
	x int64
	a []int64
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.k, tc.x))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(tc Test) string {
	levels := append([]int64(nil), tc.a...)
	sort.Slice(levels, func(i, j int) bool { return levels[i] < levels[j] })
	gaps := make([]int64, 0)
	for i := 1; i < tc.n; i++ {
		diff := levels[i] - levels[i-1]
		if diff > tc.x {
			gaps = append(gaps, (diff-1)/tc.x)
		}
	}
	sort.Slice(gaps, func(i, j int) bool { return gaps[i] < gaps[j] })
	groups := len(gaps) + 1
	k := tc.k
	for _, g := range gaps {
		if k >= g {
			k -= g
			groups--
		} else {
			break
		}
	}
	return fmt.Sprintf("%d", groups)
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(20) + 1
	k := rng.Int63n(10)
	x := rng.Int63n(10) + 1
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(50)
	}
	return Test{n: n, k: k, x: x, a: a}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		tc := genTest(rng)
		expect := expected(tc)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
