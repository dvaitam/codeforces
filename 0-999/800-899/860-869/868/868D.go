package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxK = 15

type Info struct {
	pref, suff string
	length     int
	sub        [MaxK + 1][]bool
}

func computeSubstrings(s string) Info {
	info := Info{length: len(s)}
	if len(s) <= MaxK {
		info.pref = s
		info.suff = s
	} else {
		info.pref = s[:MaxK]
		info.suff = s[len(s)-MaxK:]
	}
	bytes := []byte(s)
	for k := 1; k <= MaxK; k++ {
		size := 1 << uint(k)
		arr := make([]bool, size)
		if len(bytes) >= k {
			val := 0
			for t := 0; t < k; t++ {
				val = val<<1 | int(bytes[t]-'0')
			}
			arr[val] = true
			mask := size - 1
			for i := k; i < len(bytes); i++ {
				val = ((val << 1) & mask) | int(bytes[i]-'0')
				arr[val] = true
			}
		}
		info.sub[k] = arr
	}
	return info
}

func merge(a, b Info) Info {
	res := Info{}
	res.length = a.length + b.length
	// prefix
	if len(a.pref) == MaxK {
		res.pref = a.pref
	} else {
		tmp := a.pref + b.pref
		if len(tmp) > MaxK {
			tmp = tmp[:MaxK]
		}
		res.pref = tmp
	}
	// suffix
	if len(b.suff) == MaxK {
		res.suff = b.suff
	} else {
		tmp := a.suff + b.suff
		if len(tmp) > MaxK {
			tmp = tmp[len(tmp)-MaxK:]
		}
		res.suff = tmp
	}
	cross := a.suff + b.pref
	crossBytes := []byte(cross)
	for k := 1; k <= MaxK; k++ {
		size := 1 << uint(k)
		arr := make([]bool, size)
		for i := 0; i < size; i++ {
			if a.sub[k][i] || b.sub[k][i] {
				arr[i] = true
			}
		}
		if len(crossBytes) >= k {
			val := 0
			for t := 0; t < k; t++ {
				val = val<<1 | int(crossBytes[t]-'0')
			}
			arr[val] = true
			mask := size - 1
			for i := k; i < len(crossBytes); i++ {
				val = ((val << 1) & mask) | int(crossBytes[i]-'0')
				arr[val] = true
			}
		}
		res.sub[k] = arr
	}
	return res
}

func maxComplete(info Info) int {
	for k := MaxK; k >= 1; k-- {
		all := true
		for i := 0; i < 1<<uint(k); i++ {
			if !info.sub[k][i] {
				all = false
				break
			}
		}
		if all {
			return k
		}
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	infos := make([]Info, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		infos[i] = computeSubstrings(s)
	}
	var m int
	fmt.Fscan(in, &m)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		info := merge(infos[a], infos[b])
		infos = append(infos, info)
		ans := maxComplete(info)
		fmt.Println(ans)
	}
}
