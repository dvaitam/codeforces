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
   var W, N, M int
   if _, err := fmt.Fscan(reader, &W, &N, &M); err != nil {
       return
   }
   size := 1 << W
   nc := make([]int, size)
   for i := 0; i < N; i++ {
       var v int
       fmt.Fscan(reader, &v)
       nc[v]++
   }
   // conv: map mask to base-3 value
   conv := make([]int, size)
   for mask := 0; mask < size; mask++ {
       v := 0
       for j := W - 1; j >= 0; j-- {
           v *= 3
           if mask&(1<<j) != 0 {
               v++
           }
       }
       conv[mask] = v
   }
   // prepare nv counts
   // max base-3 value is 3^W - 1
   pw3 := 1
   for i := 0; i < W; i++ {
       pw3 *= 3
   }
   nv := make([]int64, pw3)
   for i := 0; i < size; i++ {
       ci := nc[i]
       if ci == 0 {
           continue
       }
       for j := 0; j < size; j++ {
           cj := nc[j]
           if cj == 0 {
               continue
           }
           idx := conv[i] + conv[j]
           nv[idx] += int64(ci) * int64(cj)
       }
   }
   // process queries
   buf := make([]byte, W)
   for qi := 0; qi < M; qi++ {
       // read string of length W
       var s string
       fmt.Fscan(reader, &s)
       // res(s)
       // vc dynamic slice
       var vc []int
       nr := 0
       cs := 0
       // helper pchange
       pchange := func(x int) {
           newNr := 2*nr + 1
           newVc := make([]int, newNr)
           // copy old
           for i := 0; i < nr; i++ {
               newVc[i] = vc[i]
           }
           // insert x
           newVc[nr] = x
           // copy old again
           for i := 0; i < nr; i++ {
               newVc[i+nr+1] = vc[i]
           }
           if nr > 0 {
               sloc := nr / 2
               newVc[sloc+nr+1] = -vc[sloc]
           }
           vc = newVc
           nr = newNr
       }
       // build initial cs and vc
       mul := 1
       for pi := W - 1; pi >= 0; pi-- {
           ch := s[pi]
           switch ch {
           case 'A':
               pchange(mul)
           case 'O':
               // nothing
           case 'X':
               pchange(2 * mul)
           case 'a':
               cs += 2 * mul
           case 'o':
               cs += mul
               pchange(mul)
           case 'x':
               cs += mul
           }
           mul *= 3
       }
       // accumulate answer
       var ans int64
       ans = nv[cs]
       for i := 0; i < nr; i++ {
           cs += vc[i]
           ans += nv[cs]
       }
       fmt.Fprintln(writer, ans)
   }
}
