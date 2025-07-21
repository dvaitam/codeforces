package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveB3(a []int) string {
	n := len(a)
	type node struct{ fs, sn int }
	arr := make([]node, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = node{a[i-1], i}
	}
	sort.Slice(arr[1:], func(i, j int) bool { return arr[i+1].fs < arr[j+1].fs })
	x := 0
	for i := 1; i < n; i++ {
		if arr[i].fs == arr[i+1].fs {
			x = i
			break
		}
	}
	ansX := make([]int, n+1)
	ansY := make([]int, n+1)
	anss := make([]int, n+1)
	v := make([]int, n+1)
	if x > 0 || arr[1].fs == 0 {
		if x > 0 {
			arr[1], arr[x] = arr[x], arr[1]
			arr[2], arr[x+1] = arr[x+1], arr[2]
		}
		ansX[arr[1].sn] = 1
		ansY[arr[1].sn] = 1
		v[1] = 1
		if arr[1].fs == 0 {
			anss[arr[1].sn] = arr[1].sn
		} else {
			anss[arr[1].sn] = arr[2].sn
		}
		for i := 2; i <= n; i++ {
			idx := arr[i].sn
			if arr[i].fs == 0 {
				ansX[idx] = i
				ansY[idx] = 1
				v[i] = 1
				anss[idx] = arr[i].sn
			} else if arr[i].fs < i {
				ansX[idx] = i
				y := v[i-arr[i].fs]
				ansY[idx] = y
				v[i] = y
				anss[idx] = arr[i-arr[i].fs].sn
			} else {
				ansX[idx] = i
				y := arr[i].fs - i + 2
				ansY[idx] = y
				v[i] = y
				anss[idx] = arr[1].sn
			}
		}
	} else {
		if n == 2 {
			return "NO\n"
		}
		ansX[arr[n].sn] = 1
		ansY[arr[n].sn] = 1
		anss[arr[n].sn] = arr[n-1].sn
		ansX[arr[n-1].sn] = n
		ansY[arr[n-1].sn] = 2
		anss[arr[n-1].sn] = arr[1].sn
		for i := 1; i <= n-2; i++ {
			idx := arr[i].sn
			ansX[idx] = i + 1
			ansY[idx] = 1
			anss[idx] = arr[n].sn
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", ansX[i], ansY[i]))
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", anss[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(10) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n + 1)
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
	return sb.String(), arr
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB3.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, arr := generateCase(rng)
		expect := solveB3(arr)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
