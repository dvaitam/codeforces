//　静かな海
#include <bits/stdc++.h>
using namespace std;

inline int rd() {
  int x=0,c,f=1;while(!isdigit(c=getchar()))if(c=='-')f=-1;
  for(;isdigit(c);c=getchar())x=x*10+c-'0';return x*f;
}
const int N = 1e5+3, M = 1e9+7;
typedef long long ll;

int n, k, f[N], sz[N];
int find(int x) {return x==f[x]?x:f[x]=find(f[x]);}
inline ll qpow(ll a, ll b, ll ans = 1) {
  for (a %= M; b; b>>=1, a = a * a % M)
    if (b & 1) ans = ans * a % M;
  return ans;
}

int main() {
  n = rd(), k = rd();
  for (int i=1; i<=n; i++) f[i] = i;
  for (int i=1; i<n; i++) {
    int u = rd(), v = rd(), w = rd();
    if (!w) f[find(u)] = find(v);
  }
  for (int i=1; i<=n; i++)
    sz[find(i)]++;
  ll ans = qpow(n, k);
  for (int i=1; i<=n; i++)
    ans = (ans - qpow(sz[i], k) + M) % M;
  cout << ans << endl;
}