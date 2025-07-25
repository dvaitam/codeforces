package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(binary string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	req := []int{1, 1, 2, 7, 4}
	for i := 0; i < 100; i++ {
		stocks := []int{(i%10 + 1), ((i+3)%10 + 1), ((i*2)%10 + 1), ((i*3)%10 + 1), ((i*5)%10 + 1)}
		input := fmt.Sprintf("%d %d %d %d %d\n", stocks[0], stocks[1], stocks[2], stocks[3], stocks[4])
		ans := stocks[0] / req[0]
		for j := 1; j < 5; j++ {
			v := stocks[j] / req[j]
			if v < ans {
				ans = v
			}
		}
		want := fmt.Sprint(ans)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != want {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
