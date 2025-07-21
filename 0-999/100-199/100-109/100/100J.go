package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Interval represents an interval [l, r]
type Interval struct {
   l, r int
}

// Clique maintains union and common intersection of intervals in the clique
type Clique struct {
   unionL, unionR int
   comL, comR     int
}

// Color holds a set of cliques for this color class
type Color struct {
   cliques []Clique
}

// canUse checks if interval x can be added to this color without creating a P3
func (c *Color) canUse(x Interval) bool {
   count := 0
   for _, cl := range c.cliques {
       // check if x intersects any interval in clique via union interval
       if x.r < cl.unionL || x.l > cl.unionR {
           // no intersection with any in this clique
           continue
       }
       // intersects at least one; now check if it intersects all via common region
       if x.r >= cl.comL && x.l <= cl.comR {
           // intersects common region => can join this clique
           count++
           if count > 1 {
               return false
           }
       } else {
           // intersects some but not all => would form P3
           return false
       }
   }
   return true
}

// add inserts interval x into appropriate clique or creates a new one
func (c *Color) add(x Interval) {
   for i := range c.cliques {
       cl := &c.cliques[i]
       // check common intersection
       if x.r >= cl.comL && x.l <= cl.comR {
           // join this clique
           if x.l < cl.unionL {
               cl.unionL = x.l
           }
           if x.r > cl.unionR {
               cl.unionR = x.r
           }
           if x.l > cl.comL {
               cl.comL = x.l
           }
           if x.r < cl.comR {
               cl.comR = x.r
           }
           return
       }
   }
   // no joinable clique: create new
   c.cliques = append(c.cliques, Clique{
       unionL: x.l, unionR: x.r,
       comL: x.l, comR: x.r,
   })
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   intervals := make([]Interval, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &intervals[i].l, &intervals[i].r)
   }
   // sort by start ascending, then end ascending
   sort.Slice(intervals, func(i, j int) bool {
       if intervals[i].l != intervals[j].l {
           return intervals[i].l < intervals[j].l
       }
       return intervals[i].r < intervals[j].r
   })
   var colors []Color
   // assign intervals
   for _, iv := range intervals {
       placed := false
       for ci := range colors {
           if colors[ci].canUse(iv) {
               colors[ci].add(iv)
               placed = true
               break
           }
       }
       if !placed {
           // new color
           var c Color
           c.cliques = []Clique{{unionL: iv.l, unionR: iv.r, comL: iv.l, comR: iv.r}}
           colors = append(colors, c)
       }
   }
   // output number of colors
   fmt.Println(len(colors))
}
