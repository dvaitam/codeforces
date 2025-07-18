#include <iostream>
#include <cmath>
#include <algorithm>
#include <string.h>
#include <vector>

using namespace std;
vector<int> F[60];
vector<int> IF[60];
int p[] = {2, 3, 5, 7, 11, 13, 17, 19, 23, 29};
int sol = 0;
int am[100];
int v[100];
int n, in[101];
int best_dif = 30 * 100;
vector<int> num;
int dp[101][20];

int solve(int i, int j) {
  if (i == 0) return 0;
  if (dp[i][j]) return dp[i][j] - 1;
  int ret = solve(i - 1, j) + in[i - 1] - 1;
  if (j > 0) ret = min(ret, solve(i, j - 1));
  if (j > 0) ret = min(ret, solve(i - 1, j - 1) + abs(in[i - 1] - num[j - 1]));
  dp[i][j] = ret + 1;
  return ret;
}
vector<int> best_sol;
void solve() {
  num.resize(0);
  for (int i = 1; i <= v[0]; ++i) num.push_back(v[i]);
  num.push_back(31);
  num.push_back(37);
  num.push_back(41);
  num.push_back(43);
  num.push_back(47);
  num.push_back(53);
  num.push_back(59);
  sort(num.begin(), num.end());
  // 18 X 100
  memset(dp, 0, sizeof(dp));
  int now = solve(n, num.size());
  if (now < best_dif) {
    best_dif = now;
    int X = now + 1;
    int i = n, j = num.size();
    vector<int> r;
    while (i > 0) {
      if (dp[i][j] - 1 == solve(i - 1, j) + in[i - 1] - 1) {
        r.push_back(1);
        i--;
      } else if (dp[i][j] - 1== solve(i, j - 1)) {
        j--;
      } else {
        r.push_back(num[j-1]);
        i--, j--;
      }
    }
    reverse(r.begin(), r.end());
    best_sol = r;
  }
}
void back(int k) {
  if (k == 10) {
    ++sol;
    solve();
    return;
  }
  int at_least_one = 0;
  for (int i = 0; i < IF[p[k]].size(); ++i) {
    int now = IF[p[k]][i];
    int ok = 1;
    for (int j = 0; j < F[now].size(); ++j) {
      if (am[F[now][j]]) {ok = 0; break; }
    }
    if (ok) {
      at_least_one = 1;
      for (int j = 0; j < F[now].size(); ++j) {
        am[F[now][j]] = 1;
      }
      v[++v[0]] = now;
      back(k + 1);
      v[0]--;
      for (int j = 0; j < F[now].size(); ++j) {
        am[F[now][j]] = 0;
      }
    }
  }
  if (!at_least_one) back(k + 1);
}
int in2[101];
int main() {
  cin >> n;
  for (int i = 0; i < n; ++i) { cin >> in[i]; in2[i] = in[i]; }
  sort(in, in + n);
  for (int i = 2; i < 60; ++i) {
    int X = i;
    for (int f = 2; f <= X; ++f) if (X % f == 0) {
      F[i].push_back(f);
      IF[f].push_back(i);
      while (X % f == 0) X /= f;
    }
  }
  back(0);
  for (int i = 0; i < n; ++i) {
    if (i != 0) cout << ' ';
    for (int j = 0; j < n; ++j) if (in2[i] == in[j]) {
      in[j] = -1;
      cout << best_sol[j];
      break;
    }
  }
  cout << '\n';
  return 0;
}