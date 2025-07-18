#include<cstdio>
#include<iostream>
#include<cmath>
#include<algorithm>
#include<string>
#include<cstring>
#include<cctype>
#include<queue>
#include<stack>
#include<map>
#include<set>
#include<iomanip>
#include<sstream>
#include<vector>
#include<cstdlib>
#include<ctime>
#include<list>
#include<deque>
#include<bitset>
#include<fstream>
#define ld double
#define ull unsigned long long
#define ll long long
#define pii pair<int,int >
#define iiii pair<int,pii >
#define mp make_pair
#define INF 1000000000
#define MOD 998244353
#define rep(i,x) for(int (i)=0;(i)<(x);(i)++)
inline int getint(){
	int x=0,p=1;char c=getchar();
	while (c<=32)c=getchar();
	if(c==45)p=-p,c=getchar();
	while (c>32)x=x*10+c-48,c=getchar();
	return x*p;
}
using namespace std;
//ruogu
const int N=2e5+10;
int n,m,A,b[N],res=1,inv2;
//
inline int add(int x,int y){
	x+=y;
	if(x>=MOD)x-=MOD;
	return x%MOD;
}
inline int mul(int x,int y){
	ll ans=1ll*x*y;
	return ans%MOD;
}
inline int modpow(int x,int y){
	int ans=1;
	while(y>0){
		if(y&1)ans=mul(ans,x);
		x=mul(x,x);
		y>>=1;
	}
	return ans;
}
int main(){
	n=getint();m=getint();A=getint();
	inv2=modpow(2,MOD-2);
	for(int i=1;i<=m;i++)b[i]=getint();
	for(int i=1;i<=m;i++){
		ll tmp=modpow(A,b[i]-b[i-1]);
		res=mul(res,mul(mul(tmp,add(tmp,1)),inv2));
	}
	int r=n-b[m];
	res=mul(res,modpow(A,r-b[m]));
	printf("%d\n",res);
	return 0;
}