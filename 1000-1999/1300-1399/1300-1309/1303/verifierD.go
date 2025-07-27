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

func solveCase(n int64, boxes []int64) string {
	const maxb = 61
	cnt := make([]int64, maxb)
	var sum int64
	for _, a := range boxes {
		b := 0
		for (1 << uint(b)) != a {
			b++
		}
		cnt[b]++
		sum += a
	}
	if sum < n {
		return "-1\n"
	}
	var ans int64
	var have int64
	for i := 0; i < maxb; i++ {
		have += cnt[i]
		if (n>>uint(i))&1 == 1 {
			if have > 0 {
				have--
			} else {
				j := i + 1
				for j < maxb && cnt[j] == 0 {
					j++
				}
				for k := j; k > i; k-- {
					cnt[k]--
					cnt[k-1] += 2
					ans++
				}
				have += cnt[i]
				have--
			}
		}
		have /= 2
	}
	return fmt.Sprintf("%d\n", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1_000_000) + 1
	m := rng.Intn(8) + 1
	boxes := make([]int64, m)
	for i := 0; i < m; i++ {
		b := rng.Intn(20)
		boxes[i] = 1 << uint(b)
	}
	inSb := strings.Builder{}
	fmt.Fprintf(&inSb, "1\n%d %d\n", n, m)
	for i, a := range boxes {
		if i > 0 {
			inSb.WriteByte(' ')
		}
		fmt.Fprintf(&inSb, "%d", a)
	}
	inSb.WriteByte('\n')
	out := solveCase(n, boxes)
	return inSb.String(), out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
