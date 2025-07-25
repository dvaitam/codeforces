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
   fmt.Fscan(reader, &n)
   u := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &u[i])
   }

   cnt := make(map[int]int)
   freqCnt := make(map[int]int)
   ans := 0

   for i, color := range u {
       prev := cnt[color]
       if prev > 0 {
           // decrement frequency count for prev
           freqCnt[prev]--
           if freqCnt[prev] == 0 {
               delete(freqCnt, prev)
           }
       }
       // increment count for color
       cnt[color] = prev + 1
       freqCnt[prev+1]++

       // check valid prefix of length i+1
       if validPrefix(freqCnt) {
           ans = i + 1
       }
   }
   fmt.Fprintln(writer, ans)
}

// validPrefix checks if by removing one element we can make all remaining colors have equal counts
func validPrefix(freqCnt map[int]int) bool {
   l := len(freqCnt)
   if l == 1 {
       for f, cnt := range freqCnt {
           // all colors have same frequency f
           // if f == 1 (all appear once) or only one color exists
           if f == 1 || cnt == 1 {
               return true
           }
       }
       return false
   }
   if l == 2 {
       var f1, f2 int
       var c1, c2 int
       first := true
       for f, cnt := range freqCnt {
           if first {
               f1, c1 = f, cnt
               first = false
           } else {
               f2, c2 = f, cnt
           }
       }
       // ensure f1 < f2
       if f1 > f2 {
           f1, f2 = f2, f1
           c1, c2 = c2, c1
       }
       // case: one color has freq 1
       if f1 == 1 && c1 == 1 {
           return true
       }
       // case: one color has freq one greater than others
       if f2 == f1+1 && c2 == 1 {
           return true
       }
   }
   return false
}
