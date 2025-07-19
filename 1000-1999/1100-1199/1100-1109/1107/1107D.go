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
   rowsData := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &rowsData[i])
   }
   // initialize tables
   var tofail [16]bool
   for _, v := range []int{2, 4, 5, 6, 9, 10, 11, 13} {
       tofail[v] = true
   }
   var ps, pln, ss [16]int
   for _, v := range []int{0, 1, 3, 7} {
       ps[v] = 0
   }
   for _, v := range []int{8, 12, 14, 15} {
       ps[v] = 1
   }
   pln[0], pln[15] = 0, 0
   pln[1], pln[14] = 3, 3
   pln[3], pln[12] = 2, 2
   pln[7], pln[8] = 1, 1
   for _, v := range []int{0, 8, 12, 14} {
       ss[v] = 0
   }
   for _, v := range []int{1, 3, 7, 15} {
       ss[v] = 1
   }
   // compute gcd over row groups and column bit runs
   res := 0
   // row grouping
   cnt := 0
   for i := 0; i < n; i++ {
       cnt++
       if i == n-1 || rowsData[i] != rowsData[i+1] {
           res = gcd(res, cnt)
           cnt = 0
       }
   }
   // process each row for bit runs
   for i := 0; i < n; i++ {
       row := rowsData[i]
       length := 0
       // initial state from first hex digit
       first := hexVal(row[0])
       sn := ps[first]
       segs := len(row)
       for j := 0; j < segs; j++ {
           cs := hexVal(row[j])
           if tofail[cs] {
               fmt.Println(1)
               return
           }
           if cs == 0 && sn == 0 {
               length += 4
               continue
           }
           if cs == 15 && sn == 1 {
               length += 4
               continue
           }
           if sn == ps[cs] {
               length += pln[cs]
           }
           res = gcd(res, length)
           length = 4 - pln[cs]
           sn = ss[cs]
       }
       res = gcd(res, length)
   }
   fmt.Println(res)
}

func gcd(a, b int) int {
   if a == 0 {
       return b
   }
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func hexVal(c byte) int {
   if c >= '0' && c <= '9' {
       return int(c - '0')
   }
   return int(c - 'A' + 10)
}
