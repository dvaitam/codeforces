package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Rect represents a rectangle with sorted sides a >= b >= c and its original id.
type Rect struct {
   a, b, c int
   id      int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   rects := make([]Rect, n)
   for i := 0; i < n; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       // sort sides so that a >= b >= c
       if a < b {
           a, b = b, a
       }
       if a < c {
           a, c = c, a
       }
       if b < c {
           b, c = c, b
       }
       rects[i] = Rect{a: a, b: b, c: c, id: i + 1}
   }

   // best represents the maximum diameter (2 * radius)
   best := 0
   p1 := 1
   for _, r := range rects {
       if r.c > best {
           best = r.c
           p1 = r.id
       }
   }

   // sort rectangles by a desc, then b desc, then c desc
   sort.Slice(rects, func(i, j int) bool {
       if rects[i].a != rects[j].a {
           return rects[i].a > rects[j].a
       }
       if rects[i].b != rects[j].b {
           return rects[i].b > rects[j].b
       }
       return rects[i].c > rects[j].c
   })

   flag := false
   var p2 int
   // check pairs with same a and b
   for i := 1; i < n; i++ {
       if rects[i].a == rects[i-1].a && rects[i].b == rects[i-1].b {
           sumC := rects[i].c + rects[i-1].c
           // diameter limited by b
           candidate := sumC
           if candidate > rects[i].b {
               candidate = rects[i].b
           }
           if candidate > best {
               best = candidate
               p1 = rects[i].id
               p2 = rects[i-1].id
               flag = true
           }
       }
   }

   if !flag {
       fmt.Fprintln(writer, 1)
       fmt.Fprintln(writer, p1)
   } else {
       fmt.Fprintln(writer, 2)
       fmt.Fprintln(writer, p1, p2)
   }
}
