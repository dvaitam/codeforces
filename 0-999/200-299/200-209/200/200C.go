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
   name      string
   points    int
   scored    int
   conceded  int
}

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   teams := make(map[string]*Team)
   gamesCount := make(map[string]int)
   var lines []string
   for i := 0; i < 5 && scanner.Scan(); i++ {
       lines = append(lines, scanner.Text())
   }
   // process known games
   for _, line := range lines {
       parts := strings.Fields(line)
       if len(parts) != 3 {
           continue
       }
       t1, t2 := parts[0], parts[1]
       score := parts[2]
       ss := strings.Split(score, ":")
       g1, _ := strconv.Atoi(ss[0])
       g2, _ := strconv.Atoi(ss[1])
       // init teams
       if teams[t1] == nil {
           teams[t1] = &Team{name: t1}
       }
       if teams[t2] == nil {
           teams[t2] = &Team{name: t2}
       }
       gamesCount[t1]++
       gamesCount[t2]++
       // update scored/conceded
       teams[t1].scored += g1
       teams[t1].conceded += g2
       teams[t2].scored += g2
       teams[t2].conceded += g1
       // update points
       if g1 > g2 {
           teams[t1].points += 3
       } else if g1 < g2 {
           teams[t2].points += 3
       } else {
           teams[t1].points++
           teams[t2].points++
       }
   }
   // find opponent of BERLAND
   opp := ""
   for name, cnt := range gamesCount {
       if name != "BERLAND" && cnt == 2 {
           opp = name
       }
   }
   // original stats copy
   orig := make(map[string]Team)
   for name, t := range teams {
       orig[name] = Team{name: name, points: t.points, scored: t.scored, conceded: t.conceded}
   }
   // search minimal difference and goals conceded
   found := false
   bestX, bestY := 0, 0
   // iterate difference D = X-Y
   const limit = 1000
   for d := 1; d <= limit && !found; d++ {
       for y := 0; y <= limit; y++ {
           x := y + d
           // simulate copy
           sim := make([]Team, 0, len(orig))
           for _, t := range orig {
               sim = append(sim, t)
           }
           // update BERLAND
           for i := range sim {
               if sim[i].name == "BERLAND" {
                   sim[i].points += 3
                   sim[i].scored += x
                   sim[i].conceded += y
               }
               if sim[i].name == opp {
                   sim[i].scored += y
                   sim[i].conceded += x
               }
           }
           // sort by ranking rules
           sort.Slice(sim, func(i, j int) bool {
               a, b := sim[i], sim[j]
               if a.points != b.points {
                   return a.points > b.points
               }
               da := a.scored - a.conceded
               db := b.scored - b.conceded
               if da != db {
                   return da > db
               }
               if a.scored != b.scored {
                   return a.scored > b.scored
               }
               return a.name < b.name
           })
           // find BERLAND rank
           rank := -1
           for i, t := range sim {
               if t.name == "BERLAND" {
                   rank = i
                   break
               }
           }
           if rank >= 0 && rank < 2 {
               bestX, bestY = x, y
               found = true
               break
           }
       }
   }
   if !found {
       fmt.Println("IMPOSSIBLE")
   } else {
       fmt.Printf("%d:%d\n", bestX, bestY)
   }
}
