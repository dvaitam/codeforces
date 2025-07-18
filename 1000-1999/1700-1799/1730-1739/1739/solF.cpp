#include<bits/stdc++.h>

using namespace std;

const int maxn=4005,S=(1<<12)-1;

int n,m,T,vs,es,tot,finf,ans;

int c[maxn],f[1<<12][maxn],deg[15],res[15],vis[15],mp[15][15],nxt[maxn][15],ed[maxn],fail[maxn],rec[1<<12][maxn][2];

string s;

vector<int>now;

vector<int>v[15];

queue<int>q;

inline void add(int x,int y){

	if(mp[x][y]==0)

		mp[x][y]=mp[y][x]=1,deg[x]++,deg[y]++,v[x].emplace_back(y),v[y].emplace_back(x),es++;

}

void dfs(int x){

	now.emplace_back(x),vis[x]=1;

	for(int i=0;i<v[x].size();i++)

		if(vis[v[x][i]]==0)

			dfs(v[x][i]);

}

void insert(int c){

	int p=1;

	for(int i=0;i<now.size();i++){

		int x=now[i];

		if(nxt[p][x]==0)

			nxt[p][x]=++tot;

		p=nxt[p][x];

	}

	ed[p]+=c;

}

int main(){

	scanf("%d",&n),tot=1;

	for(int i=1;i<=n;i++){

		scanf("%d",&c[i]);

		cin>>s;

		for(int j=0;j<s.size();j++){

			int c=s[j]-96;

			vs+=vis[c]==0,vis[c]=1;

			if(j>0)

				add(s[j-1]-96,s[j]-96);

		}

		int flg=vs==es+1,rt=-1;

		for(int j=1;j<=12;j++){

			if(deg[j]==1)

				rt=j;

			flg&=deg[j]<=2;

		}

		for(int j=1;j<=12;j++)

			vis[j]=0;

		if(flg)

			now.clear(),dfs(rt),insert(c[i]),reverse(now.begin(),now.end()),insert(c[i]);

		for(int j=1;j<=12;j++)

			deg[j]=vis[j]=0,v[j].clear();

		es=vs=0;

		memset(mp,0,sizeof(mp));

	}

	for(int i=1;i<=12;i++){

		if(nxt[1][i])

			q.push(nxt[1][i]),fail[nxt[1][i]]=1;

		else nxt[1][i]=1;

	}

	while(!q.empty()){

		int x=q.front();

		ed[x]+=ed[fail[x]],q.pop();

		for(int i=1;i<=12;i++){

			if(nxt[x][i])

				q.push(nxt[x][i]),fail[nxt[x][i]]=nxt[fail[x]][i];

			else nxt[x][i]=nxt[fail[x]][i];

		}

	}

	memset(f,128,sizeof(f)),finf=f[0][1],f[0][1]=0;

	for(int i=0;i<S;i++)

		for(int j=1;j<=tot;j++)

			if(f[i][j]!=finf)

				for(int c=1;c<=12;c++)

					if(((i>>(c-1))&1)==0){

						int s=i|(1<<(c-1)),t=nxt[j][c];

						if(f[s][t]<f[i][j]+ed[t])

							f[s][t]=f[i][j]+ed[t],rec[s][t][0]=i,rec[s][t][1]=j;

					}

	int x=S,y,pos=12;

	ans=finf;

	for(int i=1;i<=tot;i++)

		if(f[S][i]>ans)

			ans=f[S][i],y=i;

//	printf("ans=%d\n",ans);

	while(x){

		int nx=rec[x][y][0],ny=rec[x][y][1];

		res[pos]=__builtin_ctz(x^nx)+1,pos--;

		x=nx,y=ny;

	}

	for(int i=1;i<=12;i++)

		printf("%c",res[i]+96);

	return 0;

}