package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func feasible(m int, s string, maxc byte) bool {
    n := len(s)
    last := -1
    for i := 0; i < n; i++ {
        if s[i] <= maxc { last = i }
        if i >= m-1 {
            // window ending at i requires a pick within [i-m+1, i]
            if last < i-m+1 { return false }
        }
    }
    return true
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()
    var m int
    var s string
    if _, err := fmt.Fscan(in, &m); err != nil { return }
    if _, err := fmt.Fscan(in, &s); err != nil { return }
    n := len(s)
    // binary search minimal max character allowed
    lo, hi := byte('a'), byte('z')
    ans := byte('z')
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if feasible(m, s, mid) {
            ans = mid
            hi = mid - 1
        } else {
            lo = mid + 1
        }
    }
    // collect all chars < ans
    cnt := make([]int, 26)
    for i := 0; i < n; i++ {
        if s[i] < ans { cnt[int(s[i]-'a')]++ }
    }
    var b strings.Builder
    for c := 0; c < int(ans-'a'); c++ {
        if cnt[c] > 0 { b.WriteString(strings.Repeat(string('a'+byte(c)), cnt[c])) }
    }
    // count needed ans occurrences to cover remaining windows
    lastSmall, lastEq := -1, -1
    lastSmallPos := make([]int, n)
    lastEqPos := make([]int, n)
    for i := 0; i < n; i++ {
        if s[i] < ans { lastSmall = i }
        if s[i] == ans { lastEq = i }
        lastSmallPos[i] = lastSmall
        lastEqPos[i] = lastEq
    }
    picksAns := 0
    pos := -1
    for pos < n-m {
        reach := pos + m
        ns := lastSmallPos[reach]
        if ns > pos {
            pos = ns
        } else {
            picksAns++
            pos = lastEqPos[reach]
        }
    }
    if picksAns > 0 { b.WriteString(strings.Repeat(string(ans), picksAns)) }
    fmt.Fprintln(out, b.String())
}
