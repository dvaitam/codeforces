package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type segment struct{ l, r int }

type testCaseF struct {
	segs []segment
}

func parseTestcasesF(path string) ([]testCaseF, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cases []testCaseF
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("bad line")
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields)-1 != 2*n {
			return nil, fmt.Errorf("expected %d numbers", 2*n)
		}
		segs := make([]segment, n)
		for i := 0; i < n; i++ {
			l, _ := strconv.Atoi(fields[1+2*i])
			r, _ := strconv.Atoi(fields[2+2*i])
			segs[i] = segment{l, r}
		}
		cases = append(cases, testCaseF{segs})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func solveF(segs []segment) int {
	sort.Slice(segs, func(i, j int) bool {
		li := segs[i].r - segs[i].l
		lj := segs[j].r - segs[j].l
		if li != lj {
			return li < lj
		}
		return segs[i].l < segs[j].l
	})
	type child struct{ l, r, w int }
	n := len(segs)
	dp := make([]int, n)
	for i := 0; i < n; i++ {
		var childs []child
		for j := 0; j < i; j++ {
			if segs[i].l <= segs[j].l && segs[j].r <= segs[i].r {
				childs = append(childs, child{segs[j].l, segs[j].r, 1 + dp[j]})
			}
		}
		if len(childs) == 0 {
			dp[i] = 0
			continue
		}
		sort.Slice(childs, func(a, b int) bool { return childs[a].r < childs[b].r })
		ends := make([]int, len(childs))
		for k := range childs {
			ends[k] = childs[k].r
		}
		local := make([]int, len(childs)+1)
		for k := 1; k <= len(childs); k++ {
			lCur := childs[k-1].l
			wCur := childs[k-1].w
			q := sort.SearchInts(ends, lCur) - 1
			prev := 0
			if q >= 0 {
				prev = local[q+1]
			}
			take := wCur + prev
			if take > local[k-1] {
				local[k] = take
			} else {
				local[k] = local[k-1]
			}
		}
		dp[i] = local[len(childs)]
	}
	childs := make([]child, n)
	for k := 0; k < n; k++ {
		childs[k] = child{segs[k].l, segs[k].r, 1 + dp[k]}
	}
	sort.Slice(childs, func(a, b int) bool { return childs[a].r < childs[b].r })
	ends := make([]int, n)
	for k := range childs {
		ends[k] = childs[k].r
	}
	local := make([]int, n+1)
	for k := 1; k <= n; k++ {
		lCur := childs[k-1].l
		wCur := childs[k-1].w
		q := sort.SearchInts(ends, lCur) - 1
		prev := 0
		if q >= 0 {
			prev = local[q+1]
		}
		take := wCur + prev
		if take > local[k-1] {
			local[k] = take
		} else {
			local[k] = local[k-1]
		}
	}
	return local[n]
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcasesF("testcasesF.txt")
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		n := len(tc.segs)
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, s := range tc.segs {
			sb.WriteString(fmt.Sprintf("%d %d\n", s.l, s.r))
		}
		expected := strconv.Itoa(solveF(tc.segs))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
