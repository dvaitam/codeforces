package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MOD = 1000000007

type testCaseD struct {
	input  string
	output string
}

func solveD(r io.Reader) string {
	reader := bufio.NewReader(r)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	a := make([]int, n+2)
	pred := make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] > 0 {
			pred[a[i]] = i
		}
	}
	taken := make([]bool, n+2)
	avail := make(map[int]struct{})
	freeUnfixed := make(map[int]struct{})
	for i := 1; i <= n; i++ {
		avail[i] = struct{}{}
		if a[i] == 0 {
			freeUnfixed[i] = struct{}{}
		}
	}
	ans := 1
	for t := 1; t <= n; t++ {
		if pos := pred[t]; pos > 0 {
			if _, ok := avail[pos]; !ok {
				return "0\n"
			}
			removeD(pos, n, taken, avail, freeUnfixed)
		} else {
			k := len(freeUnfixed)
			if k == 0 {
				return "0\n"
			}
			ans = int(int64(ans) * int64(k) % MOD)
			var pos int
			for idx := range freeUnfixed {
				pos = idx
				break
			}
			removeD(pos, n, taken, avail, freeUnfixed)
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func removeD(i, n int, taken []bool, avail, freeUnfixed map[int]struct{}) {
	taken[i] = true
	delete(avail, i)
	if _, ok := freeUnfixed[i]; ok {
		delete(freeUnfixed, i)
	}
	for _, j := range []int{i - 1, i + 1} {
		if j < 1 || j > n || taken[j] {
			continue
		}
		leftTaken := j-1 >= 1 && taken[j-1]
		rightTaken := j+1 <= n && taken[j+1]
		if leftTaken && rightTaken {
			delete(avail, j)
			delete(freeUnfixed, j)
		}
	}
}

func runCaseD(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(out.String()))
	}
	return nil
}

func genCaseD(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	perm := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	used := make([]bool, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			val := perm[i] + 1
			sb.WriteString(fmt.Sprint(val))
			used[val-1] = true
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCaseD(rng)
		expect := solveD(strings.NewReader(in))
		if err := runCaseD(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
