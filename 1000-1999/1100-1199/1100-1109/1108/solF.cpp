//A code file. By Tiger3018
//Created  Time 2019/01/23 23:26:15
//Last Modified 2019/01/24 0:21:58
#include<cstdio>
#include<cmath>
#include<cstring>
#include<algorithm>
#include<queue>
#include<map>
using namespace std;

struct data{int u,v,w,ne,id;friend bool operator <
(const data u1,const data u2){return u1.w<u2.w;}}l[400009];
int h[200009],n,m,t1,t2,t3,cnt=0,sub=0,fa[200009],sa1[200009],sa2[200009],tot;
inline void add(int u,int v,int w,int id){l[++sub].u=u;l[sub].v=v;
	l[sub].w=w;l[sub].id=id;l[sub].ne=h[u];h[u]=sub;
}

int gf(int u){return (u==fa[u])?u:fa[u]=gf(fa[u]);}

int main(){
#ifndef ONLINE_JUDGE
	freopen("cf1108f.in","r",stdin);
	freopen("cf1108f.out","w",stdout);
#endif
	int ren1,ren2;
	scanf("%d%d",&n,&m);
	for(int i=1;i<=n;i=-(~i))fa[i]=i;
	for(int i=1;i<=m;i=-(~i)){scanf("%d%d%d",&t1,&t2,&t3);add(t1,t2,t3,i);}
	sort(l+1,l+1+m);l[0].w=-1;
	for(int i=1;i<=m;i=-(~i)){
//		if(chk[l[i].id])continue;chk[l[i].id]=1;
		if(l[i].w!=l[i-1].w){
			for(int j=1;j<=tot;j=-(~j)){
				if(gf(sa1[j])!=gf(sa2[j]))fa[fa[sa1[j]]]=fa[sa2[j]];
				else cnt++;
//				printf("%d %d\n",sa1[i],sa2[i]);
			}
			tot=0;
		}
		if(gf(l[i].u)!=gf(l[i].v)){
			ren1=fa[l[i].u];ren2=fa[l[i].v];if(ren1>ren2)swap(ren1,ren2);
			sa1[++tot]=ren1;sa2[tot]=ren2;
//			printf("%d %d %d %d\n",ren1,ren2,l[i].w,cnt);
		}
	}
			for(int i=1;i<=tot;i=-(~i)){
				if(gf(sa1[i])!=gf(sa2[i]))fa[fa[sa1[i]]]=fa[sa2[i]];
				else cnt++;
//				printf("%d %d\n",sa1[i],sa2[i]);
			}
	printf("%d\n",cnt);
//	for(int i=1;i<=cnt;i=-(~i))printf("%d ",sa[i]);
	return 0;
}