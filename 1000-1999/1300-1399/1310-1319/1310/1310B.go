package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Team struct {
   id  int
   fav bool
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   tot := 1 << n
   favs := make([]bool, tot+1)
   for i := 0; i < k; i++ {
       var a int
       fmt.Fscan(in, &a)
       if a >= 1 && a <= tot {
           favs[a] = true
       }
   }
   // initialize upper bracket
   ub := make([]Team, tot)
   for i := 1; i <= tot; i++ {
       ub[i-1] = Team{i, favs[i]}
   }
   // lower bracket list
   var lb []Team
   bad := 0
   // first upper bracket round
   // UB round 1
   var ubLosers []Team
   ub, ubLosers = playUBRound(ub, &bad)
   lb = ubLosers
   // rounds 2..n
   for r := 2; r <= n; r++ {
       // lower bracket first phase: LB-only matches
       lb1 := lbPlay(lb, &bad)
       // upper bracket round r
       ub, ubLosers = playUBRound(ub, &bad)
       // lower bracket second phase: LB1 survivors vs UB losers
       // sort both lists
       sort.Slice(lb1, func(i, j int) bool { return lb1[i].id < lb1[j].id })
       sort.Slice(ubLosers, func(i, j int) bool { return ubLosers[i].id < ubLosers[j].id })
       lb = make([]Team, 0, len(lb1))
       for i := 0; i < len(lb1); i++ {
           p, q := lb1[i], ubLosers[i]
           if !p.fav && !q.fav {
               bad++
           }
           var win Team
           if p.fav != q.fav {
               if p.fav {
                   win = p
               } else {
                   win = q
               }
           } else {
               if p.id < q.id {
                   win = p
               } else {
                   win = q
               }
           }
           lb = append(lb, win)
       }
   }
   // grand final
   if len(ub) == 1 && len(lb) == 1 {
       a, b := ub[0], lb[0]
       if !a.fav && !b.fav {
           bad++
       }
   }
   totalMatches := (1 << (n + 1)) - 2
   fmt.Println(totalMatches - bad)
}

// playUBRound plays one upper bracket round: returns winners and losers
func playUBRound(ub []Team, bad *int) (winners []Team, losers []Team) {
   sort.Slice(ub, func(i, j int) bool { return ub[i].id < ub[j].id })
   winners = make([]Team, 0, len(ub)/2)
   losers = make([]Team, 0, len(ub)/2)
   for i := 0; i+1 < len(ub); i += 2 {
       p, q := ub[i], ub[i+1]
       if !p.fav && !q.fav {
           (*bad)++
       }
       var win, lose Team
       if p.fav != q.fav {
           if p.fav {
               win, lose = p, q
           } else {
               win, lose = q, p
           }
       } else {
           // both fav or both non-fav: choose smaller id to win
           if p.id < q.id {
               win, lose = p, q
           } else {
               win, lose = q, p
           }
       }
       winners = append(winners, win)
       losers = append(losers, lose)
   }
   return
}

// lbPlay plays lower bracket intra-round: pairs lb, losers eliminated, winners returned
func lbPlay(lb []Team, bad *int) []Team {
   sort.Slice(lb, func(i, j int) bool { return lb[i].id < lb[j].id })
   nextLB := make([]Team, 0, len(lb)/2)
   for i := 0; i+1 < len(lb); i += 2 {
       p, q := lb[i], lb[i+1]
       if !p.fav && !q.fav {
           (*bad)++
       }
       var win Team
       if p.fav != q.fav {
           if p.fav {
               win = p
           } else {
               win = q
           }
       } else {
           // choose smaller id to win
           if p.id < q.id {
               win = p
           } else {
               win = q
           }
       }
       nextLB = append(nextLB, win)
   }
   return nextLB
}
