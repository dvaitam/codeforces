package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const maxBits = 31

func solve(n int, L int64, costs []int64) int64 {
	const inf int64 = 1 << 60
	arr := make([]int64, maxBits)
	copy(arr, costs)
	for i := n; i < maxBits; i++ {
		arr[i] = inf
	}
	for i := 1; i < maxBits; i++ {
		if arr[i] > 2*arr[i-1] {
			arr[i] = 2 * arr[i-1]
		}
	}
	ans := inf
	spent := int64(0)
	for i := maxBits - 1; i >= 0; i-- {
		size := int64(1) << uint(i)
		need := L / size
		spent += need * arr[i]
		L -= need * size
		if L > 0 {
			if cand := spent + arr[i]; cand < ans {
				ans = cand
			}
		} else if spent < ans {
			ans = spent
		}
	}
	return ans
}

func randomInput() (int, int64, []int64) {
	n := rand.Intn(5) + 1
	L := rand.Int63n(1e6) + 1
	costs := make([]int64, n)
	for i := range costs {
		costs[i] = rand.Int63n(1e6) + 1
	}
	return n, L, costs
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	const cases = 100
	for i := 0; i < cases; i++ {
		n, L, costs := randomInput()
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, L))
		for j, c := range costs {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", c))
		}
		sb.WriteByte('\n')
		input := sb.String()
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			fmt.Printf("program output:\n%s\n", string(out))
			return
		}
		got := strings.TrimSpace(string(out))
		want := fmt.Sprintf("%d", solve(n, L, costs))
		if got != want {
			fmt.Printf("case %d failed:\ninput:\n%sexpected %s got %s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Printf("OK %d cases\n", cases)
}
