#include<iostream>
#include<cstdlib>
#include<cstdio>//old c library
#include<cmath>//ceil,floor,round,M_PI,,trig,pow(ambigous)
#include<string>//sting c++ library
#include<cstring>//string c library
#include<iomanip>//set precisions
#include <algorithm>//using swap,sort using some function.
#include<fstream>//used to manipulate files.
#include<stack>//stl for stack
#include<queue>//stl fro queue
#include<map>//stl for map key,info asscociated
//#include <boost/multiprecision/cpp_int.hpp>//algerbric operations on string
//namespace mp=boost::multiprecision;

using namespace std;
int n,x,y,p[2005],pp[2005],s[2005],ps[2005],c[2005][2],mc=0,mm=0,flag=0;

int main()
{
    scanf("%d",&n);
    for(x=1;x<=n;x++)
    {
    scanf("%d",&p[x]);
    pp[p[x]]=x;
    }   
    for(x=1;x<=n;x++)
    {
    scanf("%d",&s[x]);
    ps[s[x]]=x;
    }
    for(x=1;x<=n;x++)
    {
    for(y=1;y<=n;y++)
    {
    if(p[x]==s[y])
    {
        p[x]=y;
        break;
        }}}
    for(x=1;x<=n;x++)
    {
        mc=mc+abs(pp[x]-ps[x]);
        }
    while(1)
    {
        flag=0;
        for(x=1;x<=n;x++)
        {
            if(p[x]>=x+1)
            {
            for(y=p[x];y>x;y--)
            {
                if(p[y]<=x)
                {
                    swap(p[x],p[y]);
                    c[mm][0]=x;
                    c[mm][1]=y;
                    mm++;
                    flag=1;
                    break;
                    }
                }}
            }
        if(!flag) break;
        }
    printf("%d\n%d\n",mc/2,mm);
    for(x=0;x<mm;x++)
    printf("%d %d\n",c[x][0],c[x][1]);
    return 0;
}