#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <cmath>
#include <iostream>
#include <queue>
#include <set>
#include <algorithm>

using namespace std;

typedef long long ll;
#define MEM(s,v) memset(s, v, sizeof(s))
#define inf 1e9
#define eps 1e-9
#define _N 100005
#define _M 305
#define PI acos(-1.0)
#define zero(x) (((x)>0?(x):-(x))<eps)

struct point
{
	double x, y;
}pnt[_M], res[_M];

double xmult(point p1,point p2,point p0){
	return (p1.x-p0.x)*(p2.y-p0.y)-(p2.x-p0.x)*(p1.y-p0.y);
}

int dots_inline(point p1,point p2,point p3){
	return zero(xmult(p1,p2,p3));
}
int same_side(point p1,point p2,point l1,point l2){
	return xmult(l1,p1,l2)*xmult(l1,p2,l2)>eps;
}
int dot_online_in(point p,point l1,point l2){
	return zero(xmult(p,l1,l2))&&(l1.x-p.x)*(l2.x-p.x)<eps&&(l1.y-p.y)*(l2.y-p.y)<eps;
}
int intersect_in(point u1,point u2,point v1,point v2){
	if (!dots_inline(u1,u2,v1)||!dots_inline(u1,u2,v2))
		return !same_side(u1,u2,v1,v2)&&!same_side(v1,v2,u1,u2);
	return dot_online_in(u1,v1,v2)||dot_online_in(u2,v1,v2)||dot_online_in(v1,u1,u2)||dot_online_in(v2,u1,u2);
}

int n;

double area(point p[], int nn)
{
	double s = 0;
	//p[n] = p[0];
	//long long s = 0;
	for (int i = 1; i < nn-1; ++i)
		s += xmult(p[i], p[i+1], p[0]);	//p[0]==[n-1]
	return fabs(s / 2);
}

/*bool ok(int i1, int i2, int i3, int i4)
{
	point a[5];
	a[0] = p[i1], a[1] = p[i2], a[2] = p[i3], a[3] = p[i4];
	if (!intersect_in(a[0], a[1], a[2], a[3]) && !intersect_in(a[0], a[3], a[1], a[2]))
		return 1;
	return 0;
}*/

bool operator<(const point &a,const point &b)
{
    return a.y<b.y || (a.y==b.y && a.x<b.x);
}

bool operator==(const point &a,const point &b)
{
	if (fabs(a.x-b.x)<eps && fabs(a.y-b.y)<eps)
		return 1;
	return 0;
}
/*
折线的拐向的判断（从op向sp看过去的左边）
若 (ep - sp) 叉乘 (sp - op) < 0 ,则p0p1在p1点拐向左侧后得到p1p2
若 (ep - sp) 叉乘 (sp - op) = 0 ,则 p0, p1, p2三点共线
若 (ep - sp) 叉乘 (sp - op) > 0 ,则p0p1在p1点拐向右侧后得到p1p2
*/

bool mult(point sp,point ep,point op)
{
    return (sp.x-op.x)*(ep.y-op.y)>=(ep.x-op.x)*(sp.y-op.y);///>=
}

bool vis[_M];

int graham(int n)
{
    int i,len,top=1;
    sort(pnt,pnt+n);
    res[0] = pnt[0];
    res[1] = pnt[1];
    for(i=2; i<n; i++)              ///从i=2开始！
    {
        while(top && mult(pnt[i],res[top],res[top-1]))  ///res[top-1]是假定为凸包的点
            top--;
        res[++top] = pnt[i];
    }
    len = top;
    res[++top]=pnt[n-2];            ///attention!!!
    for(i=n-3; i>=0; i--)           ///判断是否和第一个点形成凸包
    {
        while(top!=len && mult(pnt[i],res[top],res[top-1]))
            top--;
        res[++top] = pnt[i];
    }
	MEM(vis, 0);
	for (i = 0; i < top; i ++)
		for (int j = 0; j < n; j ++)
			if (pnt[j] == res[i])
				vis[j] = 1;
//    cout<<"------------------\n";
//    for(i=0; i<top; i++)
//        cout<<res[i].x<<' '<<res[i].y<<endl;
//    cout<<"------------------\n";
    return top;
}

double area_triangle(point p1,point p2,point p3){
	return fabs(xmult(p1,p2,p3))/2;
}

double cal(int id, int nc)
{
	double re = inf;
	res[nc] = res[0];
	for (int i = 0; i < nc; i ++)
	{
		re = min(re, area_triangle(pnt[id], res[i], res[i+1]));
	}
	return re;
}

int main()
{
	int i, j, k, p;
	while (~scanf("%d", &n))
	{
		for (i = 0; i < n; i ++)
			scanf("%lf%lf", &pnt[i].x, &pnt[i].y);
		point tmp[5];
		double ans = -1;
		int c = graham(n);
		double base = area(res, c);//cout<<base<<' '<<c<<endl;
		//for (i = 0; i < c; i ++)
			//cout<<res[i].x<<','<<res[i].y<<endl;
		if (c < 4)
		{
			for (i = 0; i < n;  i++)
				if (!vis[i])
					ans = max(ans, base-cal(i, c));
		}
		else
		{
			for (i = 0; i < c-3; i ++)
				for (j = i+1; j < c-2; j ++)
					for (k = j+1; k < c-1; k ++)
						for (p = k+1; p < c; p ++)
						{
							tmp[0]=res[i],tmp[1]=res[j],tmp[2]=res[k],tmp[3]=res[p];//cout<<area(tmp, 4)<<endl;
							ans = max(ans, area(tmp, 4));
						}
		}
		printf("%.9lf\n", ans);
	}
	return 0;
}