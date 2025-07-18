#include<cmath>
#include<cstdio>
#include<iostream>
#include<cstring>
#include<algorithm>
#include<cstring>
#include<vector>
using namespace std;
typedef long long ll;
int getint(){
	int x=0,f=1;
	char ch=getchar();
	while(ch<'0'||'9'<ch){if(ch=='-')f=-1;ch=getchar();}
	while('0'<=ch&&ch<='9'){x=x*10+ch-'0';ch=getchar();}
	return f*x;
}
#define N 1123456
#define mod 1000000007
int s[N];
vector<int> a[N];

int head[N],etot;
struct edge{
	int to,next;
}e[N];
void insert(int x,int y){
	e[++etot].to=y; e[etot].next=head[x]; head[x]=etot;
	e[++etot].to=x; e[etot].next=head[y]; head[y]=etot;
}
int fa[N];
int used[N];
int col[N];
void dfs(int x){
	for(int i=0;i<a[x].size();i++) {
		int y=a[x][i]; if(col[y]) used[col[y]]=x;	
	}
	int j=1;
	for(int i=0;i<a[x].size();i++ ){
		int y=a[x][i]; if(col[y]) continue;
		while(used[j]==x) j++;
		col[y]=j; j++;
	}
	for(int i=head[x];i;i=e[i].next){
		int y=e[i].to;
		if(y!=fa[x]){
			fa[y]=x; dfs(y); 
		}
	}
	return ;
}

int main(){
	int n=getint();
	int m=getint();
	int cnt=0,root=0;
	for(int i=1;i<=n;i++) {
		s[i]=getint();
		if(s[i]>cnt) root=i,cnt=s[i];
		for(int j=1;j<=s[i];j++) a[i].push_back(getint());
	}
	for(int i=1;i<n;i++) insert(getint(),getint());
	dfs(root);
	if(cnt==0) cnt=1;
	cout<<cnt<<endl;
	for(int i=1;i<=m;i++) {
		if(!col[i]) col[i]=1;
		printf("%d ",col[i]);
	}
	return 0;
}