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

type pair struct{ c, t int64 }

func expected(pl []pair, qs []int64) string {
	prefix := make([]int64, len(pl))
	var sum int64
	for i, p := range pl {
		sum += p.c * p.t
		prefix[i] = sum
	}
	var sb strings.Builder
	for i, v := range qs {
		idx := sort.Search(len(prefix), func(j int) bool {
			return prefix[j] >= v
		})
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(idx + 1))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		pl := make([]pair, n)
		var total int64
		for j := 0; j < n; j++ {
			c := rng.Int63n(5) + 1
			t := rng.Int63n(5) + 1
			pl[j] = pair{c, t}
			total += c * t
		}
		maxM := 10
		if int64(maxM) > total {
			maxM = int(total)
		}
		if maxM == 0 {
			maxM = 1
		}
		m := rng.Intn(maxM) + 1
		vals := make(map[int64]struct{})
		for len(vals) < m {
			v := rng.Int63n(total) + 1
			vals[v] = struct{}{}
		}
		qs := make([]int64, 0, m)
		for v := range vals {
			qs = append(qs, v)
		}
		sort.Slice(qs, func(a, b int) bool { return qs[a] < qs[b] })
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&sb, "%d %d\n", pl[j].c, pl[j].t)
		}
		for j, v := range qs {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := expected(pl, qs)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\n got:\n%s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
