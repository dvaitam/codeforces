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
   var sStr, tStr string
   fmt.Fscan(reader, &sStr, &tStr)
   s := []byte(sStr)
   t := []byte(tStr)
   if !isPerm(s, t) {
       fmt.Fprintln(writer, -1)
       return
   }
   var ans []int
   // Transform s into t using operations as in solE.cpp
   for i := 0; i < n; i++ {
       goal := t[n-1-i]
       // find position j
       var j int
       for j = i; j < n; j++ {
           if s[j] == goal {
               break
           }
       }
       // record operations: shift n, shift j, shift 1
       ans = append(ans, n, j, 1)
       // reverse suffix from j to end
       for l, r := j, n-1; l < r; l, r = l+1, r-1 {
           s[l], s[r] = s[r], s[l]
       }
       // rotate: move last element to front
       if n > 0 {
           last := s[n-1]
           s = append([]byte{last}, s[:n-1]...)
       }
   }
   // output answer
   fmt.Fprintln(writer, len(ans))
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}

// isPerm checks if a and b have same multiset of bytes
func isPerm(a, b []byte) bool {
   if len(a) != len(b) {
       return false
   }
   cnt := make([]int, 26)
   for _, c := range a {
       cnt[c-'a']++
   }
   for _, c := range b {
       cnt[c-'a']--
   }
   for _, v := range cnt {
       if v != 0 {
           return false
       }
   }
   return true
}
