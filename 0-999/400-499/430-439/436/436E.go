package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Node struct {
   x   int64
   id  int
   id2 int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   readInt := func() int64 {
       var x int64
       var neg bool
       b, _ := reader.ReadByte()
       for (b < '0' || b > '9') && b != '-' {
           b, _ = reader.ReadByte()
       }
       if b == '-' {
           neg = true
           b, _ = reader.ReadByte()
       }
       for b >= '0' && b <= '9' {
           x = x*10 + int64(b-'0')
           b, _ = reader.ReadByte()
       }
       if neg {
           return -x
       }
       return x
   }
   n := int(readInt())
   m := int(readInt())
   a := make([]int64, n+1)
   b := make([]int64, n+1)
   t := make([]Node, 2*n)
   for i := 1; i <= n; i++ {
       ai := readInt() << 1
       bi := readInt() << 1
       a[i] = ai
       b[i] = bi
       // two candidates
       if bi-ai < ai {
           t[i-1] = Node{bi >> 1, i, i << 1}
           t[i-1+n] = Node{bi >> 1, i + n, i<<1 | 1}
       } else {
           t[i-1] = Node{ai, i, i << 1}
           t[i-1+n] = Node{bi - ai, i + n, i<<1 | 1}
       }
   }
   sort.Slice(t, func(i, j int) bool {
       if t[i].x != t[j].x {
           return t[i].x < t[j].x
       }
       if t[i].id2 != t[j].id2 {
           return t[i].id2 < t[j].id2
       }
       return t[i].id < t[j].id
   })
   cnt := make([]int, 2*n+2)
   var ans int64
   for i := 0; i < m; i++ {
       ans += t[i].x
       cnt[t[i].id]++
   }
   // possible improvement
   if t[m-1].id <= n && cnt[t[m-1].id+n] == 0 {
       id0 := t[m-1].id
       if t[m-1].x != a[id0] {
           tmp := ans - t[m-1].x
           ans = tmp + a[id0]
           cnt[id0]--
           // collect changes
           C := []int{id0}
           Cval := []int{1}
           // search for best swap
           for i := 1; i <= n; i++ {
               if i == id0 {
                   continue
               }
               // case: not selected
               if cnt[i] == 0 {
                   if a[i]+tmp < ans {
                       ans = a[i] + tmp
                       C = []int{i}
                       Cval = []int{1}
                   }
               }
               // case: selected first, not second
               if cnt[i] > 0 && cnt[i+n] == 0 {
                   if tmp + (b[i]-a[i]) < ans {
                       ans = tmp + (b[i] - a[i])
                       C = []int{i + n}
                       Cval = []int{1}
                   }
                   if tmp - a[i] + b[id0] < ans {
                       ans = tmp - a[i] + b[id0]
                       C = []int{i, id0, id0 + n}
                       Cval = []int{-1, 1, 1}
                   }
               }
               // case: both selected
               if cnt[i] > 0 && cnt[i+n] > 0 {
                   if tmp - (b[i] - a[i]) + b[id0] < ans {
                       ans = tmp - (b[i] - a[i]) + b[id0]
                       C = []int{i + n, id0, id0 + n}
                       Cval = []int{-1, 1, 1}
                   }
               }
           }
           // apply changes
           for k, idx := range C {
               cnt[idx] += Cval[k]
           }
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans>>1)
   // per item selection count
   for i := 1; i <= n; i++ {
       sel := cnt[i] + cnt[i+n]
       // print without separator
       writer.WriteByte(byte('0' + sel))
   }
   writer.WriteByte('\n')
}
