package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func simulate(n int, a []int, queries [][3]int) []int {
	res := []int{}
	arr := append([]int(nil), a...)
	for _, q := range queries {
		t, x, y := q[0], q[1], q[2]
		if t == 1 {
			for i := 0; i < x; i++ {
				if y > arr[i] {
					arr[i] = y
				}
			}
		} else {
			money := y
			cnt := 0
			for i := x - 1; i < n; i++ {
				if money >= arr[i] {
					money -= arr[i]
					cnt++
				}
			}
			res = append(res, cnt)
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/bin")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(5) + 2
		q := rand.Intn(5) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(10) + 1
		}
		queries := make([][3]int, q)
		for i := 0; i < q; i++ {
			t := rand.Intn(2) + 1
			x := rand.Intn(n) + 1
			y := rand.Intn(10) + 1
			queries[i] = [3]int{t, x, y}
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d ", a[i]))
		}
		sb.WriteString("\n")
		for _, qq := range queries {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", qq[0], qq[1], qq[2]))
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d exec err %v\n", tcase, err)
			return
		}
		sc := bufio.NewScanner(strings.NewReader(out))
		sc.Split(bufio.ScanWords)
		sim := simulate(n, a, queries)
		for i, ans := range sim {
			if !sc.Scan() {
				fmt.Printf("test %d output missing\n", tcase)
				return
			}
			got, err := strconv.Atoi(sc.Text())
			if err != nil {
				fmt.Printf("test %d invalid int\n", tcase)
				return
			}
			if got != ans {
				fmt.Printf("test %d query %d expected %d got %d\n", tcase, i, ans, got)
				return
			}
		}
	}
	fmt.Println("All tests passed")
}
