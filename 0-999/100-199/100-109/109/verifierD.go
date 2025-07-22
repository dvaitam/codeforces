package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func lucky(x int) bool {
	if x == 0 {
		return true
	}
	for x > 0 {
		d := x % 10
		if d != 4 && d != 7 {
			return false
		}
		x /= 10
	}
	return true
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(1000) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInput(input string) []int {
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(r, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &arr[i])
	}
	return arr
}

func hasLucky(arr []int) bool {
	for _, v := range arr {
		if lucky(v) {
			return true
		}
	}
	return false
}

func runCase(bin string, input string) error {
	original := parseInput(input)
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	reader := bufio.NewReader(strings.NewReader(out.String()))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return fmt.Errorf("cannot read k: %v", err)
	}
	arr := append([]int(nil), original...)
	if k == -1 {
		if hasLucky(arr) || !sort.IntsAreSorted(append([]int(nil), arr...)) {
			return fmt.Errorf("unexpected -1 output")
		}
		if _, err := fmt.Fscan(reader, new(int)); err != io.EOF {
			return fmt.Errorf("extra output")
		}
		return nil
	}
	n := len(arr)
	if k < 0 || k > 2*n {
		return fmt.Errorf("invalid k")
	}
	for i := 0; i < k; i++ {
		var a, b int
		if _, err := fmt.Fscan(reader, &a, &b); err != nil {
			return fmt.Errorf("cannot read op %d: %v", i+1, err)
		}
		if a < 1 || a > n || b < 1 || b > n || a == b {
			return fmt.Errorf("invalid op %d", i+1)
		}
		a--
		b--
		if !lucky(arr[a]) && !lucky(arr[b]) {
			return fmt.Errorf("swap %d violates lucky condition", i+1)
		}
		arr[a], arr[b] = arr[b], arr[a]
	}
	if _, err := fmt.Fscan(reader, new(int)); err != io.EOF {
		return fmt.Errorf("extra output")
	}
	sorted := append([]int(nil), original...)
	sort.Ints(sorted)
	for i := 0; i < n; i++ {
		if arr[i] != sorted[i] {
			return fmt.Errorf("array not sorted after swaps")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{
		"1\n4\n",
		"2\n1 2\n",
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
