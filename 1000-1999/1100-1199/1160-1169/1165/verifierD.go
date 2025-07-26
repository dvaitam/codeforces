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

const numTestsD = 100

func countDivisors(x int64) int {
	cnt := 1
	n := x
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			e := 0
			for n%i == 0 {
				n /= i
				e++
			}
			cnt *= (e + 1)
		}
	}
	if n > 1 {
		cnt *= 2
	}
	return cnt
}

func solveQuery(divs []int64) int64 {
	sort.Slice(divs, func(i, j int) bool { return divs[i] < divs[j] })
	candidate := divs[0] * divs[len(divs)-1]
	n := len(divs)
	for i := 0; i < n; i++ {
		if divs[i]*divs[n-1-i] != candidate {
			return -1
		}
	}
	if countDivisors(candidate)-2 != n {
		return -1
	}
	return candidate
}

func run(binary, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func divisors(x int64) []int64 {
	var d []int64
	for i := int64(1); i*i <= x; i++ {
		if x%i == 0 {
			d = append(d, i)
			if i*i != x {
				d = append(d, x/i)
			}
		}
	}
	sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
	return d
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(4)
	for tcase := 1; tcase <= numTestsD; tcase++ {
		t := rand.Intn(3) + 1
		input := fmt.Sprintf("%d\n", t)
		var expectedLines []string
		for q := 0; q < t; q++ {
			valid := rand.Intn(2) == 0
			var n int
			var divs []int64
			var ans int64
			if valid {
				x := int64(rand.Intn(5000) + 2)
				d := divisors(x)
				// remove 1 and x
				var list []int64
				for _, v := range d {
					if v != 1 && v != x {
						list = append(list, v)
					}
				}
				divs = list
				n = len(divs)
				ans = solveQuery(append([]int64(nil), divs...))
			} else {
				// create invalid set
				x := int64(rand.Intn(5000) + 2)
				d := divisors(x)
				var list []int64
				for _, v := range d {
					if v != 1 && v != x {
						list = append(list, v)
					}
				}
				if len(list) > 0 {
					list[0]++ // modify to break
				} else {
					list = append(list, x+1)
				}
				divs = list
				n = len(divs)
				ans = -1
			}
			// shuffle divs
			rand.Shuffle(n, func(i, j int) { divs[i], divs[j] = divs[j], divs[i] })
			input += fmt.Sprintf("%d\n", n)
			for i, v := range divs {
				if i > 0 {
					input += " "
				}
				input += fmt.Sprintf("%d", v)
			}
			input += "\n"
			expectedLines = append(expectedLines, fmt.Sprintf("%d", ans))
		}
		expect := strings.Join(expectedLines, "\n") + "\n"
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("case %d failed to run: %v\noutput:%s\n", tcase, err, out)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:%sexpected:%s got:%s\n", tcase, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
