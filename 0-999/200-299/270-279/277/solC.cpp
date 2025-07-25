#include<cstdio>
#include<cstring>
#include<algorithm>
using namespace std;
const int tt=100005;
struct vl{
	int a,b,c;
	vl(){}
	vl(int aa,int bb,int cc){
		a=aa,b=bb,c=cc;
	}
	bool operator <(const vl &y) const{
		return c<y.c || c==y.c && a<y.a;
	}
} x[tt],y[tt];
struct gs{
	int ps,sg;
	gs(){}
	gs(int p,int g){
		ps=p,sg=g;
	}
} px[tt],py[tt];
int main(){
	int n,m,k,i,a,b,c,d,k1=0,k2=0,t1=0,t2=0,sg,lst,xx,yy,mi,til;
//	freopen("c.in","r",stdin);
//	freopen("c.out","w",stdout);
	scanf("%d%d%d",&n,&m,&k);
	for (i=0;i<k;i++){
		scanf("%d%d%d%d",&a,&b,&c,&d);
		if (a>c) swap(a,c);
		if (b>d) swap(b,d);
		if (b==d) x[k1++]=vl(a,c,b); else
			y[k2++]=vl(b,d,a);
	}
	sort(x,x+k1);
	sort(y,y+k2);
	sg=(n%2?0:m)^(m%2?0:n);
	for (lst=xx=0,i=0;i<=k1;i++){
		if (i==k1 && lst!=m-1 || i!=k1 && x[i].c>lst+1) xx=lst+1;
		if (i==k1 || x[i].c!=lst){
			if (i){
				sg^=n-mi;
				px[t1++]=gs(x[i-1].c,n-mi);
			}
			til=mi=0;
			if (i==k1) break;
			sg^=n;
		}
		if (x[i].b>til){
			mi+=x[i].b-max(til,x[i].a);
			til=x[i].b;
		}
		lst=x[i].c;
	}
	for (lst=yy=0,i=0;i<=k2;i++){
		if (i==k2 && lst!=n-1 || i!=k2 && y[i].c>lst+1) yy=lst+1;
		if (i==k2 || y[i].c!=lst){
			if (i){
				sg^=m-mi;
				py[t2++]=gs(y[i-1].c,m-mi);
			}
			til=mi=0;
			if (i==k2) break;
			sg^=m;
		}
		if (y[i].b>til){
			mi+=y[i].b-max(til,y[i].a);
			til=y[i].b;
		}
		lst=y[i].c;
	}
	if (xx) px[t1++]=gs(xx,n);
	if (yy) py[t2++]=gs(yy,m);
	if (sg==0) puts("SECOND"); else{
		puts("FIRST");
		for (i=0;i<t1;i++) if ((sg^px[i].sg)<px[i].sg) break;
		if (i==t1){
			for (i=0;i<t2;i++) if ((sg^py[i].sg)<py[i].sg) break;
			yy=py[i].ps,lst=py[i].sg-(sg^py[i].sg);
			printf("%d 0 ",yy);
			for (til=0,i=0;i<k2;i++) if (y[i].c==yy){
				if (y[i].a>til){
					if (lst<=y[i].a-til) break;
					else lst-=y[i].a-til;
				}
				if (y[i].b>til) til=y[i].b;
			}
			printf("%d %d\n",yy,til+lst);
		} else{
			xx=px[i].ps,lst=px[i].sg-(sg^px[i].sg);
			printf("0 %d ",xx);
			for (til=0,i=0;i<k1;i++) if (x[i].c==xx){
				if (x[i].a>til){
					if (lst<=x[i].a-til) break;
					else lst-=x[i].a-til;
				}
				if (x[i].b>til) til=x[i].b;
			}
			printf("%d %d\n",til+lst,xx);
		}
	}
	return 0;
}