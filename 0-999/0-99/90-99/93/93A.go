package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, a, b int64
   if _, err := fmt.Fscan(reader, &n, &m, &a, &b); err != nil {
       return
   }
   // rows and columns are 0-based for rows, 1-based for columns
   ra := (a - 1) / m
   rb := (b - 1) / m
   ca := (a-1)%m + 1
   cb := (b-1)%m + 1
   var res int64
   // if in same row, one rectangle suffice
   if ra == rb {
       res = 1
   } else if ca == 1 && cb == m {
       // full rows from a to b
       res = 1
   } else if ra+1 == rb {
       // adjacent rows: tail of first and head of second
       res = 2
   } else {
       // general: first row tail, full middle, last row head
       res = 3
   }
   fmt.Println(res)
}
