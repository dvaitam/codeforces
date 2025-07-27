package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solveF(a []int, k int) int64 {
	n := len(a)
	sort.Ints(a)
	prefix := make([]int64, n+1)
	for i, v := range a {
		prefix[i+1] = prefix[i] + int64(v)
	}
	ans := int64(1 << 62)
	i := 0
	for i < n {
		j := i
		for j < n && a[j] == a[i] {
			j++
		}
		x := a[i]
		cnt := j - i
		if cnt >= k {
			return 0
		}
		need := k - cnt
		left := i
		right := n - j
		if left >= need {
			cost := int64(x)*int64(need) - (prefix[left] - prefix[left-need])
			if cost < ans {
				ans = cost
			}
		}
		if right >= need {
			cost := (prefix[j+need] - prefix[j]) - int64(x)*int64(need)
			if cost < ans {
				ans = cost
			}
		}
		if left+right >= need {
			if left < need {
				leftCost := int64(x)*int64(left) - prefix[left]
				rightTake := need - left
				rightCost := (prefix[j+rightTake] - prefix[j]) - int64(x)*int64(rightTake)
				if leftCost+rightCost < ans {
					ans = leftCost + rightCost
				}
			}
			if right < need {
				rightCost := prefix[n] - prefix[j] - int64(x)*int64(right)
				leftTake := need - right
				leftCost := int64(x)*int64(leftTake) - (prefix[left] - prefix[left-leftTake])
				if leftCost+rightCost < ans {
					ans = leftCost + rightCost
				}
			}
		}
		i = j
	}
	if ans < 0 {
		ans = 0
	}
	return ans
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
		expect := fmt.Sprintf("%d", solveF(append([]int(nil), arr...), k))
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", t, err, sb.String())
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", t, expect, out, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
