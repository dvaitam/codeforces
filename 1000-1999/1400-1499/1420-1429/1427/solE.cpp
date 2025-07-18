#include <bits/stdc++.h>

using namespace std;

int main() {

 long long x;

 cin>>x;

 vector<array<long long,3> > ans;

 while (x>1) {

  int i = 0;

  while ((1ll<<i+1)<=x) {

   ans.push_back({x<<i,x<<i,0});

   ++i;

  }

  long long a = x<<i, b = a^x;

  ans.push_back({a,x,1});

  for (int j=0; j<i; ++j) ans.push_back({x<<j,a,0}),ans.push_back({x<<j,b,0}),a+=x<<j,b+=x<<j;

  ans.push_back({a,b,1});

  x = a^b;

 }

 printf("%d\n", (int)ans.size());

 for (auto &t:ans) printf("%lld %c %lld\n", t[0], "+^"[t[2]], t[1]);

}