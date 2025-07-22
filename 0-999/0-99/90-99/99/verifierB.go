package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func solveB(r *bufio.Reader) string {
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return ""
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &arr[i])
	}
	if n == 1 {
		return "Exemplary pages."
	}
	allEq := true
	for i := 1; i < n; i++ {
		if arr[i] != arr[0] {
			allEq = false
			break
		}
	}
	if allEq {
		return "Exemplary pages."
	}
	cnt := make(map[int]int)
	for _, v := range arr {
		cnt[v]++
	}
	uniq := make([]int, 0, len(cnt))
	for v := range cnt {
		uniq = append(uniq, v)
	}
	sort.Ints(uniq)
	unrecover := func() string { return "Unrecoverable configuration." }
	switch len(uniq) {
	case 2:
		if n != 2 {
			return unrecover()
		}
		x, y := uniq[0], uniq[1]
		sum := x + y
		if sum%2 != 0 {
			return unrecover()
		}
		s := sum / 2
		v := y - s
		if v <= 0 {
			return unrecover()
		}
		ia, ib := 0, 0
		for i, val := range arr {
			if val == x {
				ia = i + 1
			} else if val == y {
				ib = i + 1
			}
		}
		return fmt.Sprintf("%d ml. from cup #%d to cup #%d.", v, ia, ib)
	case 3:
		x, y, z := uniq[0], uniq[1], uniq[2]
		if z-y != y-x || y-x <= 0 {
			return unrecover()
		}
		if cnt[x] != 1 || cnt[z] != 1 || cnt[y] != n-2 {
			return unrecover()
		}
		v := y - x
		ia, ib := 0, 0
		for i, val := range arr {
			if val == x {
				ia = i + 1
			} else if val == z {
				ib = i + 1
			}
		}
		return fmt.Sprintf("%d ml. from cup #%d to cup #%d.", v, ia, ib)
	case 1:
		return "Exemplary pages."
	default:
		return unrecover()
	}
}

func generateCaseB(rng *rand.Rand) string {
	choice := rng.Intn(5)
	switch choice {
	case 0:
		// n == 1
		v := rng.Intn(21)
		return fmt.Sprintf("1\n%d\n", v)
	case 1:
		// all equal
		n := rng.Intn(5) + 2
		val := rng.Intn(21)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
		return sb.String()
	case 2:
		// two unique, recoverable
		x := rng.Intn(21)
		v := rng.Intn(10) + 1
		y := x + 2*v
		return fmt.Sprintf("2\n%d %d\n", x, y)
	case 3:
		// three unique recoverable pattern
		n := rng.Intn(4) + 3
		v := rng.Intn(5) + 1
		s := rng.Intn(20) + v
		x := s - v
		y := s
		z := s + v
		arr := make([]int, n)
		arr[0] = x
		arr[1] = z
		for i := 2; i < n; i++ {
			arr[i] = y
		}
		// shuffle
		rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", arr[i])
		}
		sb.WriteByte('\n')
		return sb.String()
	default:
		// completely random
		n := rng.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(21))
		}
		sb.WriteByte('\n')
		return sb.String()
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		expect := solveB(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
