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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([][]int64, 3)
       for i := 0; i < 3; i++ {
           a[i] = make([]int64, n)
           for j := 0; j < n; j++ {
               fmt.Fscan(reader, &a[i][j])
           }
       }
       // compute required sum = ceil(sum(a[0]) / 3)
       var sum0 int64
       for _, v := range a[0] {
           sum0 += v
       }
       need := (sum0 + 2) / 3
       perms := [6][3]int{
           {0, 1, 2}, {0, 2, 1}, {1, 0, 2},
           {1, 2, 0}, {2, 0, 1}, {2, 1, 0},
       }
       found := false
       var seg [3][2]int
       for _, p := range perms {
           ptr := 0
           ok := true
           var cur int64
           // reset segments
           for i := range seg {
               seg[i][0], seg[i][1] = 0, -1
           }
           for _, x := range p {
               start := ptr
               cur = 0
               for ptr < n && cur < need {
                   cur += a[x][ptr]
                   ptr++
               }
               if cur < need {
                   ok = false
                   break
               }
               seg[x][0] = start
               seg[x][1] = ptr - 1
           }
           if ok {
               for i := 0; i < 3; i++ {
                   // output 1-based indices
                   fmt.Fprint(writer, seg[i][0]+1, " ", seg[i][1]+1, " ")
               }
               fmt.Fprint(writer, '\n')
               found = true
               break
           }
       }
       if !found {
           fmt.Fprint(writer, "-1\n")
       }
   }
}
