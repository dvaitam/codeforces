#include<iostream>
#include<cstdio>
#include<algorithm>
#include<cstring>
#include<cmath>
#include<set>
#include<map>
#include<queue>
using namespace std;
typedef long long LL;
#define sqr(x) ((x)*(x))
#define mp make_pair
#define F first
#define S second
inline int read(){
	int x = 0; char ch = getchar(); bool positive = 1;
	for (; !isdigit(ch); ch = getchar())	if (ch == '-')	positive = 0;
	for (; isdigit(ch); ch = getchar())	x = x * 10 + ch - '0';
	return positive ? x : -x;
}
inline void write(int a){
    if(a>=10)write(a/10);
    putchar('0'+a%10);
}
inline void writeln(int a){
    if(a<0){
    	a=-a; putchar('-');
	}
	write(a); puts("");
}
#define DEBUG 0
int n,k;
inline void solve(int l,int r,int k,int L,int R){
	//if(k>0)cout<<l<<" "<<r<<" "<<k<<" "<<L<<" "<<R<<endl;
	if(k==1){
		for(int i=L;i<=R;i++){
			write(i); putchar(' ');
		}
		return;
	}
	int mid=(l+r-1)>>1,Mid=(L+R+2)>>1;
	if(2*(mid-l+1)-1>=k-2){
		solve(l,mid,k-2,Mid,R); solve(mid+1,r,1,L,Mid-1);
	}else{
		solve(l,mid,2*(mid-l+1)-1,Mid,R); solve(mid+1,r,k-2*(mid-l+1),L,Mid-1);
	}
}
int main(){
	if(DEBUG){
		freopen("1.in","r",stdin);
	}
	n=read(); k=read();
	if(k%2==0||n*2<=k){
		puts("-1"); return 0;
	}
	solve(1,n,k,1,n);
}