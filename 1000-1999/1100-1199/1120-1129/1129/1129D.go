package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353
const szb = 330
const cntb = 330

func addInt(x *int32, y int32) {
   *x += y
   if *x >= mod {
       *x -= mod
   }
}

func subInt(x *int32, y int32) {
   *x -= y
   if *x < 0 {
       *x += mod
   }
}

type pair struct { l, r int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   dp := make([]int32, n+1)
   dp[0] = 1
   cnt := make([]int32, n+2)
   freq := make([][]int32, cntb)
   for i := range freq {
       freq[i] = make([]int32, n+2)
   }
   global := make([]int32, cntb)
   total := int32(0)

   pos := make(map[int][]pair)

   // helper functions
   incseg := func(l, r int) {
       L := l / szb
       for i := l; i <= r; i++ {
           subInt(&freq[L][cnt[i]], dp[i-1])
           if cnt[i]+global[L] == int32(k) {
               subInt(&total, dp[i-1])
           }
           cnt[i]++
           addInt(&freq[L][cnt[i]], dp[i-1])
       }
   }
   decseg := func(l, r int) {
       L := l / szb
       for i := l; i <= r; i++ {
           subInt(&freq[L][cnt[i]], dp[i-1])
           cnt[i]--
           addInt(&freq[L][cnt[i]], dp[i-1])
           if cnt[i]+global[L] == int32(k) {
               addInt(&total, dp[i-1])
           }
       }
   }
   inc := func(l, r int) {
       L := l / szb
       R := r / szb
       if L == R {
           incseg(l, r)
       } else {
           incseg(l, szb*(L+1)-1)
           incseg(szb*R, r)
           for i := L + 1; i < R; i++ {
               j := int32(k) - global[i]
               if j >= 0 {
                   subInt(&total, freq[i][j])
               }
               global[i]++
           }
       }
   }
   dec := func(l, r int) {
       L := l / szb
       R := r / szb
       if L == R {
           decseg(l, r)
       } else {
           decseg(l, szb*(L+1)-1)
           decseg(szb*R, r)
           for i := L + 1; i < R; i++ {
               global[i]--
               j := int32(k) - global[i]
               if j >= 0 {
                   addInt(&total, freq[i][j])
               }
           }
       }
   }

   for i := 1; i <= n; i++ {
       L := i / szb
       // initialize this position
       addInt(&freq[L][0], dp[i-1])
       addInt(&total, dp[i-1])
       x := a[i]
       lst := pos[x]
       if len(lst) == 0 {
           inc(1, i)
           pos[x] = append(pos[x], pair{1, i})
       } else {
           last := lst[len(lst)-1]
           dec(last.l, last.r)
           inc(last.r+1, i)
           pos[x] = append(pos[x], pair{last.r + 1, i})
       }
       dp[i] = total
   }
   fmt.Fprintln(writer, dp[n])
}
