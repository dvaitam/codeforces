package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var s string
   fmt.Fscan(in, &s)
   cur := []byte(s)
   var ops [][]int

   for {
       n := len(cur)
       i, j := 0, n-1
       var pos []int
       for {
           for i < j && cur[i] != '(' {
               i++
           }
           for i < j && cur[j] != ')' {
               j--
           }
           if i >= j {
               break
           }
           pos = append(pos, i)
           pos = append(pos, j)
           i++
           j--
       }
       if len(pos) == 0 {
           break
       }
       sort.Ints(pos)
       // record operation indices (1-based)
       op := make([]int, len(pos))
       for k, p := range pos {
           op[k] = p + 1
       }
       ops = append(ops, op)
       // remove selected positions
       removed := make([]bool, n)
       for _, p := range pos {
           removed[p] = true
       }
       newCur := make([]byte, 0, n-len(pos))
       for idx, b := range cur {
           if !removed[idx] {
               newCur = append(newCur, b)
           }
       }
       cur = newCur
   }

   // output
   k := len(ops)
   fmt.Fprintln(out, k)
   for _, op := range ops {
       m := len(op)
       fmt.Fprintln(out, m)
       for i, v := range op {
           if i > 0 {
               fmt.Fprint(out, " ")
           }
           fmt.Fprint(out, v)
       }
       fmt.Fprintln(out)
   }
}
