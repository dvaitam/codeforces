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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type pair struct{ val, idx int }

func expected(a []int) string {
	n := len(a)
	ps := make([]pair, n)
	for i, v := range a {
		ps[i] = pair{v, i}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].val > ps[j].val })
	res := make([]int, n)
	for i, p := range ps {
		res[p.idx] = i + 1
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", res[i]))
	}
	return sb.String()
}

func genCase(r *rand.Rand) (string, string) {
	n := r.Intn(20) + 1
	a := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		a[i] = r.Intn(1000000000) + 1
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	return sb.String(), expected(a)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		in, exp := genCase(r)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: %v\n", i, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
