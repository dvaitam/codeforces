#include<iostream>
#include<cstdio>
#include<algorithm>
#include<vector>
#include<map>
#include<cstdlib>
#include<cstring>
#include<string>
#define maxn 101000

using namespace std;

typedef long long ll;

int init(int x,pair<int,int> p){
	return x<p.second && x>p.first;
}

bool cross(pair<int,int> p1,pair<int,int> p2){
	return (init(p1.first,p2) || init(p1.second,p2)) && (init(p2.first,p1) || init(p2.second,p1));
}


int a[110];
bool f[110][110];
bool flag;
int m;
pair<int,int> p[110];
void dfs(int x,int y){
	if(a[x]!=-1){
		if(a[x]!=y)flag=false;
		return;
	}
		
	if(!flag)return;
	a[x]=y;
	for(int i=1;i<=m;++i)if(f[x][i])dfs(i,1-y);
}

int main(){
	int n,i,j;
	scanf("%d%d",&n,&m);
	for(i=1;i<=m;++i){
		scanf("%d%d",&p[i].first,&p[i].second);
		if(p[i].first>p[i].second)swap(p[i].first,p[i].second);
	}
	for(i=1;i<=m;++i)
		for(j=1;j<=m;++j)
			if(cross(p[i],p[j]))f[i][j]=true;
	flag=true;
	for(i=1;i<=m;++i)a[i]=-1;
	for(i=1;i<=m;++i)if(a[i]==-1)dfs(i,0);
	if(flag){
		for(i=1;i<=m;++i)if(a[i]==0)printf("o");else printf("i");
		printf("\n");
	}else {
		printf("Impossible\n");
	}
	return 0;
}