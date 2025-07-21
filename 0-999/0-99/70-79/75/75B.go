package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read your name
   you, _ := reader.ReadString('\n')
   you = strings.TrimSpace(you)
   // read number of actions
   var n int
   fmt.Fscan(reader, &n)
   // consume end of line
   reader.ReadString('\n')
   // map for scores and set for names
   score := make(map[string]int)
   namesSet := make(map[string]struct{})
   for i := 0; i < n; i++ {
       line, _ := reader.ReadString('\n')
       line = strings.TrimSpace(line)
       if line == "" {
           i--
           continue
       }
       words := strings.Split(line, " ")
       if len(words) < 3 {
           continue
       }
       X := words[0]
       var Yraw string
       var points int
       action := words[1]
       switch action {
       case "posted":
           // X posted on Y's wall
           points = 15
           if len(words) > 3 {
               Yraw = words[3]
           }
       case "commented":
           // X commented on Y's post
           points = 10
           if len(words) > 3 {
               Yraw = words[3]
           }
       case "likes":
           // X likes Y's post
           points = 5
           if len(words) > 2 {
               Yraw = words[2]
           }
       default:
           continue
       }
       // Yraw is like "name's"
       Y := strings.TrimSuffix(Yraw, "'s")
       // record names
       if X != you {
           namesSet[X] = struct{}{}
       }
       if Y != you {
           namesSet[Y] = struct{}{}
       }
       // update scores only if involves you
       if X == you && Y != you {
           score[Y] += points
       } else if Y == you && X != you {
           score[X] += points
       }
   }
   // collect names
   names := make([]string, 0, len(namesSet))
   for name := range namesSet {
       names = append(names, name)
       // ensure default score entry
       if _, ok := score[name]; !ok {
           score[name] = 0
       }
   }
   // sort by descending score, then lexicographically
   sort.Slice(names, func(i, j int) bool {
       si, sj := score[names[i]], score[names[j]]
       if si != sj {
           return si > sj
       }
       return names[i] < names[j]
   })
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for _, name := range names {
       fmt.Fprintln(writer, name)
   }
}
