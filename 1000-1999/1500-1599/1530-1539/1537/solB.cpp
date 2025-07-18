#pragma GCC optimize ("trapv")
#include <bits/stdc++.h>
#include<algorithm>
#include <vector>
#include<cstring>
#include<cmath>
#include<cstdlib>
#include<string.h>
using namespace std;
#define pb push_back
#define all(v) v. begin(),v. end()
#define  rep(i,n,v) for(i=n;i<v;i++)
#define per(i,n,v) for(i=n;i>v;i--)
#define ff first 
#define ss second 
#define pp pair<ll,ll>
#define ll  long long
#define endl '\n'

void solve()
{
  ll  n, a,b=-1,m=0, c=0,k=0, i, j, l=1e9+5;
  string s;
  cin>>n>>a>>b>>m;
  c=abs(b-1)+abs(m-1);
 l=abs(m-a)+abs(n-b);
  j= abs(b-n)+abs(m-1);
  k=abs(b-1)+abs(m-a);
 // c=abs(b-1)+abs(m-a)+abs(n-b)+abs(a-m);
 i=max(max(c, l), max(j, k));
  if(i==c || i==c)
  {
    cout<<1<<" 1 "<<n<<" "<<a<<endl;
  }
  else 
  cout<<n<<" "<<1<<" "<<1<<" "<<a<<endl;
} 
int main()
{
 ios_base::sync_with_stdio(false);
  cin. tie(0);cout. tie(0);
    ll t=1;

  cin>>t;
    while(t--)
    {
      solve();
    }
    return 0;
}