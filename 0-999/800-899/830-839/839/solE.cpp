#include <cmath>
#include <cstdio>
#include <cstring>
#include <iostream>
#include <algorithm>
#include <iomanip>

typedef long long ll;
typedef long double ld;

using namespace std;

struct Solver {
private:
  static const int N = 40;
  int n, k, tot, a[N + 1][N + 1], tmp[N + 1];
  int rec[N + 1], ans[N + 1];
  double l[N + 1], ret;

  void input() {
    cin >> n >> k;
    for (int i = 1; i <= n; ++ i)
      for (int j = 1; j <= n; ++ j)
        cin >> a[i][j];
  }

  void process() {
    int T = 7000;
    for (int i = 1; i <= n; ++ i) tmp[i] = i;
    while (T --) {
      random_shuffle(tmp + 1, tmp + n + 1);
      int t = 0;
      for (int i = 1; i <= n; ++ i) {
        rec[++ t] = tmp[i];
        for (int j = 1; j < t; ++ j) {
          if (! a[rec[j]][rec[t]]) {
            -- t; break;
          }
        }
      }
      if (t > tot) {
        tot = t;
        for (int i = 1; i <= t; ++ i) ans[i] = rec[i];
      }
    }
    for (int i = 1; i <= tot; ++ i)
      l[ans[i]] = 1.0 * k / tot;
    ret = 0;
    for (int i = 1; i < n; ++ i)
      for (int j = i + 1; j <= n; ++ j)
          ret += l[i] * l[j];
    cout << setprecision(10) << fixed << ret << endl;
  }

public:
  void solve() {
    input(), process();
  }
} solver;

int main() {
#ifdef _DEBUG
  freopen("test.in", "r", stdin);
#endif
  ios::sync_with_stdio(false);
  solver.solve();
  return 0;
}