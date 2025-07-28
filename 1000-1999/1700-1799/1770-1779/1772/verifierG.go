package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solveG(n int, x, y int64, a []int64) int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		v := a[i-1] - int64(i-1)
		if v > pref[i-1] {
			pref[i] = v
		} else {
			pref[i] = pref[i-1]
		}
	}
	if x < pref[1] {
		return -1
	}
	r := x
	var games int64
	for r < y {
		k := sort.Search(len(pref), func(i int) bool { return pref[i] > r }) - 1
		if r+int64(k) >= y {
			games += y - r
			break
		}
		delta := int64(2*k - n)
		if delta <= 0 {
			return -1
		}
		if k == n {
			games += y - r
			break
		}
		nextR := pref[k+1]
		finish := y - int64(k)
		target := nextR
		if finish < target {
			target = finish
		}
		cycles := (target - r + delta - 1) / delta
		if cycles <= 0 {
			cycles = 1
		}
		r += cycles * delta
		games += cycles * int64(n)
	}
	return games
}

func expectedG(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var t int
	fmt.Sscanf(lines[0], "%d", &t)
	idx := 1
	var out strings.Builder
	for i := 0; i < t; i++ {
		var n int
		var x, y int64
		fmt.Sscanf(lines[idx], "%d %d %d", &n, &x, &y)
		idx++
		fields := strings.Fields(lines[idx])
		idx++
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Sscanf(fields[j], "%d", &arr[j])
		}
		out.WriteString(fmt.Sprintf("%d", solveG(n, x, y, arr)))
		if i+1 < t {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func genTestsG() []string {
	rand.Seed(7)
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		t := rand.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for i := 0; i < t; i++ {
			n := rand.Intn(8) + 1
			x := int64(rand.Intn(20) + 1)
			y := x + int64(rand.Intn(20)+1)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, y))
			for j := 0; j < n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", rand.Intn(20)+1))
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierG.go <binary>\n")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsG()
	for i, tcase := range tests {
		want := expectedG(tcase)
		got, err := runBinary(bin, tcase)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("Test %d failed.\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, tcase, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
