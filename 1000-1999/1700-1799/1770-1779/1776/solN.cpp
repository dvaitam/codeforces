#include<iostream>
#include<ctype.h>
#include<iomanip>
#include<vector>
#include <cmath>
using namespace std;
inline int read(){
	char c=getchar();
	int x=0;
	bool f=0;
	for(;!isdigit(c);c=getchar()) f^=!(c^45);
	for(;isdigit(c);c=getchar()) x=(x<<1)+(x<<3)+(c^48);
	if(f)x=-x;
	return x;
}
const int maxn=400005;
const long long inf=0x3f3f3f3f;
int n,n1,n2;
char s[maxn];
int a[maxn],b[maxn];
long double f[maxn];
int F(int x,int y){
	return a[x]+b[y]-x-y+1;
}
int main(){
	n=read();
	cin>>(s+1);
	for(int i=1;i<n;i++){
		if(s[i]=='<') ++n2;
		else ++n1;
	}
	int x=n1,y=0;
	for(int i=1;i<n;i++){
		a[x]=max(a[x],y);
		b[y]=max(b[y],x);
		if(s[i]=='<') ++y;
		else --x;
	}
	a[x]=max(a[x],y);
	b[y]=max(b[y],x);
	long double delt=0;
	int B=1000;
	x=n1,y=0;
	f[0]=1;
	for(int i=1;i<=n;i++){
		long double mx=0;
		for(int j=0;j<=B;j++){
			int xx=x-j,yy=y-j;
			if(xx>=0 && yy>=0){
                f[j]/=F(xx,yy);
                mx=max(mx,f[j]);
			}
			else f[j]=0;
		}
		delt+=log2(mx);
		for(int j=0;j<=B;j++) f[j]/=mx;
		if(s[i]=='<'){
			for(int j=B;j>=1;--j) f[j]+=f[j-1];
			++y;
		}
		else{
			for(int j=1;j<=B;j++) f[j-1]+=f[j];
			--x;
		}
	}
	long double res=log2(f[0])+delt;
	for(int i=1;i<=n;i++) res+=log2(1.0l*i);
	cout<<setprecision(10)<<res<<'\n';
	return 0;
}