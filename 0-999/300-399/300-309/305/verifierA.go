package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Embedded correct solver for 305A.
// For any pair to be summable, every decimal digit position must have at least
// one of the two numbers with a 0 in that position.
// Numbers 0 and 100 are always safe (all-zero in the relevant digits).
// Among others, we can pick at most one "single-digit" (1-9) and at most one
// "tens-multiple" (10,20,...,90). If we have both, they pair fine.
// Any other number (like 23) can only coexist with 0 and 100, so we pick at most one
// such number and only if we have no single-digit and no tens-multiple.
func solveA(input string) string {
	fields := strings.Fields(input)
	k, _ := strconv.Atoi(fields[0])
	nums := make([]int, k)
	for i := 0; i < k; i++ {
		nums[i], _ = strconv.Atoi(fields[i+1])
	}

	has0 := false
	has100 := false
	single := -1
	ten := -1
	other := -1

	for _, x := range nums {
		if x == 0 {
			has0 = true
		} else if x == 100 {
			has100 = true
		} else if x >= 1 && x <= 9 {
			single = x
		} else if x >= 10 && x <= 90 && x%10 == 0 {
			ten = x
		} else {
			other = x
		}
	}

	var ans []int
	if has0 {
		ans = append(ans, 0)
	}
	if has100 {
		ans = append(ans, 100)
	}

	if single != -1 && ten != -1 {
		ans = append(ans, single, ten)
	} else if single != -1 {
		ans = append(ans, single)
	} else if ten != -1 {
		ans = append(ans, ten)
	} else if other != -1 {
		ans = append(ans, other)
	}

	sort.Ints(ans)
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(ans))
	for i, v := range ans {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, v)
	}
	fmt.Fprintln(&buf)
	return strings.TrimSpace(buf.String())
}

// Validate that a set of numbers is "compatible": every pair can be summed.
func isValidSet(nums []int) bool {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if !canSum(nums[i], nums[j]) {
				return false
			}
		}
	}
	return true
}

func canSum(a, b int) bool {
	for a > 0 || b > 0 {
		if a%10 != 0 && b%10 != 0 {
			return false
		}
		a /= 10
		b /= 10
	}
	return true
}

func genTestA(rng *rand.Rand) string {
	k := rng.Intn(100) + 1
	used := make(map[int]bool)
	arr := make([]int, 0, k)
	for len(arr) < k {
		x := rng.Intn(101)
		if !used[x] {
			used[x] = true
			arr = append(arr, x)
		}
	}
	var buf bytes.Buffer
	fmt.Fprintln(&buf, k)
	for i, v := range arr {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, v)
	}
	buf.WriteByte('\n')
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 200; i++ {
		in := genTestA(rng)
		expected := solveA(in)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", i, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(stdout.String())

		// Parse candidate output and validate.
		gotFields := strings.Fields(got)
		expFields := strings.Fields(expected)

		if len(gotFields) == 0 || len(expFields) == 0 {
			fmt.Fprintf(os.Stderr, "test %d: empty output\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expected, got)
			os.Exit(1)
		}

		gotN, _ := strconv.Atoi(gotFields[0])
		expN, _ := strconv.Atoi(expFields[0])

		if gotN != expN {
			fmt.Fprintf(os.Stderr, "test %d: size mismatch expected %d got %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, expN, gotN, in, expected, got)
			os.Exit(1)
		}

		if len(gotFields) != gotN+1 {
			fmt.Fprintf(os.Stderr, "test %d: wrong number of elements\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expected, got)
			os.Exit(1)
		}

		// Parse the input set to verify candidate output is a subset
		inFields := strings.Fields(in)
		inK, _ := strconv.Atoi(inFields[0])
		inputSet := make(map[int]bool)
		for j := 1; j <= inK; j++ {
			v, _ := strconv.Atoi(inFields[j])
			inputSet[v] = true
		}

		gotNums := make([]int, gotN)
		for j := 0; j < gotN; j++ {
			v, _ := strconv.Atoi(gotFields[j+1])
			gotNums[j] = v
			if !inputSet[v] {
				fmt.Fprintf(os.Stderr, "test %d: candidate output %d not in input\ninput:\n%sgot:\n%s\n", i, v, in, got)
				os.Exit(1)
			}
		}

		if !isValidSet(gotNums) {
			fmt.Fprintf(os.Stderr, "test %d: candidate output is not a valid set\ninput:\n%sgot:\n%s\n", i, in, got)
			os.Exit(1)
		}

		if gotN != expN {
			fmt.Fprintf(os.Stderr, "test %d: not optimal: expected size %d got %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, expN, gotN, in, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
