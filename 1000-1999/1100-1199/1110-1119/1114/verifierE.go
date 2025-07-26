package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func genTests() [][]int {
	rand.Seed(5)
	tests := make([][]int, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 2
		d := rand.Intn(10) + 1
		x := rand.Intn(20)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = x + j*d
		}
		rand.Shuffle(n, func(a, b int) { arr[a], arr[b] = arr[b], arr[a] })
		tests[i] = arr
	}
	return tests
}

func runCase(bin string, arr []int) error {
	cmd := exec.Command(bin)
	inPipe, _ := cmd.StdinPipe()
	outPipe, _ := cmd.StdoutPipe()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(outPipe)
	writer := bufio.NewWriter(inPipe)
	fmt.Fprintln(writer, len(arr))
	writer.Flush()
	queries := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("process terminated: %v %s", err, stderr.String())
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			var typ, val int
			fmt.Sscanf(line, "? %d %d", &typ, &val)
			queries++
			if queries > 60 {
				return fmt.Errorf("too many queries")
			}
			if typ == 1 {
				if val < 1 || val > len(arr) {
					return fmt.Errorf("bad index")
				}
				fmt.Fprintln(writer, arr[val-1])
				writer.Flush()
			} else if typ == 2 {
				res := 0
				for _, v := range arr {
					if v > val {
						res = 1
						break
					}
				}
				fmt.Fprintln(writer, res)
				writer.Flush()
			} else {
				return fmt.Errorf("unknown query %s", line)
			}
		} else if strings.HasPrefix(line, "!") {
			var x, d int
			fmt.Sscanf(line, "! %d %d", &x, &d)
			sorted := append([]int(nil), arr...)
			sort.Ints(sorted)
			realX := sorted[0]
			realD := 0
			if len(sorted) > 1 {
				realD = sorted[1] - sorted[0]
			}
			if x != realX || d != realD {
				return fmt.Errorf("wrong answer exp %d %d got %d %d", realX, realD, x, d)
			}
			break
		} else if line == "" {
			continue
		} else {
			return fmt.Errorf("unexpected output: %s", line)
		}
	}
	inPipe.Close()
	writer.Flush()
	err := cmd.Wait()
	if err != nil {
		return fmt.Errorf("runtime error: %v %s", err, stderr.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, arr := range tests {
		if err := runCase(bin, arr); err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
