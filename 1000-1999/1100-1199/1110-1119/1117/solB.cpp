#include <algorithm>
#include <iostream>
#include <cstring>
#include <cstdio>
#include <cmath>
using namespace std;
#define isnum(x) (x<58&&x>47)
int n,m,k;
long long maxx,secx;
long long read(){
	long long ans=0;
	char ch=getchar();
	while(!isnum(ch))ch=getchar();
	while(isnum(ch))ans=(ans<<1)+(ans<<3)+(ch^48),ch=getchar();
	return ans;
}
long long ans;
char s[44];
void print(long long x){
	int len=0;
	while(x){
		s[++len]=x%10^48;
		x/=10;
	}
	while(len)
		putchar(s[len--]);
	putchar(10);
}
int main(){
	scanf("%d%d%d",&n,&m,&k);
	for(int i=1;i<=n;++i){
		long long x=read();
		if(x>maxx)secx=maxx,maxx=x;
		else if(x>secx)secx=x;
	}
	int cnt=floor(m/(k+1));
	ans=maxx*(m-cnt)+secx*cnt;
	print(ans);
	return 0;
}