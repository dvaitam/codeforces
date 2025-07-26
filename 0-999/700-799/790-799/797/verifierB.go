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

func expectedB(arr []int) int {
	sum := 0
	const inf = int(1e9)
	minPosOdd := inf
	maxNegOdd := -inf
	for _, x := range arr {
		if x > 0 {
			sum += x
		}
		if x%2 != 0 {
			if x > 0 {
				if x < minPosOdd {
					minPosOdd = x
				}
			} else {
				if x > maxNegOdd {
					maxNegOdd = x
				}
			}
		}
	}
	if sum%2 == 1 {
		return sum
	}
	best := -inf
	if minPosOdd != inf {
		if s := sum - minPosOdd; s%2 != 0 && s > best {
			best = s
		}
	}
	if maxNegOdd != -inf {
		if s := sum + maxNegOdd; s%2 != 0 && s > best {
			best = s
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	hasOdd := false
	for i := 0; i < n; i++ {
		v := rng.Intn(21) - 10
		if v%2 != 0 {
			hasOdd = true
		}
		arr[i] = v
	}
	if !hasOdd {
		idx := rng.Intn(n)
		if arr[idx]%2 == 0 {
			arr[idx]++
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	exp := expectedB(arr)
	return sb.String(), exp
}

func runCase(bin string, input string, expect int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	resStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(resStr)
	if err != nil {
		return fmt.Errorf("bad output %q", resStr)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
