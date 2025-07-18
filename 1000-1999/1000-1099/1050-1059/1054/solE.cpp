#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <algorithm>
#include <set>
#include <map>
#include <utility>
#include <queue>
#include <vector>
#include <string>
#include <cstring>

#define REP(a,n) for (int a = 0; a<(n); ++a)
#define FOR(a,b,c) for (int a = (b); a<=(c); ++a)
#define FORD(a,b,c) for (int a = (b); a>=(c); --a)
#define INIT(a,v) __typeof(v) a = (v)
#define FOREACH(a,v) for (INIT(a,(v).begin()); a!=(v).end(); ++a)

#define MP make_pair
#define PB push_back

using namespace std;

typedef pair<int, int> pii;
typedef vector<int> vi;
typedef vector<string> vs;
typedef long long LL;

///////////////////////////////

int xs, ys;

int jed[300], zer[300], chce1[300], chce0[300];

vector<pair<pii, pii> > R;
vector<pair<pii, pii> > SR;
vector<pair<pii, pii> > K;

int main() {
  char buf[200000];
  scanf("%d%d", &ys, &xs);
  REP(y, ys)
    REP(x, xs) {
      scanf("%s", buf);
      int L = strlen(buf);
      FORD(z, L-1, 0) {
        if (buf[z]=='1') {
          if (x==1) { ++jed[!y]; R.PB(MP(MP(y, x), MP(!y, x))); }
          else { ++jed[y]; R.PB(MP(MP(y, x), MP(y, 1))); }
        }
        else {
          if (x==0) { ++zer[!y]; R.PB(MP(MP(y, x), MP(!y, x))); }
          else { ++zer[y]; R.PB(MP(MP(y, x), MP(y, 0))); }
        }
      }
    }
  REP(y, ys)
    REP(x, xs) {
      scanf("%s", buf);
      int L = strlen(buf);
      FORD(z, L-1, 0) {
        if (buf[z]=='1') {
          if (x==1) { ++chce1[!y]; K.PB(MP(MP(!y, x), MP(y, x))); }
          else { ++chce1[y]; K.PB(MP(MP(y, 1), MP(y, x))); }
        }
        else {
          if (x==0) { ++chce0[!y]; K.PB(MP(MP(!y, x), MP(y, x))); }
          else { ++chce0[y]; K.PB(MP(MP(y, 0), MP(y, x))); }
        }
      }
    }
//REP(a, ys) { printf("ma0: %d, chce: %d\n", zer[a], chce0[a]); }
//REP(a, ys) { printf("ma1: %d, chce: %d\n", jed[a], chce1[a]); }
  for (int ym = 0, yd = 0; ; ) {
    while (ym<ys && jed[ym]>=chce1[ym]) ++ym;
    if (ym==ys) break;
    while (jed[yd]<=chce1[yd]) ++yd;
    --jed[yd];
    ++jed[ym];
    R.PB(MP(MP(yd, 1), MP(ym, 1)));
  }
  for (int ym = 0, yd = 0; ; ) {
    while (ym<ys && zer[ym]>=chce0[ym]) ++ym;
    if (ym==ys) break;
    while (zer[yd]<=chce0[yd]) ++yd;
    --zer[yd];
    ++zer[ym];
    R.PB(MP(MP(yd, 0), MP(ym, 0)));
  }
  printf("%d\n", R.size()+SR.size()+K.size());
  for (auto x : R)
    printf("%d %d %d %d\n", x.first.first+1, x.first.second+1, x.second.first+1, x.second.second+1);
  for (auto x : SR)
    printf("%d %d %d %d\n", x.first.first+1, x.first.second+1, x.second.first+1, x.second.second+1);
  for (auto x : K)
    printf("%d %d %d %d\n", x.first.first+1, x.first.second+1, x.second.first+1, x.second.second+1);
}