#include<iostream>
using namespace std;
int main(){
int t;cin>>t;
while(t--){
int n;cin>>n;
int b[n+1],a[n+1],d=1;
for(int i=1;i<=n-1;i++)cin>>b[i];
b[0]=b[n]=0;
for(int i=1;i<=n;i++){
a[i]=(b[i]|b[i-1]);
}
for(int i=2;i<=n;i++){if((a[i]&a[i-1])!=b[i-1])d=0;
if(!d)break;}
if(!d)cout<<-1<<endl;
else {for(int i=1;i<=n;i++)cout<<a[i]<<" ";
cout<<endl;}}}