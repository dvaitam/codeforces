#include <cstdlib>
#include <cctype>
#include <cstring>
#include <cstdio>
#include <cmath>
#include <algorithm>
#include <vector>
#include <string>
#include <iostream>
#include <sstream>
#include <map>
#include <set>
#include <queue>
#include <stack>
#include <fstream>
#include <numeric>
#include <iomanip>
#include <bitset>
#include <list>
#include <stdexcept>
#include <functional>
#include <utility>
#include <ctime>
using namespace std;
#define ll long long
#define FOR( i, a, b ) for(int i = a; i <= b; ++i )
const double pi=3.1415926535897932384626433832795;
double w,h,a;
int main()
{
    //freopen("out.txt","w",stdout);
    //freopen("in.txt","r",stdin);
    int i,j,k;
    while(scanf("%lf%lf%lf",&w,&h,&a)!=EOF)
    {
        if(w<h)swap(w,h);
        double r=w*w+h*h;
        a/=180;
        a*=pi;
        if(a>pi/2)
        {
            a=pi-a;
        }
        r=sqrt(r);
        double angle=asin(h/r);
        angle*=2;
        double angle2=a/2;
        double ans=0;
        if(a>=angle&&a<=pi-angle)
        {
            double l=h/sin(angle2);
            double hh=l/2*tan(angle2);
            ans=l*hh;
            printf("%.7f\n",ans);
        }
        else
        {
            double l=h*sin(angle2);
            ans=l*l/tan(angle2);

            double g=cos(a)+sin(a)+1.0;
            double gg=cos(a)-sin(a)+1.0;
            double x=(w-(w+h)/g*sin(a))/gg;
            if(g==0)
            x=0;
            if(gg==0)
            x=0;
            double y=(w+h)/g-x;
            if(g==0)
            y=0;
            double s1=x*sin(angle2)*x*cos(angle2)*2;
            double s2=y*sin(angle2)*y*cos(angle2)*2;
            double s3=y*cos(angle2)*2*x*cos(angle2)*2;
            printf("%.7f\n",s1+s2+s3);
        }
    }
    return 0;
}