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

func calc(x, pos int, opt []int) int {
	sum := 0
	for pos > 1 && x%2 == 0 && opt[pos-2] >= x/2 {
		sum += x
		x /= 2
		pos--
	}
	return sum + x
}

func expected(a []int) int {
	sort.Ints(a)
	var opt []int
	cnt := 1
	n := len(a)
	for i := 1; i < n; i++ {
		if a[i] != a[i-1] {
			opt = append(opt, cnt)
			cnt = 1
		} else {
			cnt++
		}
	}
	opt = append(opt, cnt)
	m := len(opt)
	if m == 1 {
		return n
	}
	sort.Ints(opt)
	ans := opt[m-1]
	if opt[m-1]%2 != 0 {
		opt[m-1]--
	}
	for opt[m-1] > 0 {
		if ans > opt[m-1]*2 {
			break
		}
		p := calc(opt[m-1], m, opt)
		if p > ans {
			ans = p
		}
		opt[m-1] -= 2
	}
	return ans
}

func runCase(exe string, input string, exp int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		return fmt.Errorf("invalid output %q", out.String())
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(200) + 1
		arr := make([]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(1000) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(arr)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
