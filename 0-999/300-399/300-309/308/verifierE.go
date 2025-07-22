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

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Sheep struct {
	l, r, idx int
}

func okCheck(n, med int, print bool, v []Sheep, T [][]bool, buf *bytes.Buffer) bool {
	ogr := make([]int, n)
	cnt := make([]int, n)
	pos := make([]int, n)
	on := make([]int, n)
	for i := 0; i < n; i++ {
		ogr[i] = n - 1
		cnt[i] = 0
		pos[i] = -1
	}
	cnt[n-1] = n
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < n; j++ {
			sum += cnt[j]
			if sum > j+1 {
				return false
			}
			if sum == j+1 && j >= i {
				break
			}
		}
		cur := -1
		for j := 0; j < n; j++ {
			if ogr[j] <= sum-1 && pos[j] == -1 {
				cur = j
				break
			}
		}
		cnt[ogr[cur]]--
		cnt[i]++
		ogr[cur] = i
		pos[cur] = i
		on[i] = cur
		for j := 0; j < n; j++ {
			if T[cur][j] {
				if i+med < ogr[j] {
					cnt[ogr[j]]--
					ogr[j] = i + med
					cnt[ogr[j]]++
				}
			}
		}
	}
	if print {
		for i := 0; i < n; i++ {
			buf.WriteString(strconv.Itoa(v[on[i]].idx + 1))
			if i+1 < n {
				buf.WriteByte(' ')
			}
		}
	}
	return true
}

func expectedE(n int, intervals []Sheep) string {
	v := make([]Sheep, n)
	copy(v, intervals)
	sort.Slice(v, func(i, j int) bool { return v[i].r < v[j].r })
	T := make([][]bool, n)
	for i := 0; i < n; i++ {
		T[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if v[i].l > v[j].r || v[j].l > v[i].r {
				T[i][j] = false
			} else {
				T[i][j] = true
			}
		}
	}
	beg, end := 0, n
	for beg < end {
		med := (beg + end) / 2
		if okCheck(n, med, false, v, T, nil) {
			end = med
		} else {
			beg = med + 1
		}
	}
	var buf bytes.Buffer
	okCheck(n, beg, true, v, T, &buf)
	return strings.TrimSpace(buf.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	intervals := make([]Sheep, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(10)
		r := l + rng.Intn(10)
		intervals[i] = Sheep{l, r, i}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", intervals[i].l, intervals[i].r)
	}
	expect := expectedE(n, intervals)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
