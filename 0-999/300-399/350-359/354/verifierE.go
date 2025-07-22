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

const (
	PERMS  = 28
	DIGITS = 19
)

var p [PERMS][7]int
var bucket [10][]int
var pow10 [DIGITS]int64
var a [6]int64

func nextDigit(i int) int {
	if i == 0 {
		return 4
	} else if i == 4 {
		return 7
	}
	return 8
}

func initPrecomp() {
	pid := 0
	for i := 0; i <= 7; i = nextDigit(i) {
		for j := i; j <= 7; j = nextDigit(j) {
			for k := j; k <= 7; k = nextDigit(k) {
				for l := k; l <= 7; l = nextDigit(l) {
					for m := l; m <= 7; m = nextDigit(m) {
						for n2 := m; n2 <= 7; n2 = nextDigit(n2) {
							sum := i + j + k + l + m + n2
							p[pid][0], p[pid][1], p[pid][2] = i, j, k
							p[pid][3], p[pid][4], p[pid][5] = l, m, n2
							p[pid][6] = sum
							bucket[sum%10] = append(bucket[sum%10], pid)
							pid++
						}
					}
				}
			}
		}
	}
	pow10[0] = 1
	for i := 1; i < DIGITS; i++ {
		pow10[i] = pow10[i-1] * 10
	}
}

func search(n int64, dig int) bool {
	if n == 0 {
		return true
	}
	if dig >= DIGITS {
		return false
	}
	d := int((n / pow10[dig]) % 10)
	for _, idx := range bucket[d] {
		sumDigit := int64(p[idx][6])
		if n < sumDigit*pow10[dig] {
			continue
		}
		for j := 0; j < 6; j++ {
			a[j] += int64(p[idx][j]) * pow10[dig]
		}
		if search(n-sumDigit*pow10[dig], dig+1) {
			return true
		}
		for j := 0; j < 6; j++ {
			a[j] -= int64(p[idx][j]) * pow10[dig]
		}
	}
	return false
}

func expected(nums []int64) []string {
	res := make([]string, len(nums))
	for idx, n := range nums {
		for i := 0; i < 6; i++ {
			a[i] = 0
		}
		if search(n, 0) {
			parts := make([]string, 6)
			for i := 0; i < 6; i++ {
				parts[i] = strconv.FormatInt(a[i], 10)
			}
			res[idx] = strings.Join(parts, " ")
		} else {
			res[idx] = "-1"
		}
	}
	return res
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	initPrecomp()
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		t := rng.Intn(5) + 1
		nums := make([]int64, t)
		for j := 0; j < t; j++ {
			nums[j] = rng.Int63n(1e12) + 1
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", t))
		for _, v := range nums {
			input.WriteString(fmt.Sprintf("%d\n", v))
		}
		expLines := expected(nums)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(gotLines) != len(expLines) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(expLines), len(gotLines), input.String())
			os.Exit(1)
		}
		for j, line := range expLines {
			if strings.TrimSpace(gotLines[j]) != line {
				fmt.Fprintf(os.Stderr, "case %d failed on line %d: expected %s got %s\ninput:\n%s", i+1, j+1, line, gotLines[j], input.String())
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
