#include<cstdio>
#define sqr(x) ((x)*(x))
double cp(double x0,double y0,double x1,double y1,double x2,double y2)
{	return (x1-x0)*(y2-y0)-(x2-x0)*(y1-y0);
}
int check(double x1,double y1,double x2,double y2,double x3,double y3)
{	double X1,Y1,X2,Y2,X3,Y3,X4,Y4,v1,v2,v3,v4,a1,b1,c1,a2,b2,c2;
	a1=(x2-x1)*2;
	b1=(y2-y1)*2;
	c1=sqr(2*x1-x2)+sqr(2*y1-y2)-sqr(x1)-sqr(y1);
	a2=(x3-2*x2+x1)*2;
	b2=(y3-2*y2+y1)*2;
	c2=sqr(x1)+sqr(y1)-sqr(x3-2*x2+2*x1)-sqr(y3-2*y2+2*y1);
	if (a1*b2==a2*b1)
		return 0;
	Y1=(c2*a1-c1*a2)/(b1*a2-b2*a1);
	X1=(c2*b1-c1*b2)/(a1*b2-a2*b1);
	X2=2*x1-X1;
	Y2=2*y1-Y1;
	X3=2*x2-2*x1+X1;
	Y3=2*y2-2*y1+Y1;
	X4=2*x3-2*x2+2*x1-X1;
	Y4=2*y3-2*y2+2*y1-Y1;
	v1=cp(X1,Y1,X2,Y2,X3,Y3);
	v2=cp(X2,Y2,X3,Y3,X4,Y4);
	v3=cp(X3,Y3,X4,Y4,X1,Y1);
	v4=cp(X4,Y4,X1,Y1,X2,Y2);
	if ((v1<0&&v2<0&&v3<0&&v4<0)||(v1>0&&v2>0&&v3>0&&v4>0))
	{	printf("YES\n%.9lf %.9lf %.9lf %.9lf %.9lf %.9lf %.9lf %.9lf\n",X1,Y1,X2,Y2,X3,Y3,X4,Y4);
		return 1;
	}
	return 0;
}
int main()
{	int t;
	double x1,y1,x2,y2,x3,y3;
	scanf("%d",&t);
	while (t--)
	{	scanf("%lf%lf%lf%lf%lf%lf",&x1,&y1,&x2,&y2,&x3,&y3);
		if ((check(x1,y1,x2,y2,x3,y3)||check(x1,y1,x3,y3,x2,y2)||check(x2,y2,x1,y1,x3,y3))==0)
			printf("NO\n\n");
	}
}