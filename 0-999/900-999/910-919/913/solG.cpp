#include <bits/stdc++.h>
using namespace std;
typedef long long               ll;
ll mod = 1;

ll mult_mod(ll a, ll b){
 ll res = a*(b>>16)%mod;
 res = (res<<16)%mod;
 res = (res + a*(b&((1<<16)-1))) % mod;
 return res;
}

ll almost_root_pow[20];
ll almost_root[20];

int main(){
 for(int i=0; i<17; ++i){
  mod *=5;
 }
 almost_root[0] = 2;
 almost_root[1] = 16;
 almost_root_pow[0] = 1;
 almost_root_pow[1] = 4;
 for(int i=2; i<=17; ++i){
  //to power 5;
  almost_root[i] = mult_mod(almost_root[i-1], almost_root[i-1]);
  almost_root[i] = mult_mod(almost_root[i], almost_root[i]);
  almost_root[i] = mult_mod(almost_root[i], almost_root[i-1]);
  almost_root_pow[i] = almost_root_pow[i-1]*5;
 }
 int n;
 cin>>n;
 for(int i=0; i<n; ++i){
  ll a;
  cin>>a;
  a *= 1000*1000;
  a += 1LL<<17;
  a &= ~((1LL<<17)-1);
  if(!(a%5)){
   a += 1LL<<17;
  }
  ll p = 60;
  ll r = (1LL<<60)%mod;
  ll curmod =5;
  for(int modp = 1; modp <= 17; ++modp){
   assert((almost_root[modp-1]-1)%curmod);
   while((r-a)%curmod){
    r = mult_mod(r, almost_root[modp-1]);
    p += almost_root_pow[modp-1];
   }
   curmod*=5;
  }
  assert(!((r-a)%mod));
  cout<<p<<'\n';
 }
 return 0;
}