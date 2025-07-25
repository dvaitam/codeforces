package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// DSU for union-find
type dsu struct {
   p []int
}
func newDSU(n int) *dsu {
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
   }
   return &dsu{p}
}
func (d *dsu) find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.find(d.p[x])
   }
   return d.p[x]
}
func (d *dsu) union(a, b int) {
   ra := d.find(a)
   rb := d.find(b)
   if ra != rb {
       d.p[rb] = ra
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   d := newDSU(n)
   // use random perfect matchings over 7 rounds
   base := make([]int, n)
   for i := 0; i < n; i++ {
       base[i] = i + 1
   }
   rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
   // temp slice
   arr := make([]int, n)
   for r := 0; r < 7; r++ {
       copy(arr, base)
       rnd.Shuffle(n, func(i, j int) {
           arr[i], arr[j] = arr[j], arr[i]
       })
       m := n / 2
       // query batch
       fmt.Fprintf(writer, "? %d", m)
       for i := 0; i < m; i++ {
           a := arr[2*i]
           b := arr[2*i+1]
           fmt.Fprintf(writer, " %d %d", a, b)
       }
       fmt.Fprintln(writer)
       writer.Flush()
       // read responses
       for i := 0; i < m; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if x == 1 {
               d.union(arr[2*i], arr[2*i+1])
           }
       }
   }
   // collect groups
   groupsMap := make(map[int][]int)
   for i := 1; i <= n; i++ {
       root := d.find(i)
       groupsMap[root] = append(groupsMap[root], i)
   }
   // output final piles (exactly three lines)
   fmt.Fprintln(writer, "!")
   // collect groups
   groups := make([][]int, 0, 3)
   for _, grp := range groupsMap {
       groups = append(groups, grp)
   }
   // pad to three piles if needed
   for len(groups) < 3 {
       groups = append(groups, []int{})
   }
   for i := 0; i < 3; i++ {
       grp := groups[i]
       if len(grp) == 0 {
           fmt.Fprintln(writer)
           continue
       }
       for j, v := range grp {
           if j > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, v)
       }
       fmt.Fprintln(writer)
   }
}
