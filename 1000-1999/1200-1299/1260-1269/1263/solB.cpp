#include <bits/stdc++.h>
using namespace std;
void solve(){
  int n;
  cin >> n;
  string s;
  map<string, int>mp;
  vector<string>v;
  for (int i = 0; i < n; i++){
    cin >> s;
    v.push_back(s);
    mp[s]++;
  }
  int cnt = 0;
  for (int i = 0; i < n; i++)
  {
    if (mp[v[i]] > 1){
      cnt++; 
      string tmp = v[i];
      bool ok = 0;
      for (int j = 0; j <= 9; j++){
        tmp[3] = (char) (j + 48);
        if (mp[tmp] == 0){ ok = 1; break;}
      }
      if (ok == 0){
        tmp = v[i];
        for (int j = 0; j <= 9; j++){
          tmp[2] = (char) (j+ 48);
          if (mp[tmp] == 0) { ok = 1; break; }
        }
      }
      mp[v[i]]--;
      v[i] = tmp;
      mp[v[i]]++;
    }
  }
  cout << cnt << '\n';
  for (auto e: v)
    cout << e << '\n';
}
int main()
{
  ios::sync_with_stdio(false);
  cin.tie(0);
  int t;
  cin >> t;
  while (t--){
    solve();
  }
  return 0;
}