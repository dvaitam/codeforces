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

func expected(n, d int) string {
	maxSum := n * (n - 1) / 2
	rem := n - 1
	last := 1
	depth := 1
	minSum := 0
	for rem > 0 {
		cnt := last * 2
		if cnt > rem {
			cnt = rem
		}
		minSum += cnt * depth
		rem -= cnt
		last = cnt
		depth++
	}
	if d < minSum || d > maxSum {
		return "NO"
	}
	parent := make([]int, n+1)
	depthArr := make([]int, n+1)
	childCnt := make([]int, n+1)
	for i := 2; i <= n; i++ {
		parent[i] = i - 1
		childCnt[i-1]++
		depthArr[i] = i - 1
	}
	depthArr[1] = 0
	avail := make([][]int, n)
	for i := 1; i <= n; i++ {
		cap := 2 - childCnt[i]
		if cap > 0 {
			d0 := depthArr[i]
			avail[d0] = append(avail[d0], i)
		}
	}
	leaves := make([][]int, n)
	for i := 1; i <= n; i++ {
		if childCnt[i] == 0 {
			d0 := depthArr[i]
			leaves[d0] = append(leaves[d0], i)
		}
	}
	diff := maxSum - d
	leafMax := n - 1
	for diff > 0 {
		for leafMax > 0 && len(leaves[leafMax]) == 0 {
			leafMax--
		}
		u := leaves[leafMax][len(leaves[leafMax])-1]
		leaves[leafMax] = leaves[leafMax][:len(leaves[leafMax])-1]
		du := depthArr[u]
		dvMin := du - diff - 1
		if dvMin < 0 {
			dvMin = 0
		}
		var dv, v int
		for dv = dvMin; dv < du; dv++ {
			if dv < len(avail) && len(avail[dv]) > 0 {
				v = avail[dv][0]
				break
			}
		}
		nd := dv + 1
		delta := du - nd
		diff -= delta
		old := parent[u]
		oldCap := 2 - childCnt[old]
		childCnt[old]--
		newCapOld := 2 - childCnt[old]
		if old != 0 {
			if oldCap == 0 && newCapOld > 0 {
				avail[depthArr[old]] = append(avail[depthArr[old]], old)
			} else if oldCap > 0 && newCapOld == 0 {
				d0 := depthArr[old]
				for i, x := range avail[d0] {
					if x == old {
						avail[d0] = append(avail[d0][:i], avail[d0][i+1:]...)
						break
					}
				}
			}
		}
		pCapOld := 2 - childCnt[v]
		childCnt[v]++
		pCapNew := 2 - childCnt[v]
		if pCapOld > 0 && pCapNew == 0 {
			for i, x := range avail[dv] {
				if x == v {
					avail[dv] = append(avail[dv][:i], avail[dv][i+1:]...)
					break
				}
			}
		}
		parent[u] = v
		depthArr[u] = nd
		leaves[nd] = append(leaves[nd], u)
		if nd > leafMax {
			leafMax = nd
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 2; i <= n; i++ {
		sb.WriteString(strconv.Itoa(parent[i]))
		if i < n {
			sb.WriteByte(' ')
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
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
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Printf("test %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		d, _ := strconv.Atoi(fields[1])
		exp := expected(n, d)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, d))

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
