package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func maxRem(arr []int) int {
	best := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] == i+1 {
			tmp := make([]int, len(arr)-1)
			copy(tmp, arr[:i])
			copy(tmp[i:], arr[i+1:])
			if val := 1 + maxRem(tmp); val > best {
				best = val
			}
		}
	}
	return best
}

func solve(reader *bufio.Reader) string {
	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return ""
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var sb strings.Builder
	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		b := append([]int(nil), a...)
		for i := 0; i < x; i++ {
			b[i] = n + 1
		}
		for i := 0; i < y; i++ {
			b[len(b)-1-i] = n + 1
		}
		sb.WriteString(fmt.Sprintf("%d\n", maxRem(b)))
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(7) + 1
	q := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(n)+1)
	}
	sb.WriteByte('\n')
	for j := 0; j < q; j++ {
		x := rng.Intn(n)
		y := rng.Intn(n - x)
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
