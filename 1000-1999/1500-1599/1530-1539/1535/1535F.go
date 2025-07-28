package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	strs := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &strs[i])
	}
	// group by sorted key
	groups := make(map[string][]string)
	for _, s := range strs {
		bs := []byte(s)
		sort.Slice(bs, func(i, j int) bool { return bs[i] < bs[j] })
		key := string(bs)
		groups[key] = append(groups[key], s)
	}
	totalPairs := int64(n) * (int64(n) - 1) / 2
	var intraPairs int64
	var sumSame int64 // sum f within groups
	for _, grp := range groups {
		g := len(grp)
		if g <= 1 {
			continue
		}
		intraPairs += int64(g) * (int64(g) - 1) / 2
		L := len(grp[0])
		// decide method based on length
		brute := func() {
			// brute pairs
			// convert strings to byte slices
			arr := make([][]byte, g)
			for i, s := range grp {
				arr[i] = []byte(s)
			}
			var cnt1 int64
			for i := 0; i < g; i++ {
				a := arr[i]
				for j := i + 1; j < g; j++ {
					b := arr[j]
					// find l,r
					l := 0
					for l < L && a[l] == b[l] {
						l++
					}
					if l == L {
						continue
					}
					r := L - 1
					for r > l && a[r] == b[r] {
						r--
					}
					// check multiset equal
					var ca [26]int
					for k := l; k <= r; k++ {
						ca[a[k]-'a']++
						ca[b[k]-'a']--
					}
					ok := true
					for _, v := range ca {
						if v != 0 {
							ok = false
							break
						}
					}
					if !ok {
						continue
					}
					// check if b[l..r] sorted
					bs := b
					sortedB := true
					for k := l; k < r; k++ {
						if bs[k] > bs[k+1] {
							sortedB = false
							break
						}
					}
					if sortedB {
						cnt1++
						continue
					}
					// check a[l..r] sorted
					as := a
					sortedA := true
					for k := l; k < r; k++ {
						if as[k] > as[k+1] {
							sortedA = false
							break
						}
					}
					if sortedA {
						cnt1++
					}
				}
			}
			// each such pair f=1, remaining pairs f=2
			sumSame += cnt1 + (int64(g)*(int64(g)-1)/2-cnt1)*2
		}
		enumeration := func() {
			// build index map
			idx := make(map[string]int, g)
			for i, s := range grp {
				idx[s] = i
			}
			L := len(grp[0])
			var cnt1 int64
			neighborSeen := make([]bool, g)
			// temp buffers
			sbuf := make([]byte, L)
			for i, s := range grp {
				bs := []byte(s)
				// reset seen
				for j := range neighborSeen {
					neighborSeen[j] = false
				}
				// try all segments
				for l := 0; l < L; l++ {
					for r := l + 1; r < L; r++ {
						// skip if segment already sorted
						sortedSeg := true
						for k := l; k < r; k++ {
							if bs[k] > bs[k+1] {
								sortedSeg = false
								break
							}
						}
						if sortedSeg {
							continue
						}
						// count frequencies
						var cnt [26]int
						for k := l; k <= r; k++ {
							cnt[bs[k]-'a']++
						}
						// build sbuf: prefix
						copy(sbuf, bs)
						// build sorted segment
						p := l
						for c := 0; c < 26; c++ {
							for t := 0; t < cnt[c]; t++ {
								sbuf[p] = byte('a' + c)
								p++
							}
						}
						// suffix already copied
						tkey := string(sbuf)
						if j, ok := idx[tkey]; ok && j > i && !neighborSeen[j] {
							neighborSeen[j] = true
							cnt1++
						}
					}
				}
			}
			sumSame += cnt1 + (int64(g)*(int64(g)-1)/2-cnt1)*2
		}
		// choose method: enumeration for small L, else brute
		if len(grp[0]) <= 20 {
			enumeration()
		} else {
			brute()
		}
	}
	// sum across groups: sumSame for within, impossible pairs get 1337
	sum := sumSame + (totalPairs-intraPairs)*1337
	fmt.Println(sum)
}
