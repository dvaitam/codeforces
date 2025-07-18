#include <bits/stdc++.h>
#define N 300005
#define rep(i,x,y) for(int i=(x);i<(y);++i)
int A[N],n,k,s;
using namespace std;
int main(){
 cin>>n>>k;
 rep(i,1,n+1){
  for(int j=i<<1;j<=n;j+=i)++A[j];
  s+=A[i]; 
  if(s>=k) {n=i; break;}
 }
 
 if(s<k) {puts("No"); return 0;} 
 s-=k;
 k=n;
 if(s)rep(i,2,n+1)if(s>=A[i]+n/i-1){s-=A[i]+n/i-1;A[i]=-1;--k;}
 printf("Yes\n%d\n",k);
 
 for(int i=1;i<=n;++i) if(~A[i])printf("%d ",i);
 puts (""); 
 return 0;
}