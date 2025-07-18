#include <bits/stdc++.h>
#define ll long long
#define ld long double
#define pii pair<int, int>
#define pll pair<ll, ll>
#define str string
#define fi first
#define se second
#define pb push_back
#define SET(a, b) memset(a, b, sizeof(a))
#define eps 1e-6
#define pi atan(1) * 4
#define mod 1000000007
#define inf 1000000000
#define llinf 1000000000000000000
#define FOR(i, a, b) for (int i = (a); i <= (b); i++)
#define FORD(i, a, b) for (int i = (a); i >= (b); i--)
#define FORl(i, a, b) for (ll i = (a); i <= (b); i++)
#define FORDl(i, a, b) for (ll i = (a); i >= (b); i--)
using namespace std;
int n;
bool ini[2005][2005], fin[2005][2005], mag[2005];
bool tmp[2005][2005];
int c1;
void column (int i) {
  FOR(j, 1, n) {
    ini[j][i] ^= mag[j];
  }
}
void row (int i) {
  FOR(j, 1, n) {
    ini[i][j] ^= mag[j];
  }
}
vector<pii> oper;
void output () {
  printf("%d\n", (int)oper.size());
  for (auto u : oper) {
    --u.se;
    if (u.fi) printf("col %d\n", u.se);
    else printf("row %d\n", u.se);
  }
}
void check () {
  FOR(i, 1, n) {
    if (ini[i][c1] != fin[i][c1]) {
      row(i); oper.pb({0, i});
    }
  }
  FOR(j, 1, n) {
    if (j == c1) continue;
    int cnt = 0, cnt2 = 0;
    FOR(i, 1, n) {
      cnt += ((ini[i][j] ^ mag[i]) == fin[i][j]);
      cnt2 += (ini[i][j] == fin[i][j]);
    }
    if (cnt == n) {
      oper.pb({1, j});
    } else if (cnt2 != n) {
      return;
    }
  }
  output(); exit(0);
}
int main () {
  ios::sync_with_stdio(false); cin.tie(0);
  cin >> n;
  FOR(i, 1, n) {
    string S;
    cin >> S;
    FOR(j, 1, n) {
      tmp[i][j] = ini[i][j] = (S[j - 1] == '1');
    }
  }
  FOR(i, 1, n) {
    string S;
    cin >> S;
    FOR(j, 1, n) {
      fin[i][j] = (S[j - 1] == '1');
    }
  }
  string S;
  cin >> S;
  FOR(j, 1, n) {
    mag[j] = (S[j - 1] == '1');
    if (mag[j]) c1 = j;
  }
  check();
  oper.clear();
  memcpy(ini, tmp, sizeof(tmp));
  column(c1); oper.pb({1, c1});
  check();
  printf("-1\n");
  return 0;
}