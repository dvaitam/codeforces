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

func expected(times []int, M int) []int {
	res := make([]int, len(times))
	prev := make([]int, 0, len(times))
	sum := 0
	for i, t := range times {
		need := sum + t - M
		if need <= 0 {
			res[i] = 0
		} else {
			tmp := append([]int(nil), prev...)
			sort.Sort(sort.Reverse(sort.IntSlice(tmp)))
			removed := 0
			acc := 0
			for _, v := range tmp {
				acc += v
				removed++
				if acc >= need {
					break
				}
			}
			res[i] = removed
		}
		prev = append(prev, t)
		sum += t
	}
	return res
}

func runCase(bin string, n, M int, times []int) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, M)
	for i, v := range times {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	expectedRes := expected(times, M)
	outs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outs) != len(times) {
		return fmt.Errorf("expected %d numbers, got %d", len(times), len(outs))
	}
	for i, tok := range outs {
		val, err := strconv.Atoi(tok)
		if err != nil || val != expectedRes[i] {
			return fmt.Errorf("position %d expected %d got %s", i, expectedRes[i], tok)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC1.txt")
	if err != nil {
		fmt.Println("could not open testcasesC1.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(sc.Text())
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		M, _ := strconv.Atoi(sc.Text())
		times := make([]int, n)
		for j := 0; j < n; j++ {
			sc.Scan()
			times[j], _ = strconv.Atoi(sc.Text())
		}
		if err := runCase(bin, n, M, times); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
