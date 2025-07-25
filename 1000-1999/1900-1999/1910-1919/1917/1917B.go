package main

import (
    "bufio"
    "fmt"
    "os"
)

const base uint64 = 911382323

func countDistinct(s string) int {
    n := len(s)
    pow := make([]uint64, n+1)
    pref := make([]uint64, n+1)
    pow[0] = 1
    for i := 0; i < n; i++ {
        pow[i+1] = pow[i] * base
        pref[i+1] = pref[i]*base + uint64(s[i])
    }
    hash := func(l, r int) uint64 {
        return pref[r] - pref[l]*pow[r-l]
    }
    uniq := make(map[uint64]struct{})
    // all suffixes
    for i := 0; i < n; i++ {
        uniq[hash(i, n)] = struct{}{}
    }
    // single characters
    for i := 0; i < n; i++ {
        uniq[hash(i, i+1)] = struct{}{}
    }
    // char + suffix
    var seen [26]bool
    for j := 2; j < n; j++ {
        seen[s[j-2]-'a'] = true
        suf := hash(j, n)
        pw := pow[n-j]
        for c := 0; c < 26; c++ {
            if seen[c] {
                ch := uint64('a' + c)
                uniq[ch*pw+suf] = struct{}{}
            }
        }
    }
    return len(uniq)
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for i := 0; i < t; i++ {
        var n int
        fmt.Fscan(in, &n)
        var s string
        fmt.Fscan(in, &s)
        fmt.Fprintln(out, countDistinct(s))
    }
}

