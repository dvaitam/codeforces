#include<cstdio>

#include<cstring>

#include<iostream>

#include<algorithm>

#include<vector>

#include<cmath>

#define pb push_back

#define ll long long

#define inf 0x3f3f3f3f

#define fill(a,b) memset(a,b,sizeof(a))

#define cpy(a,b) memcpy(a,b,sizeof(b))

template<typename T> inline bool MIN(T &a,const T &b) {return a>b? a=b,1:0;}

template<typename T> inline bool MAX(T &a,const T &b) {return a<b? a=b,1:0;}

#define mp make_pair

#define fi first

#define se second

#define N 100010

#define eps (1e-8)

#define db long double

using namespace std;

typedef vector<int> vi;

typedef pair<db,int> pii;

typedef pair<int,pii> piii;



int n,l0;



struct Point {

	db x,y;

	Point() {}

	Point(db x,db y):x(x),y(y) {}

}A[N];

Point operator + (Point a,Point b) {return Point(a.x+b.x,a.y+b.y);} 

Point operator - (Point a,Point b) {return Point(a.x-b.x,a.y-b.y);} 

Point operator * (Point a,db b) {return Point(a.x*b,a.y*b);}

Point operator / (Point a,db b) {return Point(a.x/b,a.y/b);}



db length(Point a) {return sqrt(a.x*a.x+a.y*a.y);}



pii q[N<<1];



int vis[N],S[N];



int run(Point a,Point b,db r1,db r2,Point &c1,Point &c2) {

	if (r1>r2) swap(a,b),swap(r1,r2);

	db d=length(a-b);

	if (r1+r2<d+eps||d+r1<r2+eps) return 0;

	db t=(r1*r1-r2*r2+d*d)/2/d;

	Point c=a+(b-a)/length(b-a)*t;

//	printf("t %.10lf %.10lf\n",(double)t,(double)r1);

	t=sqrt(r1*r1-t*t);

	

	Point v=a-b; swap(v.x,v.y); v.x=-v.x; v=v/length(v);

	c1=c+v*t,c2=c-v*t;

	return 1;

}





int ok(db mid) {

	int cnt_q=0,cnt_S=0;

	fill(vis,0);

	for (int i=1;i<=n;++i) {

		Point a=A[i],b=Point(l0,0);

		db r1=length(a-Point(-l0,0)),r2=mid;

		Point c1,c2;

		if (!run(a,b,r1,r2,c1,c2)) continue;

		



		q[++cnt_q]=mp(atan2(c1.y,c1.x-l0),i);

		q[++cnt_q]=mp(atan2(c2.y,c2.x-l0),i);

//		printf("c1 %.10lf %.10lf\n",(double)c1.x,(double)c1.y);

//		printf("c2 %.10lf %.10lf\n",(double)c2.x,(double)c2.y);

		

//		printf("atan2 %lf %lf\n",(double)atan2(c1.y,c1.x-l0),(double)atan2(c2.y,c2.x-l0));

//		printf("%lf %lf\n",);

	}

	sort(q+1,q+cnt_q+1);

	for (int i=1;i<=cnt_q;++i) {

		if (!vis[q[i].se]) vis[q[i].se]=1,S[++cnt_S]=q[i].se;

		else {

			if (q[i].se!=S[cnt_S]) return 1;

			--cnt_S;

		}

	}

	return 0;

}



int main()

{ 

//	freopen("A.in","r",stdin);

	scanf("%d%d",&n,&l0);

	for (int i=1;i<=n;++i) {

		int a,b; scanf("%d%d",&a,&b);

		A[i]=Point(a,b);

	

	}

//	

	db l=0,r=l0*2;

	for (int i=0;i<100;++i) {

		db mid=(l+r)/2;

		if (ok(mid)) r=mid; 

		else l=mid;

	}

	printf("%.12lf\n",(double)l);

//	printf("%d\n",ok(0.0000018));



//	for (double i=1;i<10;i+=0.1) {

//		printf("%lf %d\n",i,ok(i));

//	}

//	printf("%d\n",ok(6));

//	

//	Point c1,c2;

//	run(Point(2,0),Point(3,0),2,1.1,c1,c2);

//	

//	printf("%lf %lf\n",c1.x,c1.y);

//	printf("%lf %lf\n",c2.x,c2.y);

	return 0;

}