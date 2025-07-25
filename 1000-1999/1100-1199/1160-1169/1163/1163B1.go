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

   freq := make(map[int]int)
   freqCount := make(map[int]int)
   maxValid := 0

   for i, color := range u {
       // update frequency of color
       old := freq[color]
       if old > 0 {
           cnt := freqCount[old]
           if cnt <= 1 {
               delete(freqCount, old)
           } else {
               freqCount[old] = cnt - 1
           }
       }
       freq[color] = old + 1
       freqCount[old+1]++

       if validPrefix(freqCount) {
           maxValid = i + 1
       }
   }
   fmt.Fprintln(writer, maxValid)
}

// validPrefix checks if current prefix can be made equal freq by removing one element
func validPrefix(freqCount map[int]int) bool {
   if len(freqCount) == 1 {
       for f, cnt := range freqCount {
           // either single color, or all colors freq=1
           if f == 1 || cnt == 1 {
               return true
           }
           return false
       }
   }
   if len(freqCount) == 2 {
       // extract two frequencies
       var f1, f2, c1, c2 int
       i := 0
       for f, cnt := range freqCount {
           if i == 0 {
               f1, c1 = f, cnt
           } else {
               f2, c2 = f, cnt
           }
           i++
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
       // case: one color has freq greater by 1
       if f2 == f1+1 && c2 == 1 {
           return true
       }
   }
   return false
}
