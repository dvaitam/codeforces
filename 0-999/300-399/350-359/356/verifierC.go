package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveC(in string) string {
	reader := bufio.NewReader(strings.NewReader(in))
	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	cnt := make([]int64, 5)
	var ai int
	var S int64
	for i := int64(0); i < n; i++ {
		fmt.Fscan(reader, &ai)
		cnt[ai]++
		S += int64(ai)
	}
	cost3 := []int64{3, 2, 1, 0, 0}
	cost4 := []int64{4, 3, 2, 1, 0}
	var best int64 = -1
	c4min := int64(0)
	if S > 3*n {
		c4min = S - 3*n
	}
	c4max := S / 4
	if c4max > n {
		c4max = n
	}
	for c4 := c4min; c4 <= c4max; c4++ {
		remS := S - 4*c4
		if remS < 0 || remS%3 != 0 {
			continue
		}
		c3 := remS / 3
		if c3 < 0 || c4+c3 > n {
			continue
		}
		need4 := c4
		used4 := make([]int64, 5)
		for ai := 4; ai >= 0 && need4 > 0; ai-- {
			take := cnt[ai]
			if take > need4 {
				take = need4
			}
			used4[ai] = take
			need4 -= take
		}
		if need4 > 0 {
			continue
		}
		need3 := c3
		used3 := make([]int64, 5)
		for _, ai := range []int{3, 4, 2, 1, 0} {
			if need3 <= 0 {
				break
			}
			avail := cnt[ai] - used4[ai]
			if avail <= 0 {
				continue
			}
			take := avail
			if take > need3 {
				take = need3
			}
			used3[ai] = take
			need3 -= take
		}
		if need3 > 0 {
			continue
		}
		var cost int64
		for ai := 0; ai <= 4; ai++ {
			cost += used4[ai] * cost4[ai]
			cost += used3[ai] * cost3[ai]
		}
		if best < 0 || cost < best {
			best = cost
		}
	}
	return fmt.Sprintln(best)
}

func genTest(r *rand.Rand) string {
	n := r.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	cntPos := 0
	for i := 0; i < n; i++ {
		val := r.Intn(5)
		if val > 0 {
			cntPos++
		}
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(val))
	}
	if cntPos == 0 {
		sb.Reset()
		sb.WriteString("1\n1")
	} else {
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runBinary(path, in string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierC <path-to-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(3))
	const tests = 100
	for i := 0; i < tests; i++ {
		in := genTest(r)
		expect := strings.TrimSpace(solveC(in))
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
