package main

import (
   "bufio"
   "os"
   "sort"
   "strconv"
)

func readInt(r *bufio.Reader) int {
   var x int
   sign := 1
   b, err := r.ReadByte()
   for err == nil && b != '-' && (b < '0' || b > '9') {
       b, err = r.ReadByte()
   }
   if b == '-' {
       sign = -1
       b, _ = r.ReadByte()
   }
   for ; err == nil && b >= '0' && b <= '9'; b, err = r.ReadByte() {
       x = x*10 + int(b - '0')
   }
   return x * sign
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   n := readInt(reader)
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       a[i] = readInt(reader)
       b[i] = readInt(reader)
   }
   s1 := make([]int, n)
   s2 := make([]int, n)
   for i := 0; i < n; i++ {
       s1[i] = i
       s2[i] = i
   }
   sort.Slice(s1, func(i, j int) bool { return a[s1[i]] < a[s1[j]] })
   sort.Slice(s2, func(i, j int) bool { return b[s2[i]] < b[s2[j]] })
   x := make([]int, n)
   y := make([]int, n)
   c := make([]int, 0, n*2)
   l, r := 0, 0
   for l < n || r < n {
       var tmp int
       if l < n && r < n {
           if a[s1[l]] < b[s2[r]] {
               tmp = a[s1[l]]
           } else {
               tmp = b[s2[r]]
           }
       } else if l < n {
           tmp = a[s1[l]]
       } else {
           tmp = b[s2[r]]
       }
       c = append(c, tmp)
       idx := len(c) - 1
       for l < n && a[s1[l]] == tmp {
           x[s1[l]] = idx
           l++
       }
       for r < n && b[s2[r]] == tmp {
           y[s2[r]] = idx
           r++
       }
   }
   m := len(c)
   f := make([]int, m)
   d := make([]int, m)
   for i := 0; i < m; i++ {
       f[i] = i
   }
   // find with path compression
   var find func(int) int
   find = func(u int) int {
       for f[u] != u {
           f[u] = f[f[u]]
           u = f[u]
       }
       return u
   }
   // union and count
   for i := 0; i < n; i++ {
       xi := x[i]
       yi := y[i]
       fi := find(xi)
       fj := find(yi)
       if fi == fj {
           d[fi]++
       } else {
           d[fj] += d[fi] + 1
           f[fi] = fj
       }
   }
   // flatten
   for i := 0; i < m; i++ {
       f[i] = find(i)
   }
   // decrement
   for i := 0; i < m; i++ {
       d[f[i]]--
   }
   // check impossibility
   for i := 0; i < m; i++ {
       if f[i] == i && d[i] > 0 {
           writer.WriteString("-1")
           return
       }
   }
   v := make([]bool, m)
   for i := 0; i < m; i++ {
       if f[i] == i && d[i] == 0 {
           v[i] = true
       }
   }
   // find answer
   for i := m - 1; i >= 0; i-- {
       root := f[i]
       if !v[root] {
           v[root] = true
           continue
       }
       writer.WriteString(strconv.Itoa(c[i]))
       return
   }
}
