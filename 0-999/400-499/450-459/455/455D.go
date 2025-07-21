package main

import (
   "bufio"
   "fmt"
   "os"
)

// sqrt-decomposition with dynamic blocks
func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   // read first int n
   var n int
   {
       fmt.Fscan(in, &n)
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // build blocks
   B := 400
   blocks := [][]int{}
   cnt := []map[int]int{}
   for i := 0; i < n; i += B {
       j := i + B
       if j > n {
           j = n
       }
       bl := make([]int, j-i)
       m := make(map[int]int)
       for k := i; k < j; k++ {
           bl[k-i] = a[k]
           m[a[k]]++
       }
       blocks = append(blocks, bl)
       cnt = append(cnt, m)
   }
   // process queries
   var q int
   fmt.Fscan(in, &q)
   lastans := 0
   for qi := 0; qi < q; qi++ {
       var t, li, ri, ki int
       fmt.Fscan(in, &t, &li, &ri)
       li = (li + lastans - 1) % n + 1
       ri = (ri + lastans - 1) % n + 1
       if li > ri {
           li, ri = ri, li
       }
       l := li - 1
       r := ri - 1
       if t == 1 {
           // cyclic shift right [l, r]
           // find blockR and posR
           sum := 0
           var blockR, posR int
           for bi := range blocks {
               ln := len(blocks[bi])
               if sum+ln > r {
                   blockR = bi
                   posR = r - sum
                   break
               }
               sum += ln
           }
           // get value
           v := blocks[blockR][posR]
           // remove at blockR
           cnt[blockR][v]--
           blr := blocks[blockR]
           blocks[blockR] = append(blr[:posR], blr[posR+1:]...)
           // find blockL and posL
           sum = 0
           var blockL, posL int
           for bi := range blocks {
               ln := len(blocks[bi])
               if sum+ln > l {
                   blockL = bi
                   posL = l - sum
                   break
               }
               sum += ln
           }
           // adjust posL if same block and removal was before
           if blockL == blockR && posR < posL {
               posL--
           }
           // insert v at blockL,posL
           cnt[blockL][v]++
           bll := blocks[blockL]
           // make space
           tmp := make([]int, len(bll)+1)
           copy(tmp, bll[:posL])
           tmp[posL] = v
           copy(tmp[posL+1:], bll[posL:])
           blocks[blockL] = tmp
       } else {
           // query count k
           fmt.Fscan(in, &ki)
           ki = (ki + lastans - 1) % n + 1
           k := ki
           ans := 0
           sum := 0
           for bi := range blocks {
               ln := len(blocks[bi])
               if sum+ln-1 < l {
                   sum += ln
                   continue
               }
               if sum > r {
                   break
               }
               // block range [sum, sum+ln-1]
               if sum >= l && sum+ln-1 <= r {
                   ans += cnt[bi][k]
               } else {
                   // partial
                   st := 0
                   if l > sum {
                       st = l - sum
                   }
                   en := ln - 1
                   if r < sum+ln-1 {
                       en = r - sum
                   }
                   for i := st; i <= en; i++ {
                       if blocks[bi][i] == k {
                           ans++
                       }
                   }
               }
               sum += ln
           }
           lastans = ans
           fmt.Fprintln(out, ans)
       }
   }
}
