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

   var n, iniPos int
   fmt.Fscan(in, &n, &iniPos)
   a := make([]int, n+1)
   vals := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
       vals = append(vals, a[i])
   }
   sort.Ints(vals)
   uniq := make([]int, 0, n)
   for i, v := range vals {
       if i == 0 || v != vals[i-1] {
           uniq = append(uniq, v)
       }
   }
   MAX := len(uniq)
   mp := make(map[int]int, MAX)
   for i, v := range uniq {
       mp[v] = i + 1
   }
   // buckets v[0]..v[MAX]
   v := make([][]int, MAX+1)
   v[0] = append(v[0], iniPos)
   for i := 1; i <= n; i++ {
       idx := mp[a[i]]
       v[idx] = append(v[idx], i)
   }
   // helper arrays
   en := make([]int, n+1)
   f := make([]int, n+1)
   same := make([]int, n+1)
   Same := make([]int, n+1)
   nxt := make([]int, n+1)

   // distance on circle
   dist := func(a, b int) int {
       d := a - b
       if d < 0 {
           d = -d
       }
       if d <= n-d {
           return d
       }
       return n - d
   }
   // prepare DP
   for m := MAX; m >= 0; m-- {
       if m < MAX {
           for _, i := range v[m] {
               t := i
               if m == 0 {
                   t = 0
               }
               en[t] = -1
               for _, ii := range v[m+1] {
                   d := dist(i, ii)
                   cost := d + f[ii]
                   if en[t] == -1 || en[t] > cost {
                       en[t] = cost
                       nxt[t] = ii
                   }
               }
           }
       }
       if m == 0 {
           break
       }
       s := len(v[m])
       for j, i := range v[m] {
           // forward neighbor
           ii := v[m][(j+1)%s]
           same[i] = 1
           Same[i] = ii
           // cost going forward
           var delta int
           if ii >= i {
               delta = ii - i
           } else {
               delta = ii - i + n
           }
           backCost := n - delta
           f[i] = en[ii] + backCost
           if ii == i {
               f[i] = en[ii]
               continue
           }
           // backward neighbor
           jj := (j-1+s)%s
           ii2 := v[m][jj]
           if i >= ii2 {
               delta = i - ii2
           } else {
               delta = i - ii2 + n
           }
           backCost2 := n - delta
           cost2 := en[ii2] + backCost2
           if f[i] > cost2 {
               f[i] = cost2
               same[i] = -1
               Same[i] = ii2
           }
       }
   }
   // move printing
   move := func(a, b int) {
       d := dist(a, b)
       var diff1 int
       if b >= a {
           diff1 = b - a
       } else {
           diff1 = b + n - a
       }
       if diff1 == d {
           fmt.Fprintf(out, "%+d\n", d)
       } else {
           fmt.Fprintf(out, "%+d\n", -d)
       }
   }
   // traverse group
   var a_a func(s, cur, typ int)
   a_a = func(s, cur, typ int) {
       arr := v[s]
       if typ == -1 {
           // reverse arr
           for l, r := 0, len(arr)-1; l < r; l, r = l+1, r-1 {
               arr[l], arr[r] = arr[r], arr[l]
           }
           v[s] = arr
       }
       size := len(arr)
       I := -1
       for i := 0; ; i = (i + 1) % size {
           t := arr[i]
           if I == -1 {
               if t == cur {
                   I = i
               }
               continue
           }
           if I == i {
               return
           }
           move(cur, t)
           cur = t
       }
   }
   // solve and output
   fmt.Fprintf(out, "%d\n", en[0])
   cur := iniPos
   s := 0
   for s < MAX {
       var t int
       if s == 0 {
           t = nxt[0]
       } else {
           t = nxt[cur]
       }
       move(cur, t)
       cur = t
       s++
       a_a(s, cur, -same[cur])
       cur = Same[cur]
   }
}
