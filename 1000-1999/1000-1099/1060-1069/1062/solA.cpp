//Daniel Grzegorzewski
#include <bits/stdc++.h>
#pragma GCC optimize("O3")

#define MP make_pair
#define PB push_back
#define ST first
#define ND second

using namespace std;

typedef pair<int, int> PII;
typedef vector<int> VI;
typedef vector<PII> VII;
typedef long long LL;

void init_ios() {
     ios_base::sync_with_stdio(0);
     cin.tie(0);
}

int n, a[105];

int main() {
  init_ios();
  cin >> n;
  for (int i = 1; i <= n; ++i)
    cin >> a[i];
  if (n == 1) {
    cout<<"0\n";
    return 0;
  }
  int hm = 1;
  while (hm <= n && a[hm] == hm)
    ++hm;
  int res = max(0, hm-2);
  hm = n;
  int co = 1000;
  if (a[n] == 1000) {
    while (hm >= 1 && a[hm] == co) {
      --hm;
      --co;
    }
    res = max(res, n-hm-1);
  }
  for (int i = 1; i+1 < n; ++i)
    for (int j = i+2; j <= n; ++j)
      if (a[j]-a[i] == j-i)
        res = max(res, j-i-1);
  cout<<res<<"\n";
}