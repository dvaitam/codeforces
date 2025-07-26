package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var mod int64
var almostRoot [18]int64
var almostRootPow [18]int64

func multMod(a, b int64) int64 {
	res := (a * (b >> 16)) % mod
	res = (res << 16) % mod
	res = (res + a*(b&((1<<16)-1))) % mod
	return res
}

func prepare() {
	mod = 1
	for i := 0; i < 17; i++ {
		mod *= 5
	}
	almostRoot[0] = 2
	almostRoot[1] = 16
	almostRootPow[0] = 1
	almostRootPow[1] = 4
	for i := 2; i <= 17; i++ {
		x := multMod(almostRoot[i-1], almostRoot[i-1])
		x = multMod(x, x)
		x = multMod(x, almostRoot[i-1])
		almostRoot[i] = x
		almostRootPow[i] = almostRootPow[i-1] * 5
	}
}

func solve(ai int64) int64 {
	a := ai * 1000000
	a += 1 << 17
	a &= ^((1 << 17) - 1)
	if a%5 == 0 {
		a += 1 << 17
	}
	p := int64(60)
	r := (1 << 60) % mod
	curmod := int64(5)
	for modp := 1; modp <= 17; modp++ {
		for (r-a)%curmod != 0 {
			r = multMod(r, almostRoot[modp-1])
			p += almostRootPow[modp-1]
		}
		curmod *= 5
	}
	return p
}

func randomInput() []int64 {
	n := rand.Intn(5) + 1
	res := make([]int64, n)
	for i := range res {
		res[i] = rand.Int63n(1e6) + 1
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	prepare()
	const cases = 100
	for i := 0; i < cases; i++ {
		arr := randomInput()
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for _, v := range arr {
			sb.WriteString(fmt.Sprintf("%d\n", v))
		}
		input := sb.String()
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			fmt.Printf("program output:\n%s\n", string(out))
			return
		}
		gotLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(gotLines) != len(arr) {
			fmt.Printf("case %d failed: line count mismatch\ninput:\n%soutput:\n%s", i+1, input, string(out))
			return
		}
		for j, ai := range arr {
			want := solve(ai)
			if strings.TrimSpace(gotLines[j]) != strconv.FormatInt(want, 10) {
				fmt.Printf("case %d failed line %d expected %d got %s\n", i+1, j+1, want, gotLines[j])
				fmt.Printf("input:\n%s", input)
				return
			}
		}
	}
	fmt.Printf("OK %d cases\n", cases)
}
