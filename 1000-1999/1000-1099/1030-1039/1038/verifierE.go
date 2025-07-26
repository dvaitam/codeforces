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

const V = 4
const INF = int64(4e18)

func bitsCount(x int) int {
	cnt := 0
	for x != 0 {
		cnt++
		x &= x - 1
	}
	return cnt
}

func expected(n int, c1 []int, val []int64, c2 []int) string {
	dist := make([][]int64, V)
	for i := 0; i < V; i++ {
		dist[i] = make([]int64, V)
		for j := 0; j < V; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = INF
			}
		}
	}
	bMask := 0
	var total int64
	for i := 0; i < n; i++ {
		u := c1[i] - 1
		w := c2[i] - 1
		v := val[i]
		total += v
		bMask ^= 1 << u
		bMask ^= 1 << w
		if v < dist[u][w] {
			dist[u][w] = v
			dist[w][u] = v
		}
	}
	for k := 0; k < V; k++ {
		for i := 0; i < V; i++ {
			for j := 0; j < V; j++ {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}
	best := int64(0)
	for fMask := 0; fMask < (1 << V); fMask++ {
		bc := bitsCount(fMask)
		if bc != 0 && bc != 2 {
			continue
		}
		xMask := bMask ^ fMask
		tcnt := bitsCount(xMask)
		var cost int64
		switch tcnt {
		case 0:
			cost = 0
		case 2:
			var u, v int
			idx := 0
			for i := 0; i < V; i++ {
				if xMask&(1<<i) != 0 {
					if idx == 0 {
						u = i
					} else {
						v = i
					}
					idx++
				}
			}
			cost = dist[u][v]
		case 4:
			c1v := dist[0][1] + dist[2][3]
			c2v := dist[0][2] + dist[1][3]
			c3v := dist[0][3] + dist[1][2]
			cost = c1v
			if c2v < cost {
				cost = c2v
			}
			if c3v < cost {
				cost = c3v
			}
		default:
			continue
		}
		if cost >= INF/2 {
			continue
		}
		cur := total - cost
		if cur > best {
			best = cur
		}
	}
	return fmt.Sprintf("%d", best)
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		c1 := make([]int, n)
		c2 := make([]int, n)
		val := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			c1[i], _ = strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			val[i] = int64(v)
			scan.Scan()
			c2[i], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", c1[i], val[i], c2[i]))
		}
		exp := expected(n, c1, val, c2)
		if err := runCase(exe, sb.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
