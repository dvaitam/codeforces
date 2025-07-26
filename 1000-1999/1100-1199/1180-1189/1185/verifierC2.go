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

func expected(times []int, M int) []int {
	var cnt [101]int
	res := make([]int, len(times))
	for i, t := range times {
		L := M - t
		sum := 0
		kept := 0
		for tv := 1; tv <= 100; tv++ {
			c := cnt[tv]
			if c == 0 {
				continue
			}
			total := c * tv
			if sum+total <= L {
				sum += total
				kept += c
			} else {
				rem := (L - sum) / tv
				if rem > 0 {
					kept += rem
				}
				break
			}
		}
		res[i] = i - kept
		cnt[t]++
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
		fmt.Println("Usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC2.txt")
	if err != nil {
		fmt.Println("could not open testcasesC2.txt:", err)
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
