package main

import (
   "bufio"
   "fmt"
   "os"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   type pair struct{ pos, val int }
   adds := make([]pair, 0, n)
   lens := make([]int, 0, n)
   total := 0
   for len(a) > 0 {
       // find matching element
       j := -1
       for i := 1; i < len(a); i++ {
           if a[i] == a[0] {
               j = i
               break
           }
       }
       if j == -1 {
           fmt.Fprintln(writer, -1)
           return
       }
       // record additions
       for k := 1; k < j; k++ {
           adds = append(adds, pair{total + k + j, a[k]})
       }
       // prepare next array
       // reverse a[1:j]
       b := make([]int, 0, len(a)-2)
       for x := j - 1; x >= 1; x-- {
           b = append(b, a[x])
       }
       // append remaining
       if j+1 < len(a) {
           b = append(b, a[j+1:]...)
       }
       a = b
       lens = append(lens, j)
       total += j * 2
   }
   // output
   fmt.Fprintln(writer, len(adds))
   for _, p := range adds {
       fmt.Fprintln(writer, p.pos, p.val)
   }
   fmt.Fprintln(writer, len(lens))
   for _, l := range lens {
       fmt.Fprintln(writer, l*2)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var T int
   fmt.Fscan(reader, &T)
   for t := 0; t < T; t++ {
       solve(reader, writer)
   }
}
