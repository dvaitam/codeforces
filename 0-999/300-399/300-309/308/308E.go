package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

// Sheep represents an interval and its original index
type Sheep struct {
   l, r, idx int
}

// ok checks if a sequence can be built with given med, and prints it if print is true
func ok(n, med int, print bool, v []Sheep, T [][]bool, buf *bufio.Writer) bool {
   ogr := make([]int, n)
   cnt := make([]int, n)
   pos := make([]int, n)
   on := make([]int, n)
   for i := 0; i < n; i++ {
       ogr[i] = n - 1
       cnt[i] = 0
       pos[i] = -1
   }
   cnt[n-1] = n
   for i := 0; i < n; i++ {
       sum := 0
       for j := 0; j < n; j++ {
           sum += cnt[j]
           if sum > j+1 {
               return false
           }
           if sum == j+1 && j >= i {
               break
           }
       }
       cur := -1
       for j := 0; j < n; j++ {
           if ogr[j] <= sum-1 && pos[j] == -1 {
               cur = j
               break
           }
       }
       cnt[ogr[cur]]--
       cnt[i]++
       ogr[cur] = i
       pos[cur] = i
       on[i] = cur
       for j := 0; j < n; j++ {
           if T[cur][j] {
               if i+med < ogr[j] {
                   cnt[ogr[j]]--
                   ogr[j] = i + med
                   cnt[ogr[j]]++
               }
           }
       }
   }
   if print {
       for i := 0; i < n; i++ {
           buf.WriteString(strconv.Itoa(v[on[i]].idx + 1))
           buf.WriteByte(' ')
       }
       buf.WriteByte('\n')
       buf.Flush()
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   v := make([]Sheep, n)
   for i := 0; i < n; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       v[i] = Sheep{l, r, i}
   }
   sort.Slice(v, func(i, j int) bool {
       return v[i].r < v[j].r
   })
   T := make([][]bool, n)
   for i := 0; i < n; i++ {
       T[i] = make([]bool, n)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if v[i].l > v[j].r || v[j].l > v[i].r {
               T[i][j] = false
           } else {
               T[i][j] = true
           }
       }
   }
   beg, end := 0, n
   for beg < end {
       med := (beg + end) / 2
       if ok(n, med, false, v, T, writer) {
           end = med
       } else {
           beg = med + 1
       }
   }
   ok(n, beg, true, v, T, writer)
}
