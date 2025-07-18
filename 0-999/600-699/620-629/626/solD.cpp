#include <bits/stdc++.h>
using namespace std;
typedef long long LL;
const double pi=acos(-1.0);
struct Complex
{
    double x,y;
    Complex(double xx=0,double yy=0):x(xx),y(yy){}
    Complex operator+(const Complex&tt)const
    {
        return Complex(x+tt.x,y+tt.y);
    }
    Complex operator-(const Complex&tt)const
    {
        return Complex(x-tt.x,y-tt.y);
    }
    Complex operator*(const Complex&tt)const
    {
        return Complex(x*tt.x-y*tt.y,y*tt.x+x*tt.y);
    }
}a[32768];
int len,n,in[2005];
LL ans[10005],tmp[10005];
void change(Complex*a,int len)
{
    int i,j,t,k=len>>1;
    for(i=1,j=k;i<len-1;++i)
    {
        if(i<j)
            swap(a[i],a[j]);
        for(t=k;j>=t&&t>=1;j-=t,t>>=1);
        j+=t;
    }
}
void fft(Complex*a,int len,int f)
{
    change(a,len);
    int i,j,k;
    for(int d=2;d<=len;d<<=1)
    {
        Complex wn=Complex(cos(-2*pi*f/d),sin(-2*pi*f/d));
        for(i=0;i<len;i+=d)
        {
            Complex w=Complex(1.0);
            for(j=0;j<(d>>1);++j)
            {
                Complex p1=a[i+j];
                Complex p2=w*a[i+j+(d>>1)];
                a[i+j]=p1+p2;
                a[i+j+(d>>1)]=p1-p2;
                w=w*wn;
            }
        }
    }
    if(f==-1)
    {
        for(i=0;i<len;++i)
            a[i].x/=len;
    }
}
LL gcd(LL x,LL y)
{
    LL t;
    while(y)
    {
        t=x;
        x=y;
        y=t%y;
    }
    return x;
}
int main()
{
    int i,j;
    scanf("%d",&n);
    for(i=1;i<=n;++i)
        scanf("%d",in+i);
    for(i=1;i<n;++i)
        for(j=i+1;j<=n;++j)
            ++tmp[abs(in[i]-in[j])];
    for(len=1;len<=20000;len<<=1);
    for(i=0;i<5000;++i)
        a[i].x=tmp[i];
    fft(a,len,1);
    for(i=0;i<len;++i)
        a[i]=a[i]*a[i];
    fft(a,len,-1);
    for(i=1;i<=10000;++i)
        ans[i]=(LL)(a[i].x+0.5);
    LL x,y;
    x=y=0;
    for(i=10000;i>=1;--i)
        tmp[i]+=tmp[i+1];
    for(i=1;i<10000;++i)
        x+=ans[i]*tmp[i+1];
    y=n*(n-1);
    y>>=1;
    y=y*y*y;
    LL g=gcd(x,y);
    x/=g;
    y/=g;
    printf("%.15lf\n",1.0*x/y);
    return 0;
}