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

type testCaseD struct {
	n       int
	sectors [][]int
}

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

func generateCase() testCaseD {
	n := rand.Intn(4) + 3 // 3..6
	sectors := make([][]int, n)
	for i := 0; i < n; i++ {
		k := rand.Intn(3) + 1
		start := i*1000 + 1
		arr := make([]int, k)
		for j := 0; j < k; j++ {
			arr[j] = start + rand.Intn(900)
		}
		sort.Ints(arr)
		sectors[i] = arr
	}
	return testCaseD{n: n, sectors: sectors}
}

func compute(tc testCaseD) int {
	n := tc.n
	sectors := make([][]int, n)
	for i := 0; i < n; i++ {
		sec := append([]int(nil), tc.sectors[i]...)
		sort.Ints(sec)
		sectors[i] = sec
	}
	unstable := 0
	for i := 0; i < n; i++ {
		left := sectors[(i-1+n)%n]
		right := sectors[(i+1)%n]
		cur := sectors[i]
		for j := 0; j+1 < len(cur); j++ {
			a, b := cur[j], cur[j+1]
			l1 := sort.Search(len(left), func(x int) bool { return left[x] > a })
			l2 := sort.Search(len(left), func(x int) bool { return left[x] >= b })
			r1 := sort.Search(len(right), func(x int) bool { return right[x] > a })
			r2 := sort.Search(len(right), func(x int) bool { return right[x] >= b })
			if l2-l1 != r2-r1 {
				unstable++
			}
		}
	}
	return unstable
}

func buildInput(tc testCaseD) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, sec := range tc.sectors {
		sb.WriteString(fmt.Sprintf("%d", len(sec)))
		for _, v := range sec {
			sb.WriteString(fmt.Sprintf(" %d", v))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		tc := generateCase()
		input := buildInput(tc)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", i)
			os.Exit(1)
		}
		exp := compute(tc)
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\n", i, exp, val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
