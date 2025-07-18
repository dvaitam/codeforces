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

const int N=3e5+10;
int n,q,a[N],b[N],cnt;
inline void add(int x,int y){a[++cnt]=x,b[cnt]=y;}
void solve(){
	n=read(),cnt=0;
	while(n>4){
		q=sqrt(n)+1;
		for(int i=q+1;i<n;++i)add(i,n);
		add(n,q),add(n,(n+q-1)/q);
		n=q;
	}
	if(n==3)add(3,2),add(3,2);
	else if(n==4)add(3,4),add(4,2),add(4,2);
	enter(cnt);
	for(int i=1;i<=cnt;++i)space(a[i]),enter(b[i]);
}

int main(){
	int T=read();
	while(T--)solve();
}