package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    var s string
    fmt.Fscan(reader, &s)

    // Count distinct Pokemon types
    var totalTypes int
    var freq [128]int
    for i := 0; i < n; i++ {
        b := s[i]
        if freq[b] == 0 {
            totalTypes++
        }
        freq[b]++
    }

    // Sliding window to find minimal segment containing all types
    var windowFreq [128]int
    have := 0
    minLen := n
    left := 0
    for right := 0; right < n; right++ {
        c := s[right]
        windowFreq[c]++
        if windowFreq[c] == 1 {
            have++
        }
        // When window has all types, try to shrink from left
        for have == totalTypes {
            currLen := right - left + 1
            if currLen < minLen {
                minLen = currLen
            }
            lc := s[left]
            windowFreq[lc]--
            if windowFreq[lc] == 0 {
                have--
            }
            left++
        }
    }

    fmt.Fprintln(writer, minLen)
}
