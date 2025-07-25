package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n, k  int
	piles []int
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	k := rng.Intn(8) + 1
	piles := make([]int, n)
	for i := range piles {
		piles[i] = rng.Intn(20) + 1
	}
	return testCase{n: n, k: k, piles: piles}
}

func possible(tc testCase) bool {
	mini, maxi := tc.piles[0], tc.piles[0]
	for _, v := range tc.piles {
		if v < mini {
			mini = v
		}
		if v > maxi {
			maxi = v
		}
	}
	return maxi-mini <= tc.k
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i, v := range tc.piles {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := strings.TrimSpace(scanner.Text())
	poss := possible(tc)
	if first == "NO" {
		if poss {
			return fmt.Errorf("should be YES")
		}
		if scanner.Scan() {
			return fmt.Errorf("extra output after NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("expected YES or NO")
	}
	if !poss {
		return fmt.Errorf("should be NO")
	}
	counts := make([][]int, tc.n)
	for i := 0; i < tc.n; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("not enough lines")
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) != tc.piles[i] {
			return fmt.Errorf("pile %d expected %d numbers got %d", i+1, tc.piles[i], len(fields))
		}
		cnt := make([]int, tc.k)
		for _, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil {
				return fmt.Errorf("bad number: %v", err)
			}
			if val < 1 || val > tc.k {
				return fmt.Errorf("color out of range")
			}
			cnt[val-1]++
		}
		counts[i] = cnt
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	for c := 0; c < tc.k; c++ {
		minC, maxC := counts[0][c], counts[0][c]
		for i := 1; i < tc.n; i++ {
			if counts[i][c] < minC {
				minC = counts[i][c]
			}
			if counts[i][c] > maxC {
				maxC = counts[i][c]
			}
		}
		if maxC-minC > 1 {
			return fmt.Errorf("color difference >1 for color %d", c+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
