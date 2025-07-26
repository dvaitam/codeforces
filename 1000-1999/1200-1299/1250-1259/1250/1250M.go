package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   // block size: at most 50 blocks
   b := (n + 50 - 1) / 50
   if b <= 0 {
       b = 1
   }
   k := (n + b - 1) / b
   // prepare commands
   var rowsList [][]int
   var colsList [][]int
   for t := 0; t < k; t++ {
       start := t*b + 1
       end := (t + 1) * b
       if end > n {
           end = n
       }
       // rows in this block
       rows := make([]int, 0, end-start+1)
       for i := start; i <= end; i++ {
           rows = append(rows, i)
       }
       // mark excluded columns (slabs) for these rows
       excluded := make([]bool, n+2)
       for _, i := range rows {
           // main diagonal
           excluded[i] = true
           // slab below diagonal
           if i >= 2 {
               excluded[i-1] = true
           }
       }
       // build columns list avoiding excluded
       cols := make([]int, 0, n)
       for j := 1; j <= n; j++ {
           if !excluded[j] {
               cols = append(cols, j)
           }
       }
       if len(cols) == 0 {
           continue
       }
       rowsList = append(rowsList, rows)
       colsList = append(colsList, cols)
   }
   // output
   fmt.Fprintln(out, len(rowsList))
   for idx := range rowsList {
       rows := rowsList[idx]
       cols := colsList[idx]
       fmt.Fprint(out, len(rows))
       for _, v := range rows {
           fmt.Fprint(out, " ", v)
       }
       fmt.Fprintln(out)
       fmt.Fprint(out, len(cols))
       for _, v := range cols {
           fmt.Fprint(out, " ", v)
       }
       fmt.Fprintln(out)
   }
