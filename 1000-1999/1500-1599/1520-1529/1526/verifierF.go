package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// interactor acts as the CF judge for problem 1526F.
// It generates a hidden permutation p[1..n] with p[1]<p[2],
// feeds it to the candidate binary via interactive queries,
// and checks the final answer.
func interactor(bin string, n int, perm []int) error {
	cmd := exec.Command(bin)
	candIn, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	candOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start candidate: %v", err)
	}

	wr := bufio.NewWriter(candIn)
	rd := bufio.NewReader(candOut)

	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	median := func(a, b, c int) int {
		s := []int{a, b, c}
		sort.Ints(s)
		return s[1]
	}

	// Send: t=1, then n
	fmt.Fprintf(wr, "1\n%d\n", n)
	wr.Flush()

	queries := 0
	maxQ := 2*n + 420

	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			cmd.Process.Kill()
			return fmt.Errorf("read error: %v", err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == '?' {
			queries++
			if queries > maxQ {
				fmt.Fprintf(wr, "-1\n")
				wr.Flush()
				cmd.Process.Kill()
				return fmt.Errorf("too many queries: %d > %d", queries, maxQ)
			}
			var a, b, c int
			fmt.Sscanf(line, "? %d %d %d", &a, &b, &c)
			if a < 1 || a > n || b < 1 || b > n || c < 1 || c > n || a == b || b == c || a == c {
				fmt.Fprintf(wr, "-1\n")
				wr.Flush()
				cmd.Process.Kill()
				return fmt.Errorf("invalid query: %s", line)
			}
			res := median(abs(perm[a]-perm[b]), abs(perm[b]-perm[c]), abs(perm[a]-perm[c]))
			fmt.Fprintf(wr, "%d\n", res)
			wr.Flush()
		} else if line[0] == '!' {
			// parse answer
			parts := strings.Fields(line)
			if len(parts) != n+1 {
				fmt.Fprintf(wr, "-1\n")
				wr.Flush()
				cmd.Process.Kill()
				return fmt.Errorf("wrong answer length: got %d tokens, want %d", len(parts)-1, n)
			}
			correct := true
			for i := 1; i <= n; i++ {
				var v int
				fmt.Sscanf(parts[i], "%d", &v)
				if v != perm[i] {
					correct = false
					break
				}
			}
			if correct {
				fmt.Fprintf(wr, "1\n")
				wr.Flush()
			} else {
				fmt.Fprintf(wr, "-1\n")
				wr.Flush()
				candIn.Close()
				cmd.Wait()
				return fmt.Errorf("wrong answer")
			}
			candIn.Close()
			cmd.Wait()
			return nil
		} else {
			// ignore unknown lines
			continue
		}
	}

	candIn.Close()
	cmd.Wait()
	return fmt.Errorf("candidate exited without answering")
}

func genPerm(rng *rand.Rand, n int) []int {
	// 1-indexed permutation with p[1] < p[2]
	a := rng.Perm(n)
	// a is 0-based values 0..n-1, shift to 1..n
	perm := make([]int, n+1)
	for i := 0; i < n; i++ {
		perm[i+1] = a[i] + 1
	}
	if perm[1] > perm[2] {
		perm[1], perm[2] = perm[2], perm[1]
	}
	return perm
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Test with small n values (20 is the minimum per problem constraints)
	sizes := []int{20, 20, 20, 25, 30, 20, 20, 20, 25, 30}

	for i, n := range sizes {
		perm := genPerm(rng, n)
		err := interactor(bin, n, perm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d (n=%d) failed: %v\n", i+1, n, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(sizes))
}
