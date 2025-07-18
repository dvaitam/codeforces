#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef double db;
template<class T>inline void MAX(T &x,T y){if(y>x)x=y;}
template<class T>inline void MIN(T &x,T y){if(y<x)x=y;}
template<class T>inline void rd(T &x){
	x=0;char o,f=1;
	while(o=getchar(),o<48)if(o==45)f=-f;
	do x=(x<<3)+(x<<1)+(o^48);
	while(o=getchar(),o>47);
	x*=f;
}
const int M=1e7+5;
const int K=5e5+5;
int n,A[K],B[K];
int pre[M],prime[M],ptot;
void init(){
	for(int i=2;i<M;i++){
		if(pre[i]==0)prime[++ptot]=pre[i]=i;
		for(int j=1;j<=ptot;j++){
			int t=i*prime[j];
			if(t>=M)break;
			pre[t]=prime[j];
			if(i%prime[j]==0)break;
		}
	}
}
void solve(int x,int &ans1,int &ans2){
	if(x%2==0){
		while(x%2==0)x/=2;
		if(x==1)ans1=ans2=-1;
		else ans1=2,ans2=x;
	}
	else{
		int a=-1,b=-1;
		while(x>1){
			int t=pre[x];
			if(a==-1)a=t;
			else if(b==-1)b=t;
			while(x%t==0)x/=t;
		}
		if(a==-1||b==-1)ans1=ans2=-1;
		else ans1=a,ans2=b;
	}
}
int gcd(int a,int b){
	if(b==0)return a;
	return gcd(b,a%b);
}
int main(){
#ifndef ONLINE_JUDGE
	freopen("jiedai.in","r",stdin);
//	freopen("jiedai.out","w",stdout);
#endif
	init();
	rd(n);
	for(int i=1;i<=n;i++){
		int x;
		rd(x);
		solve(x,A[i],B[i]);
	}
	for(int i=1;i<=n;i++)printf("%d%c",A[i],i==n?'\n':' ');
	for(int i=1;i<=n;i++)printf("%d%c",B[i],i==n?'\n':' ');
	return (0-0);
}