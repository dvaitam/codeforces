package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input           string
	expectedProduct int
	expectedAssign  []int
}

func expected(arr []int) (int, []int) {
	n := len(arr) / 2
	cnt := make([]int, 90)
	for _, v := range arr {
		cnt[v-10]++
	}
	assignCnt := make([]int, 90)
	t, g1, g2 := 0, 0, 0
	for i := 0; i < 90; i++ {
		if cnt[i] == 1 {
			if t == 0 {
				assignCnt[i] = 1
				g1++
			} else {
				g2++
			}
			cnt[i] = 0
			t = (t + 1) % 2
		}
	}
	for i := 0; i < 90; i++ {
		if cnt[i] > 1 {
			assignCnt[i] = cnt[i] / 2
			g1++
			g2++
			if cnt[i]%2 == 1 {
				if t == 0 {
					assignCnt[i]++
				}
				t = (t + 1) % 2
			}
		}
	}
	assign := make([]int, 2*n)
	for i, v := range arr {
		idx := v - 10
		if assignCnt[idx] > 0 {
			assign[i] = 1
			assignCnt[idx]--
		} else {
			assign[i] = 2
		}
	}
	return g1 * g2, assign
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	arr := make([]int, 2*n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < 2*n; i++ {
		arr[i] = rng.Intn(90) + 10
		fmt.Fprintf(&sb, "%d ", arr[i])
	}
	sb.WriteByte('\n')
	prod, assign := expected(arr)
	return testCase{input: sb.String(), expectedProduct: prod, expectedAssign: assign}
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != 1+len(tc.expectedAssign) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\ninput:\n%s", i+1, 1+len(tc.expectedAssign), len(fields), tc.input)
			os.Exit(1)
		}
		var prod int
		if _, err := fmt.Sscan(fields[0], &prod); err != nil || prod != tc.expectedProduct {
			fmt.Fprintf(os.Stderr, "case %d failed: expected product %d got %s\ninput:\n%s", i+1, tc.expectedProduct, fields[0], tc.input)
			os.Exit(1)
		}
		for j, f := range fields[1:] {
			var v int
			if _, err := fmt.Sscan(f, &v); err != nil || v != tc.expectedAssign[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at position %d: expected %d got %s\ninput:\n%s", i+1, j+1, tc.expectedAssign[j], f, tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
