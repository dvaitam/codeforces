#pragma GCC optimize("O3")

#include <bits/stdc++.h>

using namespace std;

#define all(arr) (arr).begin(), (arr).end()

#define ll long long

#define ld long double

#define pb push_back

#define sz(x) int((x).size())

#define fi first

#define se second

#define endl '\n'



void solve(){

  int n;

  cin >> n;

  string a, b;

  string x, y, z;

  cin >> x >> y >> z;

  char need;

  int c1x = count(all(x), '1'), c0x = 2 * n - c1x,

  c1y = count(all(y), '1'), c0y = 2 * n - c1y,

  c1z = count(all(z), '1'), c0z = 2 * n - c1z;

  if (c1x >= n){

    if (c1y >= n) a = x, b = y, need = '1';

    if (c1z >= n) a = x, b = z, need = '1';

  }

  if (c0x >= n){

    if (c0y >= n) a = x, b = y, need = '0';

    if (c0z >= n) a = x, b = z, need = '0';

  }

  if (c1y >= n){

    if (c1z >= n) a = y, b = z, need = '1';

  }

  if (c0y >= n){

    if (c0z >= n) a = y, b = z, need = '0';

  }



  string ans;

  int j = 0;

  for (int i = 0; i < 2 * n; i++){

    if (a[i] == need){

      while (j < 2 * n && b[j] != need) ans.pb(b[j]), j++;

      j++;

    }

    ans.pb(a[i]);

  }

  for (int i = j; i < 2 * n; i++) ans.pb(b[i]);

  cout << ans << endl;

}



int main(){

  ios_base::sync_with_stdio(0);

  cin.tie(0);

  // cout.precision(20);

  int t;

  cin >> t;

  while (t--){

    solve();

  }

}