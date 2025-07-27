package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		freq := make([]int, 26)
		for i := 0; i < len(s); i++ {
			freq[s[i]-'a']++
		}
		distinct := 0
		for i := 0; i < 26; i++ {
			if freq[i] > 0 {
				distinct++
			}
		}
		if distinct == 1 {
			fmt.Fprintln(out, s)
			continue
		}
		unique := -1
		for i := 0; i < 26; i++ {
			if freq[i] == 1 {
				unique = i
				break
			}
		}
		if unique != -1 {
			var b strings.Builder
			b.WriteByte(byte('a' + unique))
			freq[unique]--
			for i := 0; i < 26; i++ {
				for freq[i] > 0 {
					b.WriteByte(byte('a' + i))
					freq[i]--
				}
			}
			fmt.Fprintln(out, b.String())
			continue
		}

		first := 0
		for freq[first] == 0 {
			first++
		}
		n := len(s)
		cntFirst := freq[first]

		if cntFirst-2 <= n-cntFirst {
			// place two of first at start
			var b strings.Builder
			b.WriteByte(byte('a' + first))
			b.WriteByte(byte('a' + first))
			freq[first] -= 2
			others := make([]byte, 0, n-cntFirst)
			for i := first + 1; i < 26; i++ {
				for j := 0; j < freq[i]; j++ {
					others = append(others, byte('a'+i))
				}
			}
			sort.Slice(others, func(i, j int) bool { return others[i] < others[j] })
			pos := 0
			for pos < len(others) {
				b.WriteByte(others[pos])
				pos++
				if freq[first] > 0 {
					b.WriteByte(byte('a' + first))
					freq[first]--
				}
			}
			for freq[first] > 0 {
				b.WriteByte(byte('a' + first))
				freq[first]--
			}
			fmt.Fprintln(out, b.String())
			continue
		}

		// cannot place two first chars at start
		second := first + 1
		for freq[second] == 0 {
			second++
		}
		if distinct == 2 {
			var b strings.Builder
			b.WriteByte(byte('a' + first))
			for i := 0; i < freq[second]; i++ {
				b.WriteByte(byte('a' + second))
			}
			for i := 0; i < freq[first]-1; i++ {
				b.WriteByte(byte('a' + first))
			}
			fmt.Fprintln(out, b.String())
			continue
		}
		third := second + 1
		for freq[third] == 0 {
			third++
		}
		var b strings.Builder
		b.WriteByte(byte('a' + first))
		b.WriteByte(byte('a' + second))
		freq[first]--
		freq[second]--
		for freq[first] > 0 {
			b.WriteByte(byte('a' + first))
			freq[first]--
		}
		b.WriteByte(byte('a' + third))
		freq[third]--
		for i := 0; i < 26; i++ {
			for freq[i] > 0 {
				b.WriteByte(byte('a' + i))
				freq[i]--
			}
		}
		fmt.Fprintln(out, b.String())
	}
}
