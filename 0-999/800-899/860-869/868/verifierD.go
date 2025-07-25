package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MaxK = 15

type Info struct {
	pref, suff string
	length     int
	sub        [MaxK + 1][]bool
}

type Test struct {
	in  string
	out string
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
	bytesArr := []byte(s)
	for k := 1; k <= MaxK; k++ {
		size := 1 << uint(k)
		arr := make([]bool, size)
		if len(bytesArr) >= k {
			val := 0
			for t := 0; t < k; t++ {
				val = val<<1 | int(bytesArr[t]-'0')
			}
			arr[val] = true
			mask := size - 1
			for i := k; i < len(bytesArr); i++ {
				val = ((val << 1) & mask) | int(bytesArr[i]-'0')
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
	if len(a.pref) == MaxK {
		res.pref = a.pref
	} else {
		tmp := a.pref + b.pref
		if len(tmp) > MaxK {
			tmp = tmp[:MaxK]
		}
		res.pref = tmp
	}
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

func oracle(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	infos := make([]Info, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		infos[i] = computeSubstrings(s)
	}
	var m int
	fmt.Fscan(in, &m)
	var sb strings.Builder
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		info := merge(infos[a], infos[b])
		infos = append(infos, info)
		ans := maxComplete(info)
		fmt.Fprintln(&sb, ans)
	}
	return strings.TrimSpace(sb.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func randomBinary(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		if rng.Intn(2) == 1 {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
	}
	return string(b)
}

func genCase(rng *rand.Rand) Test {
	n := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	strs := make([]string, n)
	for i := 0; i < n; i++ {
		s := randomBinary(rng, rng.Intn(4)+1)
		strs[i] = s
		fmt.Fprintln(&sb, s)
	}
	m := rng.Intn(3) + 1
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		a := rng.Intn(len(strs)) + 1
		b := rng.Intn(len(strs)) + 1
		fmt.Fprintf(&sb, "%d %d\n", a, b)
		strs = append(strs, "") // placeholder for new string
	}
	input := sb.String()
	out := oracle(input)
	return Test{input, out}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(4))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
