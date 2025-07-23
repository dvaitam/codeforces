package main

import (
    "bufio"
    "fmt"
    "os"
)

var (
    s    string
    memo map[[2]int]string
)

func reverseString(str string) string {
    b := []byte(str)
    for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
        b[i], b[j] = b[j], b[i]
    }
    return string(b)
}

func minString(a, b string) string {
    if a == "" {
        return b
    }
    if b == "" {
        return a
    }
    if a < b {
        return a
    }
    return b
}

func solve(pos, k int) string {
    if pos == len(s) {
        return ""
    }
    if k == 1 {
        rem := s[pos:]
        rev := reverseString(rem)
        if rem < rev {
            return rem
        }
        return rev
    }
    key := [2]int{pos, k}
    if v, ok := memo[key]; ok {
        return v
    }
    best := ""
    n := len(s)
    for i := pos + 1; i <= n; i++ {
        if k-1 > n-i {
            continue
        }
        part := s[pos:i]
        rest := solve(i, k-1)
        best = minString(best, part+rest)
        best = minString(best, reverseString(part)+rest)
    }
    memo[key] = best
    return best
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Fscan(reader, &s)
    var k int
    fmt.Fscan(reader, &k)
    memo = make(map[[2]int]string)
    res := solve(0, k)
    writer := bufio.NewWriter(os.Stdout)
    fmt.Fprintln(writer, res)
    writer.Flush()
}

