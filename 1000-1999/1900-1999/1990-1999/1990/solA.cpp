#include<bits/stdc++.h>
using namespace std;
int main()
{
int t,n,x,pp,a[100];
cin>>t;
while(t--){
cin >>n;
memset (a,0,sizeof (a));
while(n--){
cin >>x;
a[x]++;
}
n=51;
pp=0;
while(--n){
if (a[n]&1)pp=1;
}
if(pp)puts("YES");
else puts("NO");
}
}