#include<iostream>
using namespace std;
int main()
{int t;
cin>>t;
while(t--)
{int n,mx=0;
cin>>n;
int a[n];
for(int i=0;i<n;i++)
{cin>>a[i];
if(!(i&1)&&(a[i]>mx))mx=a[i];}
cout<<mx<<endl;}}