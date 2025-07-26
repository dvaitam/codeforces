package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod int64 = 1000000007

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveCase(arr []int) int64 {
	n := len(arr)
	var ans int64
	for mask := 1; mask < (1 << n); mask++ {
		g := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				if g == 0 {
					g = arr[i]
				} else {
					g = gcd(g, arr[i])
				}
			}
		}
		if g > 1 {
			k := bits.OnesCount(uint(mask))
			ans += int64(k) * int64(g)
			ans %= mod
		}
	}
	return ans % mod
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(8) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(20) + 1
	}
	return arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		arr := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(arr))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := solveCase(arr)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Printf("case %d: bad integer output %q\n", i+1, out)
			return
		}
		if got%mod != want%mod {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %d\ngot: %d\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
