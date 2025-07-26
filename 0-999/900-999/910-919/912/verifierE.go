package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const MAX int64 = 1e18

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generate(pr []int64, idx int, cur int64, res *[]int64) {
	if idx == len(pr) {
		*res = append(*res, cur)
		return
	}
	p := pr[idx]
	val := cur
	for val <= MAX {
		generate(pr, idx+1, val, res)
		if val > MAX/p {
			break
		}
		val *= p
	}
}

func countLeq(a, b []int64, x int64) int64 {
	j := len(b) - 1
	var cnt int64
	for _, v := range a {
		for j >= 0 && v > x/b[j] {
			j--
		}
		if j < 0 {
			break
		}
		cnt += int64(j + 1)
	}
	return cnt
}

func solveCase(primes []int64, k int64) string {
	n := len(primes)
	m := n / 2
	leftPr := primes[:m]
	rightPr := primes[m:]
	left := make([]int64, 0)
	right := make([]int64, 0)
	generate(leftPr, 0, 1, &left)
	generate(rightPr, 0, 1, &right)
	sort.Slice(left, func(i, j int) bool { return left[i] < left[j] })
	sort.Slice(right, func(i, j int) bool { return right[i] < right[j] })

	l, r := int64(1), MAX
	for l < r {
		mid := (l + r) / 2
		if countLeq(left, right, mid) >= k {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return fmt.Sprint(l)
}

var primeList = []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

func generateCase(rng *rand.Rand) ([]int64, int64) {
	n := rng.Intn(6) + 1
	primes := make([]int64, 0, n)
	perm := rng.Perm(len(primeList))
	for i := 0; i < n; i++ {
		primes = append(primes, primeList[perm[i]])
	}
	sort.Slice(primes, func(i, j int) bool { return primes[i] < primes[j] })
	k := int64(rng.Intn(1000) + 1)
	return primes, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		primes, k := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(primes)))
		for j, p := range primes {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(p))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", k))
		input := sb.String()
		expect := solveCase(append([]int64(nil), primes...), k)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
