#include<stdio.h>
#include<iostream>
#include<string>
using namespace std;
int a,s,d[102],f,g,h,j,k,l,i,n,m;
string x,q,w,e;
main(){
cin>>n>>m>>k;

for(i=0;i<n;i++){
if(m>0) d[i]=k; else d[i]=1;
m=d[i]-m;

}
d[n-1]=d[n-1]-m;
if(d[n-1]>0 && d[n-1]<=k)
for(i=0;i<n;i++){
cout<<d[i]<<" ";
} 
else cout<<-1;
//system("pause");
}