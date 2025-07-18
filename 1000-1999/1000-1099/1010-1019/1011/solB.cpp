#include <bits/stdc++.h>
using namespace std;
#define rep(i, n)    for(int i = 0; i < (n); ++i)
#define repA(i, a, n)  for(int i = a; i <= (n); ++i)
#define repD(i, a, n)  for(int i = a; i >= (n); --i)
#define trav(a, x) for(auto& a : x)
#define all(x) x.begin(), x.end()
#define sz(x) (int)(x).size()
#define fill(a)  memset(a, 0, sizeof (a))
#define fst first
#define snd second
#define mp make_pair
#define pb push_back
typedef long double ld;
typedef long long ll;
typedef pair<ll, int> pii;
typedef vector<int> vi;
bool sortbysec(const pair<int,int> &a,
              const pair<int,int> &b)
{
    return (a.second < b.second);
}
bool sortbyfst(const pair<int,int> &a,
              const pair<int,int> &b)
{   if(a.first == b.first){
        return (a.second < b.second);
    }
    return (a.first > b.first);
}
int main() {
  int n,m;
  cin >> n >> m;
  vi v(100,0);
  int t;
  rep(i,m){
     cin >> t;
     v[t-1]++;
   }
   if(n>m){
     cout << 0 << endl;
     return 0;
   }
   sort(v.rbegin(),v.rend());
   int d = 0;
   while(v[d]>0)d++;
   int mi = 1;
   int p = 0;
   rep(i,d)p+=(v[i]/mi);
   while(p>=n){
     mi++;
     int q = 0;
     rep(i,d)q+=(v[i]/mi);
     p = q;
   }

   cout << mi-1 << endl;

  return 0;
}