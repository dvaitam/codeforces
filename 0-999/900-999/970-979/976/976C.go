package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type item struct {
   w, h, idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]item, n)
   for i := 0; i < n; i++ {
       var w, h int
       fmt.Fscan(reader, &w, &h)
       // store negative height to sort descending by h for equal widths
       a[i] = item{w: w, h: -h, idx: i}
   }
   // sort by width asc, height desc (using stored negative h)
   sort.Slice(a, func(i, j int) bool {
       if a[i].w != a[j].w {
           return a[i].w < a[j].w
       }
       return a[i].h < a[j].h
   })
   // restore heights
   for i := range a {
       a[i].h = -a[i].h
   }
   // check for identical rectangles
   for i := 0; i+1 < n; i++ {
       if a[i].w == a[i+1].w && a[i].h == a[i+1].h {
           fmt.Fprintln(writer, a[i].idx+1, a[i+1].idx+1)
           return
       }
   }
   const INF = 1000000000
   mne := INF
   mnind := 0
   // scan from largest width to smallest
   for i := n - 1; i >= 0; i-- {
       if a[i].h >= mne {
           fmt.Fprintln(writer, a[mnind].idx+1, a[i].idx+1)
           return
       }
       if a[i].h < mne {
           mne = a[i].h
           mnind = i
       }
   }
   fmt.Fprintln(writer, "-1 -1")
}
