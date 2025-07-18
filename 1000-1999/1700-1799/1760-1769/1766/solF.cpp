#include<bits/stdc++.h>

using namespace std;

typedef long long ll;

typedef pair<int,int> pii;

typedef pair<ll,ll> pll;

#define all(x) x.begin(),x.end()



const int V=105,E=1005,inf=0x3f3f3f3f;

int S,T,e=1,fir[V],to[E],nxt[E],w[E],c[E];

inline void adde(int x,int y,int z,int t){

	//cout<<x<<' '<<y<<' '<<z<<' '<<t<<endl;

	to[++e]=y,nxt[e]=fir[x],fir[x]=e,w[e]=z,c[e]=t;

	to[++e]=x,nxt[e]=fir[y],fir[y]=e,w[e]=0,c[e]=-t;

}

ll dis[V];

int q[E];

bool vis[V];

bool spfa(){

	memset(dis,63,sizeof(dis));

	memset(vis,0,sizeof(vis));

	int l=1,r=0;

	q[++r]=T,dis[T]=0;

	while(l<=r){

		int u=q[l++];vis[u]=0;

		for(int i=fir[u],v=to[i];i;v=to[i=nxt[i]]){

			if(!w[i^1]||dis[v]<=dis[u]+c[i^1])continue;

			dis[v]=dis[u]+c[i^1];

			if(!vis[v])vis[v]=1,q[++r]=v;

		}

	}

	return dis[S]<0;

}

int cur[V];

int dfs(int u,int flow){

	if(u==T||!flow)return flow;

	vis[u]=1;

	int nowf=flow;

	for(int& i=cur[u];i;i=nxt[i]){

		int v=to[i];

		if(dis[v]+c[i]!=dis[u]||vis[v])continue;

		int f=dfs(v,min(w[i],nowf));

		w[i]-=f,w[i^1]+=f;

		if(!(nowf-=f))return flow;

	}

	return flow-nowf;

}

ll MCMF(){

	int flow=0;

	ll res=0;

	while(spfa()){

		memcpy(cur,fir,sizeof(cur));

		memset(vis,0,sizeof(vis));

		int f=dfs(S,inf);

		flow+=f,res+=dis[S]*f;

	}return res;

}



int n,m,x,y,z,t,f[V],ans[E],sum[V];

bool odd[E];

ll anss=0;

int main(){

	ios::sync_with_stdio(0);cin.tie(0);

	cin>>n>>m,S=1,T=n;

	for(int i=1;i<=m;++i){

		cin>>x>>y>>z>>t;

		adde(x,y,z>>1,t<<1);

		if(z&1)--f[x],++f[y],odd[i]=1,anss+=t;

	}

	for(int i=2;i<=n-1;++i){

		if(f[i]&1){

			cout<<"Impossible\n";

			return 0;

		}

		if(f[i]>0)adde(S,i,f[i]>>1,-inf);

		if(f[i]<0)adde(i,T,-f[i]>>1,-inf);

	}

	anss+=MCMF();

	for(int i=1;i<=m;++i){

		ans[i]=w[i*2+1]<<1|odd[i];

		sum[to[i*2+1]]-=ans[i];

		sum[to[i*2]]+=ans[i];

	}

	for(int i=2;i<=n-1;++i)if(sum[i]){

		cout<<"Impossible\n";

		return 0;

	}

	cout<<"Possible"<<endl;

	for(int i=1;i<=m;++i)cout<<ans[i]<<' ';

	cout<<'\n';

}