package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const maxn = 30

var (
	nf, ne, ns  int
	df, de      int
	rf, re_, rs float64
	U           [maxn]bool
	a           [maxn]int
	b           [maxn]float64
)

func Len(a1, b1, c1, d1 float64) float64 {
	lo := math.Max(a1, c1)
	hi := math.Min(b1, d1)
	if hi > lo {
		return hi - lo
	}
	return 0.0
}

func calc() float64 {
	Fc := 2*float64(nf)*rf*float64(df) + 2*float64(ne)*re_*float64(de)
	m := 0
	for i := 0; i < nf+ne+ns; i++ {
		if !U[i] {
			var Df, DeF float64
			xi := float64(i) / 2.0
			for j := 0; j < ns; j++ {
				Df += float64(df) * Len(xi-rf, xi+rf, float64(a[j])-rs, float64(a[j])+rs)
				DeF += float64(de) * Len(xi-re_, xi+re_, float64(a[j])-rs, float64(a[j])+rs)
			}
			Fc += Df
			b[m] = DeF - Df
			m++
		}
	}
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			if b[j] > b[i] {
				b[i], b[j] = b[j], b[i]
			}
		}
	}
	for i := 0; i < ne && i < m; i++ {
		Fc += b[i]
	}
	return Fc
}

func dfs(x, y int, ans *float64) {
	if nf+ne+y < x {
		return
	}
	if x == nf+ne+ns {
		val := calc()
		if val > *ans {
			*ans = val
		}
		return
	}
	U[x] = false
	dfs(x+1, y, ans)
	if y < ns && (x%2 == 0 || U[x-1]) {
		U[x] = true
		a[y] = x / 2
		dfs(x+1, y+1, ans)
	}
}

func solveCase() string {
	ans := 0.0
	dfs(0, 0, &ans)
	return fmt.Sprintf("%.10f", ans)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read testcasesE.txt: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "bad file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expectedOut := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		nf, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		ne, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		ns, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		rf, _ = strconv.ParseFloat(scan.Text(), 64)
		scan.Scan()
		re_, _ = strconv.ParseFloat(scan.Text(), 64)
		scan.Scan()
		rs, _ = strconv.ParseFloat(scan.Text(), 64)
		scan.Scan()
		df, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		de, _ = strconv.Atoi(scan.Text())
		expectedOut[i] = solveCase()
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s", err, errBuf.String())
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for case %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expectedOut[i] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expectedOut[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
