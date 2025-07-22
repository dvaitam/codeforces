package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveCase(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return ""
	}
	idx := 0
	n, _ := strconv.Atoi(fields[idx])
	idx++
	m, _ := strconv.Atoi(fields[idx])
	idx++
	k, _ := strconv.Atoi(fields[idx])
	idx++
	t, _ := strconv.Atoi(fields[idx])
	idx++
	edges := make([][2]int, k)
	for i := 0; i < k; i++ {
		if idx+1 >= len(fields) {
			return ""
		}
		x, _ := strconv.Atoi(fields[idx])
		idx++
		y, _ := strconv.Atoi(fields[idx])
		idx++
		edges[i] = [2]int{x - 1, y - 1}
	}
	quota := make([]int, t)
	base := k / t
	rem := k % t
	for j := 0; j < t; j++ {
		quota[j] = base
		if j < rem {
			quota[j]++
		}
	}
	L := make([][]int, n)
	for i := range L {
		L[i] = make([]int, t)
	}
	R := make([][]int, m)
	for i := range R {
		R[i] = make([]int, t)
	}
	assign := make([]int, k)
	bestUneven := 1 << 30
	var bestAssign []int
	var dfs func(int)
	dfs = func(pos int) {
		if pos == k {
			uneven := 0
			for i := 0; i < n; i++ {
				mn, mx := L[i][0], L[i][0]
				for j := 1; j < t; j++ {
					if L[i][j] < mn {
						mn = L[i][j]
					}
					if L[i][j] > mx {
						mx = L[i][j]
					}
				}
				if mx-mn > uneven {
					uneven = mx - mn
				}
			}
			for i := 0; i < m; i++ {
				mn, mx := R[i][0], R[i][0]
				for j := 1; j < t; j++ {
					if R[i][j] < mn {
						mn = R[i][j]
					}
					if R[i][j] > mx {
						mx = R[i][j]
					}
				}
				if mx-mn > uneven {
					uneven = mx - mn
				}
			}
			if uneven < bestUneven {
				bestUneven = uneven
				bestAssign = append([]int(nil), assign...)
			}
			return
		}
		e := edges[pos]
		for j := 0; j < t; j++ {
			if quota[j] == 0 {
				continue
			}
			quota[j]--
			assign[pos] = j
			L[e[0]][j]++
			R[e[1]][j]++
			dfs(pos + 1)
			L[e[0]][j]--
			R[e[1]][j]--
			quota[j]++
		}
	}
	dfs(0)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", bestUneven))
	for i := 0; i < k; i++ {
		sb.WriteString(fmt.Sprintf("%d ", bestAssign[i]+1))
	}
	return strings.TrimSpace(sb.String())
}

func run(bin string, input string) (string, error) {
	cmd := execCommand(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func execCommand(bin string) *exec.Cmd {
	return exec.Command(bin)
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
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
