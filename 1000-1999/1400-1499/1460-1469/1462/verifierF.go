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

func expected(segments [][2]int) int {
	n := len(segments)
	L := make([]int, n)
	R := make([]int, n)
	for i, seg := range segments {
		L[i] = seg[0]
		R[i] = seg[1]
	}
	sortedL := append([]int(nil), L...)
	sortedR := append([]int(nil), R...)
	sort.Ints(sortedL)
	sort.Ints(sortedR)
	maxc := 0
	for i := 0; i < n; i++ {
		x := sort.Search(n, func(j int) bool { return sortedL[j] > R[i] })
		y := sort.Search(n, func(j int) bool { return sortedR[j] >= L[i] })
		c := x - y
		if c > maxc {
			maxc = c
		}
	}
	return n - maxc
}

func runCase(bin string, segs [][2]int) error {
	n := len(segs)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintln(n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", segs[i][0], segs[i][1]))
	}

	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	exp := expected(segs)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Println("could not open testcasesF.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		segs := make([][2]int, n)
		for j := 0; j < n; j++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			l, _ := strconv.Atoi(scanner.Text())
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			r, _ := strconv.Atoi(scanner.Text())
			segs[j] = [2]int{l, r}
		}
		if err := runCase(bin, segs); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
