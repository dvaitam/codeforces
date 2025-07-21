package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const modC = 51123987

func solveCaseC(s string) int {
	n := len(s)
	nextOcc := make([][3]int, n+2)
	for x := 0; x < 3; x++ {
		nextOcc[n][x] = n + 1
		nextOcc[n+1][x] = n + 1
	}
	for i := n - 1; i >= 0; i-- {
		for x := 0; x < 3; x++ {
			nextOcc[i][x] = nextOcc[i+1][x]
		}
		var xi int
		switch s[i] {
		case 'a':
			xi = 0
		case 'b':
			xi = 1
		default:
			xi = 2
		}
		nextOcc[i][xi] = i + 1
	}
	type triple struct{ a, b, c int }
	var dists []triple
	base := n / 3
	rem := n % 3
	if rem == 0 {
		dists = append(dists, triple{base, base, base})
	} else if rem == 1 {
		dists = append(dists, triple{base + 1, base, base})
		dists = append(dists, triple{base, base + 1, base})
		dists = append(dists, triple{base, base, base + 1})
	} else {
		dists = append(dists, triple{base + 1, base + 1, base})
		dists = append(dists, triple{base + 1, base, base + 1})
		dists = append(dists, triple{base, base + 1, base + 1})
	}
	ans := 0
	for _, d := range dists {
		va, vb, vc := d.a, d.b, d.c
		dp := make([][][]int, n+1)
		for i := 0; i <= n; i++ {
			dp[i] = make([][]int, va+1)
			for ca := 0; ca <= va; ca++ {
				dp[i][ca] = make([]int, vb+1)
			}
		}
		dp[0][0][0] = 1
		for pos := 0; pos < n; pos++ {
			for ca := 0; ca <= va; ca++ {
				for cb := 0; cb <= vb; cb++ {
					cur := dp[pos][ca][cb]
					if cur == 0 {
						continue
					}
					ccDone := pos - ca - cb
					for x := 0; x < 3; x++ {
						var remCap int
						switch x {
						case 0:
							remCap = va - ca
						case 1:
							remCap = vb - cb
						default:
							remCap = vc - ccDone
						}
						if remCap <= 0 {
							continue
						}
						nxt := nextOcc[pos][x]
						if nxt > n {
							continue
						}
						Lmin := nxt - pos
						maxLen := n - pos
						if remCap < maxLen {
							maxLen = remCap
						}
						for L := Lmin; L <= maxLen; L++ {
							pos2 := pos + L
							ca2 := ca
							cb2 := cb
							if x == 0 {
								ca2 += L
							} else if x == 1 {
								cb2 += L
							}
							dp[pos2][ca2][cb2] += cur
							if dp[pos2][ca2][cb2] >= modC {
								dp[pos2][ca2][cb2] -= modC
							}
						}
					}
				}
			}
		}
		ans = (ans + dp[n][va][vb]) % modC
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		s := scan.Text()
		expected[i] = fmt.Sprintf("%d", solveCaseC(s))
		_ = n // length already known from s
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
