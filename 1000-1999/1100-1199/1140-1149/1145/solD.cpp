//#pragma comment(linker, "/stack:200000000")
//#pragma GCC optimize("Ofast")
//#pragma GCC target("sse, sse2, sse3, ssse3, sse4, popcnt, abm, mmx, avx, tune = native")
#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
#define loop(i, n) for(int i = 0; i < n; i++)
#define loop1(i, n) for(int i = 1; i <= n; i++)
int main(){
  ios::sync_with_stdio(false);
  cin.tie(0);
  cout.tie(0);
  int n;
  cin >> n;
  int a[n], m = 1e9;
  loop(i, n){
    cin >> a[i];
    m = min(m, a[i]);
  }
  cout << 2 + (m ^ a[2]);
  return 0;
}