package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseC struct {
	arr []int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func expectedProduct(a []int) *big.Int {
	zeros := 0
	negatives := []int{}
	subset := []int{}
	for _, v := range a {
		if v == 0 {
			zeros++
		} else {
			subset = append(subset, v)
			if v < 0 {
				negatives = append(negatives, v)
			}
		}
	}
	if zeros == len(a) {
		return big.NewInt(0)
	}
	if len(negatives)%2 != 0 {
		idx := -1
		minAbs := 0
		for i, v := range subset {
			if v < 0 {
				if idx == -1 || abs(v) < minAbs {
					idx = i
					minAbs = abs(v)
				}
			}
		}
		if idx >= 0 {
			subset = append(subset[:idx], subset[idx+1:]...)
		}
		if len(subset) == 0 {
			return big.NewInt(0)
		}
	}
	prod := big.NewInt(1)
	for _, v := range subset {
		prod.Mul(prod, big.NewInt(int64(v)))
	}
	return prod
}

func applyOps(a []int, ops []string) (*big.Int, error) {
	n := len(a)
	vals := make([]*big.Int, n)
	alive := make([]bool, n)
	for i, v := range a {
		vals[i] = big.NewInt(int64(v))
		alive[i] = true
	}
	removeUsed := 0
	for _, op := range ops {
		parts := strings.Fields(op)
		if len(parts) == 0 {
			continue
		}
		if parts[0] == "1" {
			if len(parts) != 3 {
				return nil, fmt.Errorf("bad op %q", op)
			}
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			x--
			y--
			if x < 0 || x >= n || y < 0 || y >= n {
				return nil, fmt.Errorf("index out of range")
			}
			if !alive[x] || !alive[y] {
				return nil, fmt.Errorf("using removed index")
			}
			vals[y].Mul(vals[y], vals[x])
			alive[x] = false
			removeUsed++
		} else if parts[0] == "2" {
			if len(parts) != 2 {
				return nil, fmt.Errorf("bad op %q", op)
			}
			x, _ := strconv.Atoi(parts[1])
			x--
			if x < 0 || x >= n {
				return nil, fmt.Errorf("index out of range")
			}
			if !alive[x] {
				return nil, fmt.Errorf("removing already removed")
			}
			alive[x] = false
			removeUsed++
		} else {
			return nil, fmt.Errorf("bad op %q", op)
		}
	}
	if removeUsed != n-1 {
		return nil, fmt.Errorf("wrong number of operations")
	}
	var res *big.Int
	count := 0
	for i := 0; i < n; i++ {
		if alive[i] {
			res = vals[i]
			count++
		}
	}
	if count != 1 {
		return nil, fmt.Errorf("invalid final state")
	}
	return new(big.Int).Set(res), nil
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCaseC {
	rng := rand.New(rand.NewSource(44))
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rng.Intn(8) + 2
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(11) - 5
		}
		tests[i] = testCaseC{arr}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for idx, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(tc.arr))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		expected := expectedProduct(tc.arr)
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", idx+1, err)
			return
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		prod, err := applyOps(tc.arr, lines)
		if err != nil {
			fmt.Printf("test %d: invalid operations: %v\n", idx+1, err)
			return
		}
		if prod.Cmp(expected) != 0 {
			fmt.Printf("test %d failed:\ninput:%sexpected %s got %s\n", idx+1, sb.String(), expected.String(), prod.String())
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
