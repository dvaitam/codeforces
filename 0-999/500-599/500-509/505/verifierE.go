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

func expected(n, m, k int, p int64, hArr, aArr []int64) int64 {
	mi := int64(m)
	var maxA int64
	var upper int64

	fArr := make([]int64, n)
	gArr := make([]int64, n)

	for i := 0; i < n; i++ {
		g := mi * aArr[i]
		f := hArr[i] + g
		gArr[i] = g
		fArr[i] = f
		if aArr[i] > maxA {
			maxA = aArr[i]
		}
		if f > upper {
			upper = f
		}
	}

	ceilPos := func(x, y int64) int64 {
		if x <= 0 {
			return 0
		}
		return (x + y - 1) / y
	}

	totalCap := int64(m * k)

	check := func(H int64) bool {
		if H < maxA {
			return false
		}
		cnt := make([]int64, m+1)
		var total int64

		for i := 0; i < n; i++ {
			t := ceilPos(fArr[i]-H, p)
			total += t
			if total > totalCap {
				return false
			}

			u := ceilPos(gArr[i]-H, p)
			a := aArr[i]

			for r := int64(0); r < u; r++ {
				d := int((H + r*p) / a)
				cnt[d]++
			}
			cnt[m] += t - u
		}

		var used int64
		kk := int64(k)
		for s := 1; s <= m; s++ {
			used += cnt[s]
			if used > int64(s)*kk {
				return false
			}
		}
		return true
	}

	l, r := maxA-1, upper
	for l+1 < r {
		mid := (l + r) >> 1
		if check(mid) {
			r = mid
		} else {
			l = mid
		}
	}
	return r
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		k := rng.Intn(4) + 1
		p := int64(rng.Intn(20) + 1)
		hArr := make([]int64, n)
		aArr := make([]int64, n)
		for i := 0; i < n; i++ {
			hArr[i] = int64(rng.Intn(20))
			aArr[i] = int64(rng.Intn(10) + 1)
		}
		input := fmt.Sprintf("%d %d %d %d\n", n, m, k, p)
		for i := 0; i < n; i++ {
			input += fmt.Sprintf("%d %d\n", hArr[i], aArr[i])
		}
		exp := expected(n, m, k, p, hArr, aArr)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", tc+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
