//g++ -std=c++14 -g -O2 -o ./a ./A.cpp
#include <bits/stdc++.h>
using namespace std;
#define ff first
#define ss second
#define nl '\n'
typedef long long ll;
//////////////////////////////////////////////////////////////////////

const int N = 1000100;
char s[N],t[N];
int n,m;
const int C = 26;
int tc[C],sc[C+1];

bool poss(ll suit){
  ll tot = 0;
  for(int i=0;i<C;i++){
    ll temp = suit * tc[i] - sc[i];
    if(temp < 0)continue;
    tot += temp;
  }
  //cerr << suit << ": " << tot  << endl;
  return tot <= sc[C];
}

int main(){
  ios::sync_with_stdio(0);cin.tie(0);cout.tie(0);

  cin>>1+s>>1+t;
  n = strlen(1+s);
  m = strlen(1+t);
  fill(tc,tc+C,0);
  
  for(int i=1;i<=m;i++)
    tc[t[i]-'a']++;

  fill(sc,sc+C+1,0);
  for(int i=1;i<=n;i++){
    if(s[i]=='?')sc[C]++;
    else sc[s[i]-'a']++;
  }
  
  ll lo = 0 , hi = N;
  while(lo < hi){
    ll mid = (lo+hi+1)/2;
    if(poss(mid))lo = mid;
    else hi=mid-1;
  }
  
  for(int i=0;i<C;i++)
    sc[i] = max(tc[i] * lo - sc[i],0ll);
  
  //cerr << lo << endl;
  
  for(int i=1,j=0;i<=n;i++)
    if(s[i]=='?'){
      while(j<C and sc[j]==0)j++;
      if(j==C){
	s[i]='t';
	continue;
      }
      s[i] = 'a'+j;
      sc[j]--;
    }

  cout << 1+s << nl;
  
  return 0;
}