#include <iostream>
#include <string>
#include <cstdio>
#include <vector>
#include <algorithm>
#include <functional>
#include <map>
#include <set>
#include <fstream>

using namespace std;

typedef long long i64;
typedef int int_t;
typedef vector<int_t> vi;
typedef vector<vi> vvi;

typedef pair<int_t, int_t> pi;

#define tr(c, i) for(typeof((c).begin()) i = (c).begin(); i < (c).end(); i++)
#define pb push_back
#define sz(a) int((a).size())
#define all(c) (c).begin(), (c).end()
#define REP(s, e, i) for(i=(s); i < (e); ++i)

int main(int argc, char *argv[]) {
  ifstream ifs("input.txt");
  ofstream ofs("output.txt");

  int n, m, i;
  vi a, b;
  ifs >> n;
  a = vi(n);
  REP(0, n, i)
    ifs >> a[i];
  ifs >> m;
  b = vi(m);
  REP(0, m, i)
    ifs >> b[i];
  
  vi indexes(n + m);
  vi merged(n + m);
  vi flips;
  int ia, ib, im;
  ia = n-1;
  ib = m-1;
  im = n + m - 1;
  int c = 0;
  bool bf = true;
  while(ia >= 0 || ib >= 0) {
    while(ia >= 0 && a[ia] == c) {
      merged[im] = a[ia];
      indexes[im] = ia + 1;
      im--;
      ia--;
    }
    while(ib >= 0 && b[ib] == c) {
      merged[im] = b[ib];
      indexes[im] = ib + 1 + n;
      im--;
      ib--;
    }

    flips.pb(im + 1);
    
    c = (c == 0 ? 1 : 0);
  }

  REP(0, sz(indexes), i) {
    ofs << indexes[i] << " ";
  }
  ofs << endl;

  flips.erase(flips.end()-1);
  ofs << sz(flips) << endl;
  REP(0, sz(flips), i) {
    if(i != 0)
      ofs << " ";
    ofs << flips[sz(flips) - i - 1];
  }

  ofs.close();
  ifs.close();

  return 0;
}