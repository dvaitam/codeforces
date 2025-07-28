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

func solve(tp []int, dmg []int64) int64 {
	var fire, frost []int64
	for i, t := range tp {
		if t == 0 {
			fire = append(fire, dmg[i])
		} else {
			frost = append(frost, dmg[i])
		}
	}
	sort.Slice(fire, func(i, j int) bool { return fire[i] > fire[j] })
	sort.Slice(frost, func(i, j int) bool { return frost[i] > frost[j] })
	var sumFire, sumFrost int64
	for _, x := range fire {
		sumFire += x
	}
	for _, x := range frost {
		sumFrost += x
	}
	if len(fire) == len(frost) {
		if len(fire) == 0 {
			return 0
		}
		minVal := fire[len(fire)-1]
		if frost[len(frost)-1] < minVal {
			minVal = frost[len(frost)-1]
		}
		return 2*(sumFire+sumFrost) - minVal
	}
	if len(fire) < len(frost) {
		fire, frost = frost, fire
		sumFire, sumFrost = sumFrost, sumFire
	}
	m := len(frost)
	var extra int64
	for i := 0; i < m; i++ {
		extra += fire[i] + frost[i]
	}
	return sumFire + sumFrost + extra
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) == 0 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		need := 1 + 2*n
		if len(parts) != need {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, need, len(parts))
			os.Exit(1)
		}
		tp := make([]int, n)
		dmg := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[1+i])
			tp[i] = v
		}
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(parts[1+n+i], 10, 64)
			dmg[i] = v
		}
		expect := solve(tp, dmg)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", tp[i]))
		}
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", dmg[i]))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != fmt.Sprintf("%d", expect) {
			fmt.Printf("test %d failed: expected %d got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
