#include<bits/stdc++.h>
using namespace std;
#define ll long long
inline ll read() {
    ll x=0,f=1; char ch=getchar();
    for(;ch<'0'||ch>'9';ch=getchar())
        if(ch=='-')f=-f;
    for(;ch>='0'&&ch<='9';ch=getchar())
        x=x*10+ch-'0';
    return x*f;
}
inline void chkmin( int &a,int b ) { if(a>b) a=b; }
inline void chkmax( int &a,int b ) { if(a<b) a=b; }
#define _ read()
const int N=2e5+5;
struct Edge { int to,nxt; } e[N<<1];  
int head[N],cnt,n,m,u[N],v[N],col[N],flag;
inline void insert( int u,int v ) { e[++cnt]=(Edge){v,head[u]}; head[u]=cnt; }
inline void ins( int u,int v ) {
	insert(u,v); insert(v,u);
}
inline void dfs( int u ){
    if(flag) return;
    for(int i=head[u];i;i=e[i].nxt){
        if(!col[e[i].to]){
            col[e[i].to]=3-col[u];
            dfs(e[i].to);
        }else if(col[e[i].to]==col[u]){
            flag=1;return;
        }
    }
}
int main() {
	n=_; m=_;
	for( int i=1;i<=m;i++ ) {
		u[i]=_; v[i]=_;
		ins(u[i],v[i]);
	}
	col[1]=1; dfs(1);
	if(flag) puts("NO");
	else {
		puts("Yes");
		for( int i=1;i<=m;i++ ) 
			putchar(col[u[i]]-1+'0');
	}
}