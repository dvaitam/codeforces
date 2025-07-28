package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func maxCoins(a []int, n, k int) int {
	type state struct{ mask, ptr int }
	memo := make(map[state]int)
	var dfs func(mask, ptr int) int
	dfs = func(mask, ptr int) int {
		st := state{mask, ptr}
		if v, ok := memo[st]; ok {
			return v
		}
		if mask == 0 && ptr >= n {
			return 0
		}
		indices := make([]int, 0, bits.OnesCount(uint(mask)))
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				indices = append(indices, i)
			}
		}
		best := 0
		for i := 0; i < len(indices); i++ {
			for j := i + 1; j < len(indices); j++ {
				m := mask &^ (1 << indices[i]) &^ (1 << indices[j])
				p := ptr
				if p < n {
					m |= 1 << p
					p++
				}
				if p < n {
					m |= 1 << p
					p++
				}
				gain := 0
				if a[indices[i]] == a[indices[j]] {
					gain = 1
				}
				val := gain + dfs(m, p)
				if val > best {
					best = val
				}
			}
		}
		memo[st] = best
		return best
	}
	startMask := (1 << k) - 1
	return dfs(startMask, k)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Println("failed to open testcasesF.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		a := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[2+i])
			a[i] = v
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		want := strconv.Itoa(maxCoins(a, n, k))
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", idx, want, strings.TrimSpace(got), sb.String())
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
