package main

import (
   "bufio"
   "os"
   "sort"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // read n, m
   line, _ := reader.ReadBytes('\n')
   parts := splitInts(line)
   n, m := parts[0], parts[1]
   // read matrix
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       line, _ = reader.ReadBytes('\n')
       row := splitInts(line)
       for j := 0; j < m; j++ {
           a[i][j] = row[j]
       }
   }
   // row and column ranks
   b := make([][]int, n)
   for i := range b {
       b[i] = make([]int, m)
   }
   c := make([][]int, n)
   for i := range c {
       c[i] = make([]int, m)
   }
   f := make([]int, n)
   g := make([]int, m)
   type pair struct{ v, idx int }
   // rows
   for i := 0; i < n; i++ {
       pairs := make([]pair, m)
       for j := 0; j < m; j++ {
           pairs[j] = pair{a[i][j], j}
       }
       sort.Slice(pairs, func(i, j int) bool { return pairs[i].v < pairs[j].v })
       rank := 0
       for k := 0; k < m; k++ {
           if k == 0 || pairs[k].v != pairs[k-1].v {
               rank++
           }
           b[i][pairs[k].idx] = rank
       }
       f[i] = rank
   }
   // columns
   for j := 0; j < m; j++ {
       pairs := make([]pair, n)
       for i := 0; i < n; i++ {
           pairs[i] = pair{a[i][j], i}
       }
       sort.Slice(pairs, func(i, j int) bool { return pairs[i].v < pairs[j].v })
       rank := 0
       for k := 0; k < n; k++ {
           if k == 0 || pairs[k].v != pairs[k-1].v {
               rank++
           }
           c[pairs[k].idx][j] = rank
       }
       g[j] = rank
   }
   // output
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           ri := b[i][j]
           cj := c[i][j]
           t := ri
           if cj > t {
               t = cj
           }
           x1 := f[i] - ri
           x2 := g[j] - cj
           if x2 > x1 {
               x1 = x2
           }
           res := t + x1
           writer.WriteString(strconv.Itoa(res))
           if j+1 < m {
               writer.WriteByte(' ')
           }
       }
       writer.WriteByte('\n')
   }
}

// splitInts splits a line into integers, handling spaces
func splitInts(line []byte) []int {
   var res []int
   n := len(line)
   i := 0
   for i < n {
       // skip non-digit and non-minus
       for i < n && (line[i] < '0' && line[i] != '-' || line[i] > '9') {
           i++
       }
       if i >= n {
           break
       }
       sign := 1
       if line[i] == '-' {
           sign = -1
           i++
       }
       val := 0
       for i < n && line[i] >= '0' && line[i] <= '9' {
           val = val*10 + int(line[i]-'0')
           i++
       }
       res = append(res, sign*val)
   }
   return res
}
