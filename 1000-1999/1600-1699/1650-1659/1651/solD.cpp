#include<bits/stdc++.h>
using namespace std;
int n,ma[200001],mi[200001];
int asw[200001],x[200001],y[200001];
struct nood {
	int x,y,id;
} g[200005];
bool cmp1(nood a,nood b) {
	if(a.y==b.y) return a.x<b.x;
	return a.y<b.y;
}
bool cmp2(nood a,nood b) {
	if(a.y==b.y) return a.x>b.x;
	return a.y<b.y;
}
int main() {
	scanf("%d",&n);
	for(int i=1; i<=n; i++) scanf("%d%d",&g[i].x,&g[i].y),g[i].id=i;
	sort(g+1,g+1+n,cmp1);
	for(int xm,ASW,X,Y,i=1; i<=n; i++) {
		if(g[i].y!=g[i-1].y) ASW=1e8,xm=g[i].x-1;
		else {
			ASW+=g[i].x-g[i-1].x;
			if(g[i-1].x+1!=g[i].x) xm=g[i].x-1;
		}
		asw[g[i].id]=g[i].x-xm,x[g[i].id]=xm,y[g[i].id]=g[i].y;
		if(ma[g[i].x]+1!=g[i].y) mi[g[i].x]=g[i].y-1;
		ma[g[i].x]=g[i].y;
		if(g[i].y-mi[g[i].x]<ASW) {
			ASW=g[i].y-mi[g[i].x],X=g[i].x,Y=mi[g[i].x];
		}
		if(ASW<asw[g[i].id]) asw[g[i].id]=ASW,x[g[i].id]=X,y[g[i].id]=Y;
	}
	memset(ma,0,sizeof(ma)),memset(mi,0,sizeof(mi));
	for(int xm,ASW,X,Y,i=n; i; i--) {
		if(g[i].y!=g[i+1].y) ASW=1e8,xm=g[i].x+1;
		else {
			ASW+=g[i+1].x-g[i].x;
			if(g[i+1].x-1!=g[i].x) xm=g[i].x+1;
		}
		if(xm-g[i].x<asw[g[i].id]) asw[g[i].id]=xm-g[i].x,x[g[i].id]=xm,y[g[i].id]=g[i].y;
		if(ma[g[i].x]-1!=g[i].y) mi[g[i].x]=g[i].y+1;
		ma[g[i].x]=g[i].y;
		if(mi[g[i].x]-g[i].y<ASW) {
			ASW=mi[g[i].x]-g[i].y,X=g[i].x,Y=mi[g[i].x];
		}
		if(ASW<asw[g[i].id]) asw[g[i].id]=ASW,x[g[i].id]=X,y[g[i].id]=Y;
	}
	sort(g+1,g+1+n,cmp2);
	memset(ma,0,sizeof(ma)),memset(mi,0,sizeof(mi));
	for(int xm,ASW,X,Y,i=1; i<=n; i++) {
		if(g[i].y!=g[i-1].y) ASW=1e8,xm=g[i].x+1;
		else {
			ASW+=g[i-1].x-g[i].x;
			if(g[i-1].x-1!=g[i].x) xm=g[i].x+1;
		}
		if(xm-g[i].x<asw[g[i].id]) asw[g[i].id]=xm-g[i].x,x[g[i].id]=xm,y[g[i].id]=g[i].y;
		if(ma[g[i].x]+1!=g[i].y) mi[g[i].x]=g[i].y-1;
		ma[g[i].x]=g[i].y;
		if(g[i].y-mi[g[i].x]<ASW) {
			ASW=g[i].y-mi[g[i].x],X=g[i].x,Y=mi[g[i].x];
		}
		if(ASW<asw[g[i].id]) asw[g[i].id]=ASW,x[g[i].id]=X,y[g[i].id]=Y;
	}
	memset(ma,0,sizeof(ma)),memset(mi,0,sizeof(mi));
	for(int xm,ASW,X,Y,i=n; i; i--) {
		if(g[i].y!=g[i+1].y) ASW=1e8,xm=g[i].x-1;
		else {
			ASW+=g[i].x-g[i+1].x;
			if(g[i+1].x+1!=g[i].x) xm=g[i].x-1;
		}
		if(g[i].x-xm<asw[g[i].id]) asw[g[i].id]=g[i].x-xm,x[g[i].id]=xm,y[g[i].id]=g[i].y;
		if(ma[g[i].x]-1!=g[i].y) mi[g[i].x]=g[i].y+1;
		ma[g[i].x]=g[i].y;
		if(mi[g[i].x]-g[i].y<ASW) {
			ASW=mi[g[i].x]-g[i].y,X=g[i].x,Y=mi[g[i].x];
		}
		if(ASW<asw[g[i].id]) asw[g[i].id]=ASW,x[g[i].id]=X,y[g[i].id]=Y;
	}
	for(int i=1; i<=n; i++) printf("%d %d\n",x[i],y[i]);
}