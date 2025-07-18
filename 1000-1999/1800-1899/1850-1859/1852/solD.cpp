/**
 *    author:  tourist
 *    created: 23.07.2023 11:30:59       
**/
#include <bits/stdc++.h>

using namespace std;

#ifdef LOCAL
#include "algo/debug.h"
#else
#define debug(...) 42
#endif

int main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
  int tt;
  cin >> tt;
  while (tt--) {
    int n, k;
    cin >> n >> k;
    string foo;
    cin >> foo;
    vector<int> s(n);
    for (int i = 0; i < n; i++) {
      s[i] = (int) (foo[i] - 'A');
    }
    int init = 0;
    for (int i = 0; i < n - 1; i++) {
      init += (s[i] != s[i + 1]);
    }
    const int inf = (int) 1e9;
    vector<array<array<int, 2>, 2>> dp_min(n);
    vector<array<array<int, 2>, 2>> dp_max(n);
    for (int i = 0; i < n; i++) {
      for (int c = 0; c < 2; c++) {
        for (int p = 0; p < 2; p++) {
          dp_min[i][c][p] = inf;
          dp_max[i][c][p] = -inf;
        }
      }
    }
    for (int c = 0; c < 2; c++) {
      int val = init + (c != s[0]);
      dp_min[0][c][val & 1] = val;
      dp_max[0][c][val & 1] = val;
    }
    for (int i = 1; i < n; i++) {
      for (int c = 0; c < 2; c++) {
        for (int t = 0; t < 2; t++) {
          int add = 0;
          add += (c != t);
          add += (c != s[i]);
          for (int p = 0; p < 2; p++) {
            dp_min[i][c][p] = min(dp_min[i][c][p], dp_min[i - 1][t][(p + add) & 1] + add);
            dp_max[i][c][p] = max(dp_max[i][c][p], dp_max[i - 1][t][(p + add) & 1] + add);
          }
        }
      }
    }
    vector<int> res(n, -1);
    for (int c = 0; c < 2; c++) {
      if (dp_min[n - 1][c][k & 1] <= k && k <= dp_max[n - 1][c][k & 1]) {
        res[n - 1] = c;
        break;
      }
    }
    if (res[n - 1] == -1) {
      cout << "NO" << '\n';
      continue;
    }
    cout << "YES" << '\n';
    for (int i = n - 1; i > 0; i--) {
      int c = res[i];
      int p = k & 1;
      for (int t = 0; t < 2; t++) {
        int add = 0;
        add += (c != t);
        add += (c != s[i]);
        if (dp_min[i - 1][t][(p + add) & 1] <= k - add && k - add <= dp_max[i - 1][t][(p + add) & 1]) {
          res[i - 1] = t;
          k -= add;
          break;
        }
      }
      assert(res[i - 1] != -1);
    }
    for (int i = 0; i < n; i++) {
      cout << (char) (res[i] + 'A');
    }
    cout << '\n';
  }
  return 0;
}