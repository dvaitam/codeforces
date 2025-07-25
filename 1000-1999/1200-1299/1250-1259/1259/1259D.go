package main

import (
	"bufio"
	"fmt"
	"os"
)

func reverseString(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		words := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &words[i])
		}

		present := make(map[string]struct{}, n)
		for _, w := range words {
			present[w] = struct{}{}
		}

		idx01 := []int{}
		idx10 := []int{}
		idx00 := []int{}
		idx11 := []int{}
		for i, w := range words {
			if w[0] == '0' && w[len(w)-1] == '1' {
				idx01 = append(idx01, i)
			} else if w[0] == '1' && w[len(w)-1] == '0' {
				idx10 = append(idx10, i)
			} else if w[0] == '0' && w[len(w)-1] == '0' {
				idx00 = append(idx00, i)
			} else {
				idx11 = append(idx11, i)
			}
		}

		if len(idx01) == 0 && len(idx10) == 0 {
			if len(idx00) > 0 && len(idx11) > 0 {
				fmt.Fprintln(out, -1)
				continue
			}
			fmt.Fprintln(out, 0)
			continue
		}

		diff := len(idx01) - len(idx10)
		if diff == 0 || diff == 1 || diff == -1 {
			fmt.Fprintln(out, 0)
			continue
		}

		candidate01 := []int{}
		for _, i := range idx01 {
			r := reverseString(words[i])
			if _, ok := present[r]; !ok {
				candidate01 = append(candidate01, i+1)
			}
		}

		candidate10 := []int{}
		for _, i := range idx10 {
			r := reverseString(words[i])
			if _, ok := present[r]; !ok {
				candidate10 = append(candidate10, i+1)
			}
		}

		ans := []int{}
		if diff > 1 {
			need := diff / 2
			if len(candidate01) < need {
				fmt.Fprintln(out, -1)
				continue
			}
			ans = candidate01[:need]
		} else if diff < -1 {
			need := (-diff) / 2
			if len(candidate10) < need {
				fmt.Fprintln(out, -1)
				continue
			}
			ans = candidate10[:need]
		}

		fmt.Fprintln(out, len(ans))
		if len(ans) > 0 {
			for i, v := range ans {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			fmt.Fprintln(out)
		}
	}
}
