package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"time"
)

func expectedC(a []int) int64 {
	sort.Ints(a)
	n := len(a)
	var ans int64
	pref := 0
	for i := 0; i < n; {
		j := i
		for j < n && a[j] == a[i] {
			j++
		}
		pref += j - i
		if pref < n {
			v := int64(pref) * int64(n-pref)
			if v > ans {
				ans = v
			}
		}
		i = j
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]int64, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(9) + 2
		fmt.Fprintln(&input, n)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(20) + 1
			fmt.Fprint(&input, arr[j])
			if j+1 < n {
				fmt.Fprint(&input, " ")
			}
		}
		fmt.Fprintln(&input)
		expected[i] = expectedC(append([]int(nil), arr...))
	}

	cmd := exec.Command(cand)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("candidate run error:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(outBytes))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("missing output for test", i+1)
			os.Exit(1)
		}
		var ans int64
		fmt.Sscan(scanner.Text(), &ans)
		if ans != expected[i] {
			fmt.Printf("wrong answer on test %d: expected %d got %d\n", i+1, expected[i], ans)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
