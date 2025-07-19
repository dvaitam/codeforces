package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var zh [4][4]int
   var zhu [4]int
   // read pairwise counts
   for i := 0; i <= 2; i++ {
       for j := i + 1; j <= 3; j++ {
           if _, err := fmt.Fscan(reader, &zh[i][j]); err != nil {
               return
           }
           zh[j][i] = zh[i][j]
           zhu[i] += zh[i][j]
           zhu[j] += zh[j][i]
       }
   }
   // find minimum sum
   res := zhu[0]
   for i := 1; i < 4; i++ {
       if zhu[i] < res {
           res = zhu[i]
       }
   }
   // prepare builders
   var sb [4]*strings.Builder
   for i := 0; i < 4; i++ {
       sb[i] = &strings.Builder{}
   }
   t := 0
   // error output and exit
   los := func() {
       fmt.Fprintln(writer, -1)
       writer.Flush()
       os.Exit(0)
   }
   // sing: balance to equal res
   sing := func(i, m int) {
       if m < 0 {
           los()
       }
       for j := 0; j < 4; j++ {
           if j == i {
               sb[j].WriteString(strings.Repeat("b", m))
           } else {
               sb[j].WriteString(strings.Repeat("a", m))
               zh[i][j] -= m
               zh[j][i] -= m
           }
       }
       t += m
   }
   // solve: final arrangement
   solve := func(i, m int) {
       if m < 0 {
           los()
       }
       for j := 0; j < 4; j++ {
           if j == 0 || j == i {
               sb[j].WriteString(strings.Repeat("b", m))
           } else {
               sb[j].WriteString(strings.Repeat("a", m))
           }
       }
       t += m
   }
   // check parity and perform sing
   for i := 0; i < 4; i++ {
       if (zhu[i]&1) != (res&1) {
           los()
       }
       sing(i, (zhu[i]-res)/2)
   }
   // final solve stage
   as := zh[0][1] + zh[0][2] + zh[0][3]
   if as&1 != 0 {
       los()
   }
   solve(3, (zh[0][1]+zh[0][2]-zh[0][3])/2)
   solve(2, (zh[0][1]+zh[0][3]-zh[0][2])/2)
   solve(1, (zh[0][3]+zh[0][2]-zh[0][1])/2)
   // output
   fmt.Fprintln(writer, t)
   for i := 0; i < 4; i++ {
       fmt.Fprintln(writer, sb[i].String())
   }
}
