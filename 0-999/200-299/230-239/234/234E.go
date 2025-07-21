package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Team holds the name and rating of a football team
type Team struct {
   name   string
   rating int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, x, a, b, c int
   fmt.Fscan(reader, &n)
   fmt.Fscan(reader, &x, &a, &b, &c)

   teams := make([]Team, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &teams[i].name, &teams[i].rating)
   }
   // sort all teams by descending rating
   sort.Slice(teams, func(i, j int) bool {
       return teams[i].rating > teams[j].rating
   })

   m := n / 4
   // initialize baskets
   baskets := make([][]Team, 4)
   for i := 0; i < 4; i++ {
       baskets[i] = append([]Team(nil), teams[i*m:(i+1)*m]...)
   }

   // random number generator
   rng := func() int {
       x = (x*a + b) % c
       return x
   }

   // groups to hold the draw result
   groups := make([][]Team, m)
   // perform m-1 draws
   for gi := 0; gi < m-1; gi++ {
       var group []Team
       for bi := 0; bi < 4; bi++ {
           k := rng()
           idx := k % len(baskets[bi])
           group = append(group, baskets[bi][idx])
           // remove selected team from basket
           baskets[bi] = append(baskets[bi][:idx], baskets[bi][idx+1:]...)
       }
       // sort group by descending rating
       sort.Slice(group, func(i, j int) bool {
           return group[i].rating > group[j].rating
       })
       groups[gi] = group
   }
   // last group: remaining one in each basket
   var last []Team
   for bi := 0; bi < 4; bi++ {
       last = append(last, baskets[bi][0])
   }
   sort.Slice(last, func(i, j int) bool {
       return last[i].rating > last[j].rating
   })
   groups[m-1] = last

   // output groups
   for i, group := range groups {
       // group letter
       fmt.Fprintln(writer, string('A'+i))
       for _, t := range group {
           fmt.Fprintln(writer, t.name)
       }
   }
