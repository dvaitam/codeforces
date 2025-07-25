package main

import (
   "bufio"
   "fmt"
   "os"
)

// Point or vector
type Vec struct {
   x, y int64
   id   int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   pts := make([]Vec, n)
   for i := 0; i < n; i++ {
       var x, y int64
       fmt.Fscan(reader, &x, &y)
       pts[i] = Vec{x: x, y: y, id: i + 1}
   }
   var s string
   fmt.Fscan(reader, &s)
   // visited flags
   used := make([]bool, n)
   ans := make([]int, n)
   // find start: smallest x, then y
   start := 0
   for i := 1; i < n; i++ {
       if pts[i].x < pts[start].x || (pts[i].x == pts[start].x && pts[i].y < pts[start].y) {
           start = i
       }
   }
   // initialize
   used[start] = true
   ans[0] = pts[start].id
   curr := pts[start]
   // initial direction: positive x-axis
   prev := Vec{x: 1, y: 0}
   // greedy pick
   for k := 0; k < len(s); k++ {
       best := -1
       var bestVec Vec
       for i := 0; i < n; i++ {
           if used[i] {
               continue
           }
           // vector from curr to candidate
           v := Vec{x: pts[i].x - curr.x, y: pts[i].y - curr.y, id: pts[i].id}
           if best < 0 {
               best = i
               bestVec = v
               continue
           }
           // compare prev->bestVec (u) and prev->v
           u := bestVec
           // cross products of prev with u and v
           cdu := prev.x*u.y - prev.y*u.x
           cdv := prev.x*v.y - prev.y*v.x
           if s[k] == 'L' {
               // choose minimal positive angle -> smaller angle
               useV := false
               if (cdu >= 0) != (cdv >= 0) {
                   if cdv >= 0 {
                       useV = true
                   }
               } else {
                   // same half: compare cross(v, u) > 0 => v angle < u angle
                   if v.x*u.y - v.y*u.x > 0 {
                       useV = true
                   }
               }
               if useV {
                   best = i
                   bestVec = v
               }
           } else {
               // 'R': choose maximal angle -> larger angle
               useV := false
               if (cdu >= 0) != (cdv >= 0) {
                   if cdv < 0 {
                       useV = true
                   }
               } else {
                   // same half: compare cross(u, v) > 0 => v angle > u angle
                   if u.x*v.y - u.y*v.x > 0 {
                       useV = true
                   }
               }
               if useV {
                   best = i
                   bestVec = v
               }
           }
       }
       used[best] = true
       ans[k+1] = pts[best].id
       curr = pts[best]
       prev = bestVec
   }
   // last remaining point
   for i := 0; i < n; i++ {
       if !used[i] {
           ans[n-1] = pts[i].id
           break
       }
   }
   // output
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
