package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func expected(guest, host, pile string) string {
	if len(pile) != len(guest)+len(host) {
		return "NO"
	}
	cnt := make([]int, 26)
	for _, c := range guest {
		cnt[c-'A']++
	}
	for _, c := range host {
		cnt[c-'A']++
	}
	for _, c := range pile {
		cnt[c-'A']--
	}
	for _, v := range cnt {
		if v != 0 {
			return "NO"
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("failed to open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	tests := 0
	for {
		var guest, host, pile string
		if !scanner.Scan() {
			break
		}
		guest = scanner.Text()
		if !scanner.Scan() {
			fmt.Println("unexpected EOF")
			os.Exit(1)
		}
		host = scanner.Text()
		if !scanner.Scan() {
			fmt.Println("unexpected EOF")
			os.Exit(1)
		}
		pile = scanner.Text()
		expect := expected(guest, host, pile)
		input := fmt.Sprintf("%s\n%s\n%s\n", guest, host, pile)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tests+1, err)
			os.Exit(1)
		}
		res := string(bytes.TrimSpace(out))
		if res != expect {
			fmt.Printf("test %d failed: expected %s got %s (guest=%s host=%s pile=%s)\n", tests+1, expect, res, guest, host, pile)
			os.Exit(1)
		}
		tests++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", tests)
}
