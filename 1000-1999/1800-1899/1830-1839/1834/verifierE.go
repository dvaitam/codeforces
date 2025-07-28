package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const limitE = 300005

type testCaseE struct {
	arr []int64
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	g := gcd(a, b)
	res := a / g * b
	if res > int64(limitE) {
		return int64(limitE + 1)
	}
	return res
}

func expectedE(tc testCaseE) string {
	seen := make([]bool, limitE+2)
	prev := []int64{}
	for _, val := range tc.arr {
		curMap := make(map[int64]struct{})
		if val <= limitE {
			curMap[val] = struct{}{}
		}
		for _, v := range prev {
			l := lcm(v, val)
			if l <= limitE {
				curMap[l] = struct{}{}
			}
		}
		prev = prev[:0]
		for k := range curMap {
			prev = append(prev, k)
			seen[int(k)] = true
		}
	}
	ans := 1
	for ans <= limitE && seen[ans] {
		ans++
	}
	return fmt.Sprint(ans)
}

func genTestsE() []testCaseE {
	rand.Seed(5)
	tests := make([]testCaseE, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(5) + 1
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = int64(rand.Intn(20) + 1)
		}
		tests = append(tests, testCaseE{arr: arr})
	}
	return tests
}

func runCase(bin string, tc testCaseE) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", len(tc.arr)))
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedE(tc)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
