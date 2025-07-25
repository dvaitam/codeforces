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

func solveF(nums []int) string {
	maxVal := 0
	for _, x := range nums {
		if x > maxVal {
			maxVal = x
		}
	}
	freq := make([]int, maxVal+1)
	for _, x := range nums {
		if x >= 0 {
			freq[x]++
		}
	}
	pref := make([]int, maxVal+1)
	for i := 1; i <= maxVal; i++ {
		pref[i] = pref[i-1] + freq[i]
	}
	var ans int64
	for v := 1; v <= maxVal; v++ {
		if freq[v] == 0 {
			continue
		}
		var total int64
		for j := v; j <= maxVal; j += v {
			r := j + v - 1
			if r > maxVal {
				r = maxVal
			}
			cnt := pref[r] - pref[j-1]
			total += int64(cnt) * int64(j)
		}
		if total > ans {
			ans = total
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func genCaseF(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := range arr {
		arr[i] = rng.Intn(50) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String(), solveF(arr)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseF(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
