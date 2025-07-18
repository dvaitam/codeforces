#include<iostream>
#include<cstdio>
#include<cmath>
#include<cstring>
#include<queue>
#include<stack>
#include<vector>
#include<algorithm>
#include<numeric>
#include<string>
#include<sstream>
#include<set>
#include<map>
#include<complex>
#include<cassert>
using namespace std;
#define rep(i,s,e) for(int i=(s),___e=(e);i<___e;++i)
#define REP(i,n) rep(i,0,n)
#define ITER(c) __typeof((c).begin())
#define FOR(i,c) for(ITER(c) i=(c).begin(),___e=(c).end();i!=___e;++i)
#define ALL(c) (c).begin(), (c).end()
#define MP make_pair
typedef unsigned int ui;
typedef long double ld;
typedef long long ll;
const double PI = atan(1.0) * 4.0;
int in_c() { int c; while ((c = getchar()) <= ' ') { if (!~c) throw ~0; } return c; }
int in() {
  int x = 0, c;
  while ((ui)((c = getchar()) - '0') >= 10) { if (c == '-') return -in(); if (!~c) return ~0; }
  do { x = 10 * x + (c - '0'); } while ((ui)((c = getchar()) - '0') < 10);
  return x;
}

struct S {
  int p, id;
  S(int p, int id) : p(p), id(id) {}
  bool operator < (const S& s) const {
    return p > s.p;
  }
};

int main()
{
  int N = in(), V = in();
  vector<S> ks, cs;
  REP(i, N) {
    int t = in(), p = in();
    if(t == 1)
      ks.push_back(S(p, i + 1));
    else
      cs.push_back(S(p, i + 1));
  }
  sort(ALL(ks)); // size 1
  sort(ALL(cs)); // size 2
  vector<int> sum_c(cs.size() + 1), sum_k(ks.size() + 1);
  sum_c[0] = sum_k[0] = 0;
  REP(i, cs.size()) sum_c[i+1] = sum_c[i] + cs[i].p;
  REP(i, ks.size()) sum_k[i+1] = sum_k[i] + ks[i].p;
//  cout << cs.size() << ' ' << ks.size() << endl;
//  FOR(it, ks) cout << *it << endl;
  int ans = 0, used_c = 0, used_k = 0;
  for(int use_k = 0; use_k <= min(V, (int)ks.size()); ++use_k) {
    int use_c = min((V - use_k) / 2, (int)cs.size());
    int ps = 0;
//     REP(i, use_c) ps += cs[i].p;
//     REP(i, use_k) ps += ks[i].p;
    ps = sum_c[use_c] + sum_k[use_k];
//    cout << use_c << ' ' << use_k << ' ' << ps << endl;
    if(ans < ps) {
      ans = ps;
      used_c = use_c;
      used_k = use_k;
    }
  }
  cout << ans << endl;
//  cout << used << endl;
  
  vector<int> ids;
  REP(i, used_c) ids.push_back(cs[i].id);
  REP(i, used_k) ids.push_back(ks[i].id);
  if(!ids.empty()) {
    sort(ALL(ids));
    REP(i, ids.size() - 1) cout << ids[i] << ' ';
    cout << ids.back() << endl;
  } else {
    cout << endl;
  }
  return 0;
}