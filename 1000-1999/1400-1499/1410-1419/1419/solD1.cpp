#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
inline int read(){
	register int res=0;
	register char c=getchar(),f=1;
	while(c<48||c>57){if(c=='-')f=0;c=getchar();}
	while(c>=48&&c<=57)res=(res<<3)+(res<<1)+(c&15),c=getchar();
	return f?res:-res;
}
inline void write(int x){
	register char c[21],len=0;
	if(!x)return putchar('0'),void();
	if(x<0)x=-x,putchar('-');
	while(x)c[++len]=x%10,x/=10;
	while(len)putchar(c[len--]+48);
}
#define space(x) write(x),putchar(' ')
#define enter(x) write(x),putchar('\n')

const int N=1e5+10;
int n,a[N],m;
int main(){
	n=read();
	for(int i=1;i<=n;++i)a[i]=read();
	sort(a+1,a+n+1);
	enter(m=(n-1)>>1);
	for(int i=1;i<=m;++i)space(a[m+i]),space(a[i]);
	if(n&1)enter(a[n]);
	else space(a[n-1]),enter(a[n]);
}