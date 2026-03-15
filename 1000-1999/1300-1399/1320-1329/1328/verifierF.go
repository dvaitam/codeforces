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
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// solveF is the correct embedded solver for CF 1328/F.
func solveF(input string) string {
	fields := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(fields[idx])
		idx++
		return v
	}

	n := nextInt()
	k := nextInt()

	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = nextInt()
	}

	sort.Ints(a)

	vals := make([]int64, 0, n)
	cnts := make([]int64, 0, n)
	for _, v := range a {
		x := int64(v)
		if len(vals) == 0 || vals[len(vals)-1] != x {
			vals = append(vals, x)
			cnts = append(cnts, 1)
		} else {
			cnts[len(cnts)-1]++
		}
	}

	m := len(vals)
	prefCnt := make([]int64, m)
	prefSum := make([]int64, m)
	for i := 0; i < m; i++ {
		prefCnt[i] = cnts[i]
		prefSum[i] = cnts[i] * vals[i]
		if i > 0 {
			prefCnt[i] += prefCnt[i-1]
			prefSum[i] += prefSum[i-1]
		}
	}

	totalCnt := prefCnt[m-1]
	totalSum := prefSum[m-1]
	needAll := int64(k)
	ans := int64(1 << 62)

	for i := 0; i < m; i++ {
		c := cnts[i]
		if c >= needAll {
			ans = 0
			break
		}

		t := needAll - c
		var lcnt, lsum int64
		if i > 0 {
			lcnt = prefCnt[i-1]
			lsum = prefSum[i-1]
		}
		rcnt := totalCnt - prefCnt[i]
		rsum := totalSum - prefSum[i]
		x := vals[i]

		bl := (x-1)*lcnt - lsum
		br := rsum - (x+1)*rcnt

		if lcnt >= t {
			cur := bl + t
			if cur < ans {
				ans = cur
			}
		}
		if rcnt >= t {
			cur := br + t
			if cur < ans {
				ans = cur
			}
		}

		cur := bl + br + t
		if cur < ans {
			ans = cur
		}
	}

	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(6)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(20) + 1
		k := rand.Intn(n) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(50) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveF(input)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", t, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", t, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
