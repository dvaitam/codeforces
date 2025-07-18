#include<iostream>

#include<cstdio>

#include<algorithm>

#include<string.h>

#include<cmath>

#include<stack>

using namespace std;

const int maxn=500005;

const int maxv=100005;

int n,t,s,d,i,j,k,u,v,cost,ans,ls,m,cntege,lt,z;

int tem[maxv];

int to[maxn],nxt[maxn],adj[maxv],ne;

bool vis[maxn];

int apr[maxv];

bool vise[maxn];

int posj[maxv];



void addege(int u,int v){

	to[ne]=v;

	nxt[ne]=adj[u];

	adj[u]=ne;

	ne++;

}

void dfs(int s){

	int j;//printf("s %d\n",s);

	for(;posj[s]!=-1;posj[s]=nxt[posj[s]]){

		j=posj[s];//printf("pos1 pos2 %d %d\n",posj[1],posj[2]);if(j==0)break;

		if(!vis[j]){

			vis[j]=true;

			vis[j^1]=true;

			dfs(to[j]);

			if(++z%2){

				vise[j]=true;

			}

			else{

				vise[j^1]=true;

			}

		}

		if(posj[s]==-1)break;

	}

}

int main(){



	memset(adj,-1,sizeof(adj));

	ne=2;cntege=0;

	scanf("%d %d",&n,&m);

	for(i=1;i<=m;i++){

		scanf("%d %d",&u,&v);

		addege(u,v);

		addege(v,u);

		apr[u]++;

		apr[v]++;

		cntege++;

	}

	lt=0;

	for(i=1;i<=n;i++){

		if(apr[i]%2)

			tem[++lt]=i;

			

	}//printf("%d\n",nxt[2]);

	for(i=1;i<=lt;i+=2){

		if(i==lt){

			addege(tem[i],tem[i]);

			addege(tem[i],tem[i]);

		}

		else{

			addege(tem[i],tem[i+1]);

			addege(tem[i+1],tem[i]);			

		}

		cntege++;

	}z=0;

//for(i=1;i<=n;i++)for(;posj[i]!=-1;posj[i]=nxt[posj[i]])printf("i j %d %d\n",i,to[posj[i]]);

	for(i=1;i<=n;i++)posj[i]=adj[i];

	dfs(1);

	//



	if(cntege%2){

		addege(1,1);

		addege(1,1);

		vise[ne-1]=true;

		cntege++;

	}



	

	printf("%d\n",cntege);

	for(i=1;i<=n;i++){

		for(j=adj[i];j!=-1;j=nxt[j]){

			if(vise[j])printf("%d %d\n",i,to[j]);

		}

	}



	return 0;

}