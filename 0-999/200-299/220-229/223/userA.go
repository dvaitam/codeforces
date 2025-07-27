package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    sBytes, _, _ := reader.ReadLine()
    s := string(sBytes)
    n := len(s)

    dp := make([]int, n)
    pref := make([]int, n+1)
    for i := 0; i < n; i++ {
        pref[i+1] = pref[i]
        if s[i] == '[' {
            pref[i+1]++
        }
    }

    maxW := 0
    bestEnd := -1

    for i := 0; i < n; i++ {
        if s[i] == ')' || s[i] == ']' {
            var match byte
            if s[i] == ')' {
                match = '('
            } else {
                match = '['
            }
            prevLen := 0
            if i > 0 {
                prevLen = dp[i-1]
            }
            j := i - 1 - prevLen
            if j >= 0 && s[j] == match {
                dp[i] = prevLen + 2
                if j-1 >= 0 {
                    dp[i] += dp[j-1]
                }
                l := i - dp[i] + 1
                w := pref[i+1] - pref[l]
                if w > maxW {
                    maxW = w
                    bestEnd = i
                }
            }
        }
    }

    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    fmt.Fprintln(writer, maxW)
    if maxW == 0 {
        fmt.Fprintln(writer, "")
    } else {
        l := bestEnd - dp[bestEnd] + 1
        fmt.Fprintln(writer, s[l:bestEnd+1])
    }
}