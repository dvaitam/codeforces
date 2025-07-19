#include <bits/stdc++.h>
using namespace std;
int main()
{
int T;
cin>>T;
while(T--)
{
int n,k,ans=0;
cin>>n>>k;
if(k>=n)
k-=n,ans++;
for(int i=n-1;i>=1;i--)
{
if(k>=i)
k-=i,ans++;
if(k>=i)
k-=i,ans++;
}
cout<<ans<<endl;
}
}