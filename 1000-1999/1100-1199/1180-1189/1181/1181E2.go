package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Rect struct {
   a, b, c, d int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   rects := make([]Rect, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &rects[i].a, &rects[i].b, &rects[i].c, &rects[i].d)
   }
   ids0 := make([]int, n)
   for i := 0; i < n; i++ {
       ids0[i] = i
   }
   ok := true
   // stack of segments to check
   stack := [][]int{ids0}
   for len(stack) > 0 {
       // pop
       ids := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       m := len(ids)
       if m <= 1 {
           continue
       }
       // try vertical cut
       idsA := make([]int, m)
       copy(idsA, ids)
       sort.Slice(idsA, func(i, j int) bool {
           return rects[idsA[i]].a < rects[idsA[j]].a
       })
       maxC := rects[idsA[0]].c
       split := -1
       for i := 0; i < m-1; i++ {
           if rects[idsA[i]].c > maxC {
               maxC = rects[idsA[i]].c
           }
           // next a
           if maxC <= rects[idsA[i+1]].a {
               split = i
               break
           }
       }
       if split >= 0 {
           // two parts: [0..split], [split+1..]
           left := make([]int, split+1)
           right := make([]int, m-split-1)
           copy(left, idsA[:split+1])
           copy(right, idsA[split+1:])
           stack = append(stack, left, right)
           continue
       }
       // try horizontal cut
       idsB := make([]int, m)
       copy(idsB, ids)
       sort.Slice(idsB, func(i, j int) bool {
           return rects[idsB[i]].b < rects[idsB[j]].b
       })
       maxD := rects[idsB[0]].d
       for i := 0; i < m-1; i++ {
           if rects[idsB[i]].d > maxD {
               maxD = rects[idsB[i]].d
           }
           if maxD <= rects[idsB[i+1]].b {
               split = i
               break
           }
       }
       if split >= 0 {
           left := make([]int, split+1)
           right := make([]int, m-split-1)
           copy(left, idsB[:split+1])
           copy(right, idsB[split+1:])
           stack = append(stack, left, right)
           continue
       }
       ok = false
       break
   }
   if ok {
       fmt.Fprintln(out, "YES")
   } else {
       fmt.Fprintln(out, "NO")
   }
}
