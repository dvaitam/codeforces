package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   prev := "0"
   for i := 0; i < n; i++ {
       bi := b[i]
       // try same length
       L := len(prev)
       found := ""
       // prefix sums
       ps := make([]int, L+1)
       for j := 0; j < L; j++ {
           ps[j+1] = ps[j] + int(prev[j]-'0')
       }
       for pos := L - 1; pos >= 0 && found == ""; pos-- {
           orig := int(prev[pos] - '0')
           for d := orig + 1; d <= 9; d++ {
               sumPre := ps[pos] + d
               rem := bi - sumPre
               remPos := L - pos - 1
               if rem < 0 || rem > 9*remPos {
                   continue
               }
               // build t
               t := make([]byte, L)
               copy(t, []byte(prev))
               t[pos] = byte('0' + d)
               // fill rest minimal
               for k := pos + 1; k < L; k++ {
                   // minimal digit at k: x such that remaining fits
                   // x = max(0, rem - 9*(positions after k))
                   after := L - k - 1
                   x := rem - 9*after
                   if x < 0 {
                       x = 0
                   }
                   if x > 9 {
                       x = 9
                   }
                   t[k] = byte('0' + x)
                   rem -= x
               }
               found = string(t)
               break
           }
       }
       if found != "" {
           prev = found
       } else {
           // need longer length
           // minimal length for sum bi
           needLen := bi/9
           if bi%9 != 0 {
               needLen++
           }
           newLen := needLen
           if newLen <= len(prev) {
               newLen = len(prev) + 1
           }
           rem := bi
           t := make([]byte, newLen)
           for k := 0; k < newLen; k++ {
               remPos := newLen - k - 1
               // minimal digit
               minD := 0
               if k == 0 {
                   minD = 1
               }
               x := rem - 9*remPos
               if x < minD {
                   x = minD
               }
               if x > 9 {
                   x = 9
               }
               t[k] = byte('0' + x)
               rem -= x
           }
           prev = string(t)
       }
       // output
       fmt.Fprintln(writer, prev)
   }
}
