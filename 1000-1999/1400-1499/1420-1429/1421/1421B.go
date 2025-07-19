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
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       grid := make([]string, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &grid[i])
       }
       // positions: (1,2), (2,1), (n,n-1), (n-1,n)
       a := grid[0][1]
       b := grid[1][0]
       x := grid[n-1][n-2]
       y := grid[n-2][n-1]
       type pair struct{ r, c int }
       opt1 := make([]pair, 0, 4) // make a=b='0', x=y='1'
       if a != '0' {
           opt1 = append(opt1, pair{1, 2})
       }
       if b != '0' {
           opt1 = append(opt1, pair{2, 1})
       }
       if x != '1' {
           opt1 = append(opt1, pair{n, n - 1})
       }
       if y != '1' {
           opt1 = append(opt1, pair{n - 1, n})
       }
       opt2 := make([]pair, 0, 4) // make a=b='1', x=y='0'
       if a != '1' {
           opt2 = append(opt2, pair{1, 2})
       }
       if b != '1' {
           opt2 = append(opt2, pair{2, 1})
       }
       if x != '0' {
           opt2 = append(opt2, pair{n, n - 1})
       }
       if y != '0' {
           opt2 = append(opt2, pair{n - 1, n})
       }
       // choose minimal
       res := opt1
       if len(opt2) < len(opt1) {
           res = opt2
       }
       // output
       fmt.Fprintln(writer, len(res))
       for _, p := range res {
           fmt.Fprintln(writer, p.r, p.c)
       }
   }
}
