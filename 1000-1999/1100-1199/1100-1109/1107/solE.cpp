#include <bits/stdc++.h>
using namespace std;
typedef long long LL;

long long dp[101][101][101];


class Solution {
  int n, m;
  string s;
  long long B[100];
  vector<long long> A;
  long long go(int i, int j, int k) {
    if (dp[i][j][k] != -1) return dp[i][j][k];
    if (j - i == 0) return A[k];
    if (j - i == 1) return A[B[i] + k];
    long long ans = go(i, j - 1, 0) + A[B[j - 1] + k];
    for (int l = j - 2; l >= i; l -= 2) {
      ans = max(ans, go(i, l, B[j - 1] + k) + go(l, j - 1, 0));
    }
    return dp[i][j][k] = ans;
  }
public:
  void run() {
    memset(dp, -1, sizeof(dp));
    cin >> n;
    A.resize(n + 1);
    cin >> s;
    for (int i = 1; i <= n; ++i) {
      cin >> A[i];
      for (int j = 1; j < i; ++j) {
        A[i] = max(A[i], A[j] + A[i - j]);
      }
    }
    m = 0;
    for (int i = 0, j = 0; i < n; i = j) {
      while (j < n and s[i] == s[j]) ++j;
      B[m++] = j - i;
    }
    cout << go(0, m, 0) << '\n';
  }
};

int main() {
  ios::sync_with_stdio(0); cin.tie(0);
  Solution().run();
}