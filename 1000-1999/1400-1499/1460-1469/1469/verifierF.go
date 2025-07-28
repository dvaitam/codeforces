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

func run(bin, input string) (string, error) {
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

func can(R int64, lengths []int64, k int64) bool {
	if R == 0 {
		return k <= 1
	}
	n := len(lengths)
	arr := make([]int64, n)
	copy(arr, lengths)
	sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
	cnt := make([]int64, R+1)
	cnt[0] = 1
	total := int64(1)
	idx := int64(0)
	for _, L := range arr {
		for idx < R {
			if cnt[idx] == 0 {
				idx++
				continue
			}
			m := R - idx - 1
			if m <= 0 {
				idx++
				continue
			}
			break
		}
		if idx >= R || total >= k {
			break
		}
		d := idx
		cnt[d]--
		total--
		m := R - d - 1
		if L <= 2*m+1 {
			x := (L - 1) / 2
			for i := int64(1); i <= x && d+1+i <= R; i++ {
				cnt[d+1+i] += 2
			}
			if L%2 == 0 && d+1+L/2 <= R {
				cnt[d+1+L/2]++
			}
			total += L - 1
		} else {
			for i := int64(1); i <= m && d+1+i <= R; i++ {
				cnt[d+1+i] += 2
			}
			total += 2 * m
		}
		if total >= k {
			return true
		}
	}
	return total >= k
}

func expected(lengths []int64, k int64) int64 {
	sum := int64(1)
	for _, L := range lengths {
		sum += L - 2
	}
	if sum < k {
		return -1
	}
	l, r := int64(1), int64(200000)
	for l < r {
		mid := (l + r) / 2
		if can(mid, lengths, k) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

func generateCase(rng *rand.Rand) ([]int64, int64) {
	n := rng.Intn(8) + 1
	lengths := make([]int64, n)
	for i := 0; i < n; i++ {
		lengths[i] = int64(rng.Intn(8) + 3)
	}
	var sum int64 = 1
	for _, L := range lengths {
		sum += L - 2
	}
	k := int64(rng.Intn(int(sum+5)) + 2)
	return lengths, k
}

func runCase(bin string, lengths []int64, k int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(lengths), k))
	for i, v := range lengths {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := expected(lengths, k)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		lengths, k := generateCase(rng)
		if err := runCase(bin, lengths, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
