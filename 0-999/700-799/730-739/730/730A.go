package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   type pair struct {
       val int
       idx int
   }
   r := make([]pair, n)
   sum := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &r[i].val)
       r[i].idx = i
       sum += r[i].val
   }
   // initialize strings s and g
   s := make([]byte, n)
   g := make([]byte, n)
   for i := 0; i < n; i++ {
       s[i] = '0'
       g[i] = '0'
   }
   // moves as pairs of indices
   var moves [][2]int
   cnt := 0
   ch := false
   // process
   for {
       sort.Slice(r, func(i, j int) bool { return r[i].val < r[j].val })
       if r[0].val == r[n-1].val {
           break
       }
       cnt++
       minv := r[0].val
       maxv := r[n-1].val
       // special case
       if maxv == minv+1 && sum-n*minv == 3 {
           ch = true
           // mark top three (last three largest)
           for k := n - 1; k >= n-3; k-- {
               g[r[k].idx] = '1'
           }
           break
       }
       // record a move: decrement two largest
       i1 := r[n-1].idx
       i2 := r[n-2].idx
       moves = append(moves, [2]int{i1, i2})
       // decrement values
       if r[n-1].val > 0 {
           r[n-1].val--
           sum--
       }
       if r[n-2].val > 0 {
           r[n-2].val--
           sum--
       }
   }
   // output final equal value and move count
   fmt.Fprintln(writer, r[0].val)
   fmt.Fprintln(writer, cnt)
   if ch {
       writer.Write(g)
       writer.WriteByte('\n')
   }
   // print moves: for each, set s[i]=s[j]='1', print, reset
   for _, mv := range moves {
       i, j := mv[0], mv[1]
       s[i] = '1'
       s[j] = '1'
       writer.Write(s)
       writer.WriteByte('\n')
       s[i] = '0'
       s[j] = '0'
   }
}
