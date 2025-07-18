#include <bits/stdc++.h>
using namespace std;

namespace PROG
{
	inline int read(){
		int x=0,f=1;char cc=getchar();
		while(cc<'0' || cc>'9') {if(cc=='-') f=-1;cc=getchar();}
		while(cc>='0' && cc<='9') {x=x*10+cc-'0';cc=getchar();}
		return x*f;
	}
	
	#define maxn 300030
	#define ll long long
	int n,m,t[maxn*2],d[maxn*2],next[maxn*2],con[maxn];
	int s[maxn],g[maxn],v[maxn],num,flag;
	ll f[maxn],h[maxn],u[maxn],ans;
	
	inline void ins(int k,int x,int y,int z){
		t[k]=y;
		d[k]=z;
		next[k]=con[x];
		con[x]=k;
	}
	
	inline void dfs(int k,int fat){
		num++;
		s[num]=k;
		if(k==n){
			flag=1;
			return;
		}
		int p=con[k];
		while(p>0){
			f[num+1]=f[num]+d[p];
			if(t[p]!=fat) dfs(t[p],k);
			if(flag) return;
			p=next[p];
		}
		num--;
	}
	
    int main(void) {
    	ios::sync_with_stdio(false);

        n=read();m=read();
        for(int i=1;i<=n-1;i++){
        	int x=read();int y=read();int z=read();
        	ins(2*i-1,x,y,z);
        	ins(2*i,y,x,z);
		}
		num=0;f[1]=0;
		dfs(1,1);
		for(int i=1;i<=num;i++) h[i]=f[num]-f[i];
		for(int i=1;i<=num;i++) v[s[i]]=1;
		for(int i=1;i<=num;i++){
			int p=con[s[i]];
			while(p>0){
				if(!v[t[p]]){
					v[t[p]]=1;
					g[i]=d[p];
					break;
				}
				p=next[p];
			}
		}
		ans=2000000001;
		for(int i=1;i<=n;i++) if(!v[i]) ans=0;
		for(int i=1;i<=num;i++){
			u[i]=u[i-1];
			if(g[i]) u[i]=max(u[i],f[i]+g[i]);
		}
		for(int i=num;i>=2;i--) if(g[i] && u[i-1]) ans=min(ans,f[num]-u[i-1]-h[i]-g[i]);
		for(int i=1;i<=num-2;i++) ans=min(ans,f[num]-f[i]-h[i+2]);
		for(int i=1;i<=num-1;i++){
			if(g[i]) ans=min(ans,f[num]-f[i]-h[i+1]-g[i]);
			if(g[i+1]) ans=min(ans,f[num]-f[i]-h[i+1]-g[i+1]);
		}
		ans=max(ans,0ll);
		for(int i=1;i<=m;i++){
			int x=read();
			if(x>=ans) printf("%I64d\n",f[num]);
			else printf("%I64d\n",f[num]-ans+x);
		}

    	return 0;
    }
}

int main(void) {
    PROG::main();
    return 0;
}