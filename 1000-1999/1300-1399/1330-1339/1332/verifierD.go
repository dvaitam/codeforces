package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) int64 {
	return rng.Int63n(100001) // 0 to 100000
}

func runCase(bin string, k int64) error {
	input := fmt.Sprintf("%d\n", k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	
	output := strings.TrimSpace(out.String())
	return checkOutput(k, output)
}

func checkOutput(k int64, output string) error {
	lines := strings.Fields(output)
	if len(lines) < 2 {
		return fmt.Errorf("insufficient output")
	}
	n, _ := strconv.Atoi(lines[0])
	m, _ := strconv.Atoi(lines[1])
	
	if n < 1 || n > 500 || m < 1 || m > 500 {
		return fmt.Errorf("invalid matrix size %dx%d", n, m)
	}
	
	// Expected tokens: 2 (dimensions) + n*m (matrix)
	if len(lines) != 2 + n*m {
		return fmt.Errorf("expected %d tokens, got %d", 2+n*m, len(lines))
	}
	
	matrix := make([][]int, n)
	idx := 2
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, m)
		for j := 0; j < m; j++ {
			val, err := strconv.Atoi(lines[idx])
			if err != nil {
				return fmt.Errorf("invalid token at (%d,%d): %s", i+1, j+1, lines[idx])
			}
			if val < 0 || val > 300000 {
				return fmt.Errorf("invalid value at (%d,%d): %d", i+1, j+1, val)
			}
			matrix[i][j] = val
			idx++
		}
	}
	
bob := solveBob(n, m, matrix)
	real := solveReal(n, m, matrix)
	
	if int64(real - bob) != k {
		return fmt.Errorf("diff is %d (real %d, bob %d), expected %d", real-bob, real, bob, k)
	}
	return nil
}

func solveBob(n, m int, a [][]int) int {
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, m)
	}
	
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			prev := -1
			if i > 0 {
				prev = max(prev, dp[i-1][j])
			}
			if j > 0 {
				prev = max(prev, dp[i][j-1])
			}
			if prev == -1 { // (0,0)
				dp[i][j] = a[i][j]
			} else {
				dp[i][j] = prev & a[i][j]
			}
		}
	}
	return dp[n-1][m-1]
}

func solveReal(n, m int, a [][]int) int {
    // dp[i][j] stores the set of reachable AND-sums at (i,j).
    // To be efficient, we keep only pareto-optimal values (supermasks).
    dp := make([][][]int, n)
    for i := range dp {
        dp[i] = make([][]int, m)
    }
    
    dp[0][0] = []int{a[0][0]}
    
    for i := 0; i < n; i++ {
        for j := 0; j < m; j++ {
            current := prune(dp[i][j])
            dp[i][j] = current
            
            // Move down
            if i+1 < n {
                nextVal := a[i+1][j]
                for _, v := range current {
                    dp[i+1][j] = append(dp[i+1][j], v & nextVal)
                }
            }
            // Move right
            if j+1 < m {
                nextVal := a[i][j+1]
                for _, v := range current {
                    dp[i][j+1] = append(dp[i][j+1], v & nextVal)
                }
            }
        }
    }
    
    res := 0
    for _, v := range dp[n-1][m-1] {
        if v > res {
            res = v
        }
    }
    return res
}

func prune(vals []int) []int {
    if len(vals) == 0 { return vals }
    
    // Remove duplicates
    unique := make(map[int]bool)
    for _, v := range vals {
        unique[v] = true
    }
    
    var candidates []int
    for v := range unique {
        candidates = append(candidates, v)
    }
    
    // Filter out submasks
    // v is a submask of u if (u & v) == v
    // We keep u if it is NOT a submask of any other element (unless equal)
    // Actually, we want to MAXIMIZE the result.
    // If u has more bits than v (u is supermask), u is always better or equal.
    // So we discard v if there exists u such that (u & v) == v and u != v.
    
    var res []int
    for _, v := range candidates {
        isSub := false
        for _, u := range candidates {
            if u != v && (u & v) == v {
                isSub = true
                break
            }
        }
        if !isSub {
            res = append(res, v)
        }
    }
    return res
}

func max(a, b int) int {
	if a > b { return a }
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		k := generateCase(rng)
		if err := runCase(bin, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}