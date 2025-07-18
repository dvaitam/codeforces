#include <bits/stdc++.h>

using namespace std;

const int N = 2020;

int n, h;
long double p, p2;
int x[N];
long double dp[N][N][2][2];
long double res;
int L[N], R[N];

int get(int l, int r, int x, int y) {
  l = max(l, x);
  r = min(r, y);
  return r - l;
}

int main() {
  ios_base::sync_with_stdio(0);
  cin >> n >> h >> p;
  for (int i = 1; i <= n; i++) {
    cin >> x[i];
  }
  p2 = 1.0 - p;
  x[0] = -2e9;
  x[n + 1] = 2e9;
  sort(x + 1, x + n + 1);
  L[1] = 1;
  for (int i = 2; i <= n; i++) {
    if (x[i] - x[i - 1] < h) {
      L[i] = L[i - 1];
    } else {
      L[i] = i;
    }
  }
  R[n] = n;
  for (int i = n - 1; i; i--) {
    if (x[i + 1] - x[i] < h) {
      R[i] = R[i + 1];
    } else {
      R[i] = i;
    }
  }
  dp[1][n][0][1] = 1;
  for (int i = 1; i <= n; i++) {
    for (int j = n; j >= i; j--) {
      for (int a = 0; a < 2; a++) {
        for (int b = 0; b < 2; b++) {
          long double pro = dp[i][j][a][b];
          int ll, rr;
          if (a) {
            ll = x[i - 1] + h;
          } else {
            ll = x[i - 1];
          }
          if (b) {
            rr = x[j + 1];
          } else {
            rr = x[j + 1] - h;
          }
          //left
          res += pro * 0.5 * p * get(x[i] - h, x[i], ll, rr);
          dp[i + 1][j][0][b] += pro * 0.5 * p;
          res += pro * 0.5 * p2 * get(x[i], x[R[i]] + h, ll, rr);
          dp[R[i] + 1][j][1][b] += pro * 0.5 * p2;
          //right
          res += pro * 0.5 * p * get(x[L[j]] - h, x[j], ll, rr);
          dp[i][L[j] - 1][a][0] += pro * 0.5 * p;
          res += pro * 0.5 * p2 * get(x[j], x[j] + h, ll, rr);
          dp[i][j - 1][a][1] += pro * 0.5 * p2;
        }
      }
    }
  }
  printf("%0.17f", (double)res);
  return 0;
}