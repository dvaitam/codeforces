package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)

	type state struct {
		idx    int
		curVal byte
		build  []byte
	}

	n := len(s)
	res := make(map[string]struct{})
	stack := []state{{1, s[0], []byte{}}}
	for len(stack) > 0 {
		st := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if st.idx == n {
			final := append(st.build, st.curVal)
			res[string(final)] = struct{}{}
			continue
		}
		ch := s[st.idx]
		// Option 1: merge with current block
		st1 := state{st.idx + 1, maxByte(st.curVal, ch), append([]byte(nil), st.build...)}
		stack = append(stack, st1)
		// Option 2: start new block
		nb := append(append([]byte(nil), st.build...), st.curVal)
		st2 := state{st.idx + 1, ch, nb}
		stack = append(stack, st2)
	}

	const mod int64 = 1000000007
	ans := int64(len(res)) % mod
	fmt.Println(ans)
}

func maxByte(a, b byte) byte {
	if a < b {
		return b
	}
	return a
}
