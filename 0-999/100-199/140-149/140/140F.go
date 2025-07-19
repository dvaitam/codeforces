package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type point struct {
   x, y int
}
type pair struct {
   x, y int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   pts := make([]point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &pts[i].x, &pts[i].y)
   }
   if k >= n {
       fmt.Println(-1)
       return
   }
   sort.Slice(pts, func(i, j int) bool {
       if pts[i].x != pts[j].x {
           return pts[i].x < pts[j].x
       }
       return pts[i].y < pts[j].y
   })
   found := make(map[pair]bool)
   for ll := 0; ll <= k; ll++ {
       for rr := 0; ll+rr <= k; rr++ {
           sumIdx := ll + rr
           var cx, cy int
           if sumIdx == n-1 {
               cx = pts[ll].x*2
               cy = pts[ll].y*2
               found[pair{cx, cy}] = true
               continue
           }
           cx = pts[ll].x + pts[n-1-rr].x
           cy = pts[ll].y + pts[n-1-rr].y
           tmp := sumIdx
           ok := true
           l, r := ll, n-1-rr
           for l <= r {
               sx := pts[l].x + pts[r].x
               if sx != cx {
                   if tmp == k {
                       ok = false
                       break
                   }
                   tmp++
                   if sx > cx {
                       r--
                   } else {
                       l++
                   }
                   continue
               }
               sy := pts[l].y + pts[r].y
               if sy != cy {
                   if tmp == k {
                       ok = false
                       break
                   }
                   tmp++
                   if sy > cy {
                       r--
                   } else {
                       l++
                   }
                   continue
               }
               l++
               r--
           }
           if ok {
               found[pair{cx, cy}] = true
           }
       }
   }
   // collect and sort results
   res := make([]pair, 0, len(found))
   for p := range found {
       res = append(res, p)
   }
   sort.Slice(res, func(i, j int) bool {
       if res[i].x != res[j].x {
           return res[i].x < res[j].x
       }
       return res[i].y < res[j].y
   })
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(res))
   for _, p := range res {
       // print center coordinates with one decimal
       fmt.Fprintf(w, "%.1f %.1f\n", float64(p.x)/2.0, float64(p.y)/2.0)
   }
}
