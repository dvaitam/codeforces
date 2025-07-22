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

func expected(line string) string {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return ""
	}
	n, _ := strconv.Atoi(parts[0])
	_ = n
	m, _ := strconv.Atoi(parts[1])
	t, _ := strconv.Atoi(parts[2])
	vols := make([]int, t)
	for i := 0; i < t; i++ {
		v, _ := strconv.Atoi(parts[3+i])
		vols[i] = v
	}
	mod := m
	total := 0
	for mask := 0; mask < (1 << t); mask++ {
		reachable := map[int]bool{0: true}
		used := false
		for i := 0; i < t; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			used = true
			v := vols[i] % mod
			next := make(map[int]bool)
			for r := range reachable {
				next[(r+v)%mod] = true
				next[(r-v+mod)%mod] = true
			}
			reachable = next
		}
		valid := false
		if used {
			if reachable[0] {
				valid = true
			}
		}
		if !valid {
			total++
		}
	}
	const MOD = 1000000007
	return fmt.Sprint(total % MOD)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
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
		expect := expected(line)
		input := line + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
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
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
