#include <cstdio>
#include <cmath>
const double eps=1e-4;
const double PI=acos(-1.0);
double gcd(double x,double y)
{
    return y>eps? gcd(y,x-floor(x/y)*y):x;
}

double bcos(double a,double b,double c)
{
     return acos((a*a+b*b-c*c)/(2*a*b));
}
int main()
{
    double ax,ay,bx,by,cx,cy;
    scanf("%lf%lf%lf%lf%lf%lf",&ax,&ay,&bx,&by,&cx,&cy);
    double a=sqrt((ax-bx)*(ax-bx)+(ay-by)*(ay-by));
    double b=sqrt((ax-cx)*(ax-cx)+(ay-cy)*(ay-cy));
    double c=sqrt((bx-cx)*(bx-cx)+(by-cy)*(by-cy));
    double p=(a+b+c)/2;
    double s=sqrt(p*(p-a)*(p-b)*(p-c));
    double R=(a*b*c)/(4*s);
    double A=bcos(b,c,a);
    double B=bcos(a,c,b);
    double C=bcos(a,b,c);
    double n=PI/gcd(A,gcd(B,C));
    printf("%.11lf\n",R*R*sin(2*PI/n)*n/2);
    return 0;
}