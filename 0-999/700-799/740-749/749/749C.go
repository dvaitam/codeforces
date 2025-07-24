package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   s := make([]byte, n)
   // read the string of fractions
   // skip whitespace then read
   var line string
   if _, err := fmt.Fscan(reader, &line); err != nil {
       return
   }
   for i := 0; i < n; i++ {
       s[i] = line[i]
   }

   // initialize queues with capacity 2*n
   total := n * 2
   dq := make([]int, total)
   rq := make([]int, total)
   hd, td := 0, 0
   hr, tr := 0, 0
   for i := 0; i < n; i++ {
       if s[i] == 'D' {
           dq[td] = i
           td++
       } else {
           rq[tr] = i
           tr++
       }
   }

   // simulate banning process
   for hd < td && hr < tr {
       di := dq[hd]
       hd++
       ri := rq[hr]
       hr++
       if di < ri {
           // D bans R
           dq[td] = di + n
           td++
       } else {
           // R bans D
           rq[tr] = ri + n
           tr++
       }
   }
   if hd < td {
       fmt.Println("D")
   } else {
       fmt.Println("R")
   }
}
