package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

type Team struct {
   name   string
   points int
   gd      int
   gs      int
}

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanLines)
   // read n
   if !scanner.Scan() {
       return
   }
   n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
   if err != nil {
       return
   }
   teams := make([]*Team, n)
   idx := make(map[string]int, n)
   // read team names
   for i := 0; i < n && scanner.Scan(); i++ {
       name := strings.TrimSpace(scanner.Text())
       teams[i] = &Team{name: name}
       idx[name] = i
   }
   // process matches
   total := n*(n-1)/2
   for k := 0; k < total && scanner.Scan(); k++ {
       line := strings.TrimSpace(scanner.Text())
       // format: name1-name2 num1:num2
       parts := strings.SplitN(line, " ", 2)
       teamsPart := parts[0]
       scorePart := parts[1]
       names := strings.SplitN(teamsPart, "-", 2)
       name1, name2 := names[0], names[1]
       scores := strings.SplitN(scorePart, ":", 2)
       s1, _ := strconv.Atoi(scores[0])
       s2, _ := strconv.Atoi(scores[1])
       t1 := teams[idx[name1]]
       t2 := teams[idx[name2]]
       // update goals scored and goal difference
       t1.gs += s1
       t1.gd += s1 - s2
       t2.gs += s2
       t2.gd += s2 - s1
       // update points
       if s1 > s2 {
           t1.points += 3
       } else if s1 < s2 {
           t2.points += 3
       } else {
           t1.points++
           t2.points++
       }
   }
   // sort by points, gd, gs
   sort.Slice(teams, func(i, j int) bool {
       a, b := teams[i], teams[j]
       if a.points != b.points {
           return a.points > b.points
       }
       if a.gd != b.gd {
           return a.gd > b.gd
       }
       if a.gs != b.gs {
           return a.gs > b.gs
       }
       return false
   })
   // select top n/2 teams
   m := n / 2
   selected := make([]string, m)
   for i := 0; i < m; i++ {
       selected[i] = teams[i].name
   }
   sort.Strings(selected)
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for _, name := range selected {
       fmt.Fprintln(writer, name)
   }
}
