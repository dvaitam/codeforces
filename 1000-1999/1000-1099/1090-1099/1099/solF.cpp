#include <bits/stdc++.h>
#include<bits/stdc++.h>
using namespace std;
long long fc[1000100],ft[1000100],V[100100],x[100100],t[100100];
vector<pair<int,long long> >node[100100];
void add(long long f[],long long v,long long x){
  for(int i=x;i<=1000000;i+=(i&-i))f[i]+=v;
}
long long sum(long long f[],int x){
  long long s=0;
  for(int i=x;i>0;i-=(i&-i))s+=f[i];
  return s;
}
long long bs(long long T){
  if(sum(ft,1000000)<=T)return sum(fc,1000000);
  int l=1,r=1000000,mid;
  long long s;
  while(l<r){
    mid=(l+r)/2;
    if(sum(ft,mid)>=T)r=mid;
    else l=mid+1;
  }
  s=sum(ft,l);
  if((s-T)%l==0)return sum(fc,l)-(s-T)/l;
  return sum(fc,l)-(s-T)/l-1;
}
long long F(int X,long long T){
  if(T<=0)return (long long)0;
  add(ft,x[X]*t[X],t[X]);
  add(fc,x[X],t[X]);
  V[X]=bs(T);
  long long M1=0,M2=0;
  for(int i=0;i<node[X].size();i++){
    long long z=F(node[X][i].first,T-node[X][i].second*2);
    if(M1<=z){
      M2=M1;
      M1=z;
    }
    else if(M2<z){
      M2=z;
    }
  }
  add(ft,x[X]*t[X]*-1,t[X]);
  add(fc,x[X]*-1,t[X]);
  if(X==1)return max(M1,V[X]);
  return max(M2,V[X]);
}
int main(){
  int n,i,p;
  long long T,l;
  scanf("%d%lld",&n,&T);
  for(i=1;i<=n;i++)scanf("%lld",&x[i]);
  for(i=1;i<=n;i++)scanf("%lld",&t[i]);
  for(i=2;i<=n;i++){
    scanf("%d%lld",&p,&l);
    node[p].push_back(make_pair(i,l));
  }
  printf("%lld",F(1,T));
  return 0;
}