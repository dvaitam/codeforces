package main

import (
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

func run(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    out, err := cmd.CombinedOutput()
    return string(out), err
}

func existsDivisible(nums []int) int {
    for i, v := range nums {
        ok := true
        for j, other := range nums {
            if i == j {
                continue
            }
            if v%other != 0 {
                ok = false
                break
            }
        }
        if ok {
            return 1
        }
    }
    return 0
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rand.Seed(42)
    for t := 0; t < 100; t++ {
        n := rand.Intn(9) + 2
        nums := make([]int, n)
        for i := 0; i < n; i++ {
            nums[i] = rand.Intn(100) + 1
        }
        fields := make([]string, n)
        for i, v := range nums {
            fields[i] = fmt.Sprintf("%d", v)
        }
        in := strings.Join(fields, " ") + "\n"
        want := fmt.Sprintf("%d", existsDivisible(nums))
        out, err := run(bin, in)
        if err != nil {
            fmt.Printf("test %d runtime error: %v\n", t+1, err)
            os.Exit(1)
        }
        out = strings.TrimSpace(out)
        if out != want {
            fmt.Printf("test %d failed: expected %q got %q\n", t+1, want, out)
            os.Exit(1)
        }
    }
    fmt.Println("OK")
}
