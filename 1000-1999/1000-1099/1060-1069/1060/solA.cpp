#include <bits/stdc++.h>
using namespace std;

int main ()
{
  int n, ei = 0, z = 0, ans;
  string s;
  cin >> n;
  getline (cin,s);
  getline(cin,s);
  for (int i = 0; i < n; ++i){
    if (s[i] == '0')++z;
    else if (s[i] == '8')++ei;
  }
  if (n >= 11){
    n = n / 11;
    ans = min(n, ei);
    cout << ans << endl;
  }
  else cout << 0 << endl;

}