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

type Cup struct {
	c int64
	w int64
}

func expected(n, m int, d int64, phys, info []Cup) int64 {
	sort.Slice(phys, func(i, j int) bool { return phys[i].c > phys[j].c })
	sort.Slice(info, func(i, j int) bool { return info[i].c > info[j].c })
	wp := make([]int64, n+1)
	sp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		wp[i] = wp[i-1] + phys[i-1].w
		sp[i] = sp[i-1] + phys[i-1].c
	}
	wi := make([]int64, m+1)
	si := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		wi[i] = wi[i-1] + info[i-1].w
		si[i] = si[i-1] + info[i-1].c
	}
	ans := int64(0)
	for i := 1; i <= n; i++ {
		remain := d - wp[i]
		if remain < 0 {
			continue
		}
		j := sort.Search(len(wi), func(k int) bool { return wi[k] > remain }) - 1
		if j >= 1 {
			val := sp[i] + si[j]
			if val > ans {
				ans = val
			}
		}
	}
	for j := 1; j <= m; j++ {
		remain := d - wi[j]
		if remain < 0 {
			continue
		}
		i := sort.Search(len(wp), func(k int) bool { return wp[k] > remain }) - 1
		if i >= 1 {
			val := si[j] + sp[i]
			if val > ans {
				ans = val
			}
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

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	if rng.Float64() < 0.1 {
		n = rng.Intn(10) + 1
		m = rng.Intn(10) + 1
	}
	d := int64(rng.Intn(20) + 1)
	phys := make([]Cup, n)
	for i := 0; i < n; i++ {
		phys[i] = Cup{int64(rng.Intn(10) + 1), int64(rng.Intn(5) + 1)}
	}
	info := make([]Cup, m)
	for i := 0; i < m; i++ {
		info[i] = Cup{int64(rng.Intn(10) + 1), int64(rng.Intn(5) + 1)}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", n, m, d)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d %d\n", phys[i].c, phys[i].w)
	}
	for i := 0; i < m; i++ {
		fmt.Fprintf(&b, "%d %d\n", info[i].c, info[i].w)
	}
	return b.String(), expected(n, m, d, phys, info)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
