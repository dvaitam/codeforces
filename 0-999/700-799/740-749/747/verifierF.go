package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

var t int
var memo map[state]uint64

type state struct {
	rem  uint8
	defs [16]uint8
}

func count(rem int, defs [16]uint8) uint64 {
	if rem == 0 {
		return 1
	}
	st := state{uint8(rem), defs}
	if v, ok := memo[st]; ok {
		return v
	}
	var total uint64
	for d := 0; d < 16; d++ {
		if int(defs[d]) < t {
			defs[d]++
			total += count(rem-1, defs)
			defs[d]--
		}
	}
	memo[st] = total
	return total
}

func kth(k uint64, tt int) string {
	t = tt
	memo = make(map[state]uint64)
	var length int
	for l := 1; l <= 16*t; l++ {
		var cnt uint64
		for d := 1; d < 16; d++ {
			if t > 0 {
				var defs [16]uint8
				defs[d] = 1
				cnt += count(l-1, defs)
			}
		}
		if cnt >= k {
			length = l
			break
		}
		k -= cnt
	}
	var defs [16]uint8
	rem := length
	res := make([]byte, 0, length)
	for pos := 0; pos < length; pos++ {
		for d := 0; d < 16; d++ {
			if pos == 0 && d == 0 {
				continue
			}
			if int(defs[d]) >= t {
				continue
			}
			defs[d]++
			cnt := count(rem-1, defs)
			if cnt >= k {
				var ch byte
				if d < 10 {
					ch = byte('0' + d)
				} else {
					ch = byte('a' + byte(d-10))
				}
				res = append(res, ch)
				rem--
				break
			} else {
				k -= cnt
				defs[d]--
			}
		}
	}
	return string(res)
}

func runCase(bin string, k uint64, t int) error {
	input := fmt.Sprintf("%d %d\n", k, t)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	got := strings.TrimSpace(string(out))
	expected := kth(k, t)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for tcase := 1; tcase <= 120; tcase++ {
		tt := rand.Intn(3) + 1
		k := uint64(rand.Intn(300) + 1)
		if err := runCase(bin, k, tt); err != nil {
			fmt.Printf("Test %d failed: %v\n", tcase, err)
			return
		}
	}
	fmt.Println("OK")
}
