package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func gFunc(a []int, i, j int) int {
	if i > j {
		return 0
	}
	required := make(map[int]struct{})
	for p := i; p <= j; p++ {
		required[a[p-1]] = struct{}{}
	}
	for x := j; x >= 1; x-- {
		if _, ok := required[a[x-1]]; ok {
			delete(required, a[x-1])
			if len(required) == 0 {
				return x
			}
		}
	}
	return 0
}

func expected(n, q int, a []int, queries [][4]int) []string {
	res := make([]string, q)
	for idx, qu := range queries {
		l, r, x, y := qu[0], qu[1], qu[2], qu[3]
		ans := 0
		for i := l; i <= r; i++ {
			for j := x; j <= y; j++ {
				if i <= j {
					ans += gFunc(a, i, j)
				}
			}
		}
		res[idx] = fmt.Sprintf("%d", ans)
	}
	return res
}

func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		var n, q int
		if _, err := fmt.Fscan(reader, &n, &q); err != nil {
			fmt.Printf("bad test file at case %d: %v\n", caseIdx, err)
			os.Exit(1)
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		queries := make([][4]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(reader, &queries[i][0], &queries[i][1], &queries[i][2], &queries[i][3])
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, qu := range queries {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", qu[0], qu[1], qu[2], qu[3]))
		}
		exp := expected(n, q, a, queries)
		got, err := runProg(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", caseIdx, err, got)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(strings.NewReader(got))
		outScan.Split(bufio.ScanWords)
		for i := 0; i < q; i++ {
			if !outScan.Scan() {
				fmt.Printf("case %d: missing output\n", caseIdx)
				os.Exit(1)
			}
			if outScan.Text() != exp[i] {
				fmt.Printf("case %d failed: expected %s got %s\n", caseIdx, exp[i], outScan.Text())
				os.Exit(1)
			}
		}
		if outScan.Scan() {
			fmt.Printf("case %d: extra output\n", caseIdx)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
