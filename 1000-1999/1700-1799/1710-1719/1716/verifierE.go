package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Node struct {
	sum, pref, suff, ans int64
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func merge(l, r Node) Node {
	res := Node{}
	res.sum = l.sum + r.sum
	res.pref = max64(l.pref, l.sum+r.pref)
	res.suff = max64(r.suff, r.sum+l.suff)
	cross := l.suff + r.pref
	res.ans = max64(max64(l.ans, r.ans), cross)
	return res
}

func build(a []int64, lvl, start int) []Node {
	if lvl == 0 {
		x := a[start]
		n := Node{sum: x}
		if x > 0 {
			n.pref = x
			n.suff = x
			n.ans = x
		}
		return []Node{n}
	}
	half := 1 << (lvl - 1)
	left := build(a, lvl-1, start)
	right := build(a, lvl-1, start+half)
	size := 1 << lvl
	res := make([]Node, size)
	maskLower := half - 1
	for m := 0; m < size; m++ {
		sub := m & maskLower
		if m&half == 0 {
			res[m] = merge(left[sub], right[sub])
		} else {
			res[m] = merge(right[sub], left[sub])
		}
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
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

func verifyCase(bin string, n int, arr []int64, qs []int) error {
	nodes := build(arr, n, 0)
	mask := 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(qs)))
	for _, k := range qs {
		sb.WriteString(fmt.Sprintf("%d\n", k))
	}
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return err
	}
	expectedVals := make([]string, len(qs))
	for i, k := range qs {
		mask ^= 1 << k
		expectedVals[i] = fmt.Sprint(nodes[mask].ans)
	}
	expected := strings.Join(expectedVals, "\n")
	if strings.TrimSpace(out) != expected {
		return fmt.Errorf("expected:\n%s\n got:\n%s", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		size := 1 << n
		arr := make([]int64, size)
		for j := 0; j < size; j++ {
			arr[j] = rng.Int63n(200) - 100
		}
		q := rng.Intn(20) + 1
		qs := make([]int, q)
		for j := 0; j < q; j++ {
			qs[j] = rng.Intn(n)
		}
		if err := verifyCase(bin, n, arr, qs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
