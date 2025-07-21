package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // positions of each char
   pos := make([][]int, 26)
   for i, ch := range s {
       pos[ch-'a'] = append(pos[ch-'a'], i)
   }
   // count pairs available
   totalPairs := 0
   for c := 0; c < 26; c++ {
       totalPairs += len(pos[c]) / 2
   }
   // determine target pairs (half-length) for length 100 palindrome
   target := totalPairs
   want100 := false
   if totalPairs >= 50 {
       target = 50
       want100 = true
   }
   l, r := 0, n-1
   left := make([]byte, 0, target)
   // greedy pick pairs
   for i := 0; i < target; i++ {
       bestC := -1
       bestL, bestR := 0, 0
       bestWidth := -1
       for c := 0; c < 26; c++ {
           arr := pos[c]
           // find first >= l
           il := sort.SearchInts(arr, l)
           if il >= len(arr) {
               continue
           }
           // find last <= r
           // upper bound of r
           ir := sort.SearchInts(arr, r+1) - 1
           if ir < il {
               continue
           }
           width := arr[ir] - arr[il]
           if width > bestWidth {
               bestWidth = width
               bestC = c
               bestL = arr[il]
               bestR = arr[ir]
           }
       }
       if bestC < 0 {
           break
       }
       left = append(left, byte('a'+bestC))
       l = bestL + 1
       r = bestR - 1
   }
   // build result
   var res []byte
   // left half
   res = append(res, left...)
   // center if needed
   if !want100 && l <= r {
       // any single char
       res = append(res, s[l])
   }
   // right half (mirror)
   for i := len(left) - 1; i >= 0; i-- {
       res = append(res, left[i])
   }
   // if want exactly 100 but got more (shouldn't), trim
   if want100 && len(res) > 100 {
       res = res[:100]
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   writer.Write(res)
}
