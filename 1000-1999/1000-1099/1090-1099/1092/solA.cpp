#include<bits/stdc++.h>
using namespace std;
int main()
{
int t;
cin>>t;
while(t--)
{
int n,k;
int val=96;
cin>>n;
cin>>k;
int p=n/k;
for(int i=0;i<k;i++)
{
val=val+1;
for(int j=0;j<p;j++)
{
cout<<char(val);
}
}
for(int i=0;i<n-k*p;i++)
cout<<char(val);
cout<<"\n";
}
}