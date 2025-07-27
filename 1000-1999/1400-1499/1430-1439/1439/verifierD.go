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

func brute(n, m, p int) int {
	res := 0
	a := make([]int, m)
	b := make([]int, m)
	var dfs func(int)
	dfs = func(idx int) {
		if idx == m {
			seats := make([]bool, n)
			sum := 0
			for i := 0; i < m; i++ {
				pos := a[i]
				if b[i] == 0 { // L
					seat := -1
					for j := pos; j >= 0; j-- {
						if !seats[j] {
							seat = j
							break
						}
					}
					if seat == -1 {
						return
					}
					seats[seat] = true
					if seat != pos {
						if seat > pos {
							sum += seat - pos
						} else {
							sum += pos - seat
						}
					}
				} else {
					seat := -1
					for j := pos; j < n; j++ {
						if !seats[j] {
							seat = j
							break
						}
					}
					if seat == -1 {
						return
					}
					seats[seat] = true
					if seat != pos {
						if seat > pos {
							sum += seat - pos
						} else {
							sum += pos - seat
						}
					}
				}
			}
			res = (res + sum) % p
			return
		}
		for i := 0; i < n; i++ {
			a[idx] = i
			for t := 0; t < 2; t++ {
				b[idx] = t
				dfs(idx + 1)
			}
		}
	}
	dfs(0)
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/bin")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(3) + 1
		m := rand.Intn(n) + 1
		p := 1000000007
		input := fmt.Sprintf("%d %d %d\n", n, m, p)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d exec err %v\n", tcase, err)
			return
		}
		sc := bufio.NewScanner(strings.NewReader(out))
		sc.Split(bufio.ScanWords)
		if !sc.Scan() {
			fmt.Printf("test %d no output\n", tcase)
			return
		}
		got, err := strconv.Atoi(sc.Text())
		if err != nil {
			fmt.Printf("test %d bad int\n", tcase)
			return
		}
		expect := brute(n, m, p)
		if got != expect {
			fmt.Printf("test %d expected %d got %d\n", tcase, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
