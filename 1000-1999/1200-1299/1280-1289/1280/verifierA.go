package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod int64 = 1e9 + 7

func solve(x int, s string) int64 {
	arr := []byte(s)
	lenMod := int64(len(arr))
	for i := 0; i < x; i++ {
		c := int(arr[i] - '0')
		diff := (lenMod - int64(i+1)) % mod
		if diff < 0 {
			diff += mod
		}
		lenMod = (lenMod + diff*int64(c-1)) % mod
		if len(arr) < x {
			suffix := make([]byte, len(arr)-i-1)
			copy(suffix, arr[i+1:])
			for j := 0; j < c-1 && len(arr) < x; j++ {
				need := x - len(arr)
				if need >= len(suffix) {
					arr = append(arr, suffix...)
				} else {
					arr = append(arr, suffix[:need]...)
				}
			}
		}
	}
	return lenMod % mod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	t := 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		x := rand.Intn(20) + 1
		l := rand.Intn(10) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('1' + rand.Intn(3))
		}
		s := string(b)
		fmt.Fprintln(&input, x)
		fmt.Fprintln(&input, s)
		expected[i] = fmt.Sprintf("%d", solve(x, s))
	}

	cmd := exec.Command(bin)
	cmd.Stdin = &input
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary error:", err)
		fmt.Print(string(out))
		return
	}
	outputs := strings.Fields(string(out))
	if len(outputs) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(outputs))
		fmt.Print(string(out))
		return
	}
	for i := 0; i < t; i++ {
		if outputs[i] != expected[i] {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, expected[i], outputs[i])
			return
		}
	}
	fmt.Println("All tests passed!")
}
