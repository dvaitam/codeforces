#include<bits/stdc++.h>
using namespace std;
int main()
{int t;
cin>>t;
while(t--)
{long long n;
cin>>n;
if(__builtin_popcountll(n)==1)cout<<1<<endl<<n<<endl;
else
{cout<<__builtin_popcountll(n)+1<<endl;for(int i=63;i>=0;i--)
{if((n&(1LL<<i)))cout<<(n^(1LL<<i))<<" ";}
cout<<n<<endl;}}}