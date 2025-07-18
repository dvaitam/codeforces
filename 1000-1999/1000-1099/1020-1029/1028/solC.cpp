#include<cstdio>
#include<cstring>
#include<algorithm>
#include<cmath>
#include<vector>
#include<map>
typedef long long ll;
using namespace std;
const int N=2e5+5;
const int inf=2e9;
int read(){
	int x=0,w=1;
	char ch=0;
	while (ch<'0' || ch>'9'){
		  if (ch=='-') w=-1;
		  ch=getchar();
	}
	while (ch<='9' && ch>='0'){
		  x=(x<<1)+(x<<3)+ch-'0';
		  ch=getchar();
	}
	return x*w;
}
int n,u[N],d[N],l[N],r[N];
int m[5],s[5];
int main(){
	n=read();
	m[1]=m[2]=s[1]=s[2]=inf;
	m[3]=m[4]=s[3]=s[4]=-inf;
	for (int i=1;i<=n;++i){
		d[i]=read();l[i]=read();
		if (m[3]<=d[i]){
			s[3]=m[3];
			m[3]=d[i];
		}else if (d[i]>=s[3]) s[3]=d[i];
		if (m[4]<=l[i]){
			s[4]=m[4];
			m[4]=l[i];
		}else if (l[i]>=s[4]) s[4]=l[i];
		u[i]=read();r[i]=read();	
		if (m[1]>=u[i]){
			s[1]=m[1];
			m[1]=u[i];	
		}else if (u[i]<=s[1]) s[1]=u[i];
		if (m[2]>=r[i]){
			s[2]=m[2];
			m[2]=r[i];	
		}else if (r[i]<=s[2]) s[2]=r[i];
	}
	for (int i=1;i<=n;++i){
		int U=m[1],D=m[3],L=m[4],R=m[2];
		if (u[i]==m[1]) U=s[1];
		if (r[i]==m[2]) R=s[2];
		if (d[i]==m[3]) D=s[3];
		if (l[i]==m[4]) L=s[4];
		if (L<=R && D<=U){
			printf("%d %d\n",D,L);
			return 0;
		}
	}
	return 0;
}