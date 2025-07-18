#include<bits/stdc++.h>
using namespace std;
#define pc putchar

const int D[4][2]={{1,0},{0,1},{0,-1},{-1,0}};
struct node{int x,y;bool in(){return 0<=x&&x<=410&&0<=y&&y<=410;}}a[410],P[4],s,t,q[514*514];
bool operator==(const node&a,const node&b){return a.x==b.x&&a.y==b.y;}
int n,p[514*514],fr[514*514],pre[514*514];
bitset<514>c[514],v[514];

bool bfs(){
	q[1]=s,v[s.x][s.y]=1;
	for(int l=1,r=1;l<=r;++l)for(int i=0;i^4;++i){
		node u={q[l].x+D[i][0],q[l].y+D[i][1]};
		if(u.in()&&!v[u.x][u.y]){
			q[++r]=u,fr[r]=i,pre[r]=l,v[u.x][u.y]=1;
			if(u==t){
				for(;r>1;r=pre[r])p[++p[0]]=fr[r];
				reverse(p+1,p+p[0]+1);
				return 1;
			}
		}
	}
	return 0;
}

bool move(int i){
	s.x+=D[i][0],s.y+=D[i][1],t.x+=D[i][0],t.y+=D[i][1],pc("RUDL"[i]);
	if(t.in()&&c[t.x][t.y]){t.x-=D[i][0],t.y-=D[i][1];return 0;}
	return 1;
}

void getout(){
	while((0<=s.x&&s.x<=410)||(0<=t.x&&t.x<=410))move(0);
	while((0<=s.y&&s.y<=410)||(0<=t.y&&t.y<=410))move(1);
}
void ud(int d){
	getout();
	while((d==2&&t.y>=0)||(d==1&&t.y<=410))move(d);
	while(t.x^P[d].x)move(t.x<P[d].x?0:3);
	while(t.y^s.y)move(3-d);
}
void lr(int d){
	getout();
	while((d==3&&t.x>=0)||(d==0&&t.x<=410))move(d);
	while(t.y^P[d].y)move(t.y<P[d].y?1:2);
	while(t.x^s.x)move(3-d);
}

int main(){
	ios::sync_with_stdio(false),cin.tie(0),cout.tie(0);
	cin>>s.x>>s.y>>t.x>>t.y>>n;
	if(!n)return pc('-'),pc('1'),pc('\n'),0;
	s.x+=200,s.y+=200,t.x+=200,t.y+=200,P[0].x=P[1].y=-1e9,P[2].y=P[3].x=1e9;
	for(int i=1;i<=n;i++){
        cin>>a[i].x>>a[i].y;
        a[i].x+=200,a[i].y+=200;
        v[a[i].x][a[i].y]=c[a[i].x][a[i].y]=1;
        P[3]=a[i].x<P[3].x?a[i]:P[3];
        P[0]=a[i].x>P[0].x?a[i]:P[0];
        P[2]=a[i].y<P[2].y?a[i]:P[2];
        P[1]=a[i].y>P[1].y?a[i]:P[1];
    }
    if(!bfs())return pc('-'),pc('1'),pc('\n'),0;
	for(int i=1;i<=p[0];i++){
		if(move(p[i]))p[++p[0]]=p[i];
		if(s==t)return pc('\n'),0;
		if(!s.in())break;
	}
    if(s.y^t.y)ud(1+(bool)(s.y<t.y));
	if(s.x^t.x)lr(3*(bool)(s.x<t.x));
	return pc('\n'),0;
}/*1694643741.2081487*/