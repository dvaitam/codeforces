package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

func nondecreasing(nums []int) bool {
	for i := 1; i < len(nums); i++ {
		if nums[i] < nums[i-1] {
			return false
		}
	}
	return true
}

func expectedAnswer(arr []int) bool {
	n := len(arr)
	lastDigit := -1
	for i := 0; i < n; i++ {
		if arr[i] < 10 {
			lastDigit = i
		}
	}
	for pos := lastDigit + 1; pos <= n; pos++ {
		digits := make([]int, 0, pos*2)
		for i := 0; i < pos; i++ {
			if arr[i] >= 10 {
				digits = append(digits, arr[i]/10, arr[i]%10)
			} else {
				digits = append(digits, arr[i])
			}
		}
		if !nondecreasing(digits) {
			continue
		}
		good := true
		for i := pos + 1; i < n; i++ {
			if arr[i] < arr[i-1] {
				good = false
				break
			}
		}
		if good {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(42)
	const t = 100
	ns := make([]int, t)
	arrays := make([][]int, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(49) + 2 // 2..50
		ns[i] = n
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(100) // 0..99
		}
		arrays[i] = arr
	}

	var input bytes.Buffer
	fmt.Fprintf(&input, "%d\n", t)
	for i := 0; i < t; i++ {
		fmt.Fprintf(&input, "%d\n", ns[i])
		for j, v := range arrays[i] {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", i+1)
			os.Exit(1)
		}
		got := scanner.Text()
		want := "NO"
		if expectedAnswer(arrays[i]) {
			want = "YES"
		}
		if got != want {
			fmt.Printf("case %d: expected %s, got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("warning: extra output detected")
	}
	fmt.Println("All tests passed!")
}
