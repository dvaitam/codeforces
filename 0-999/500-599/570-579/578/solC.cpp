#include<cstdio>
#include<cstring>
#include<algorithm>
#include<cmath>
#include<iostream>
using namespace std;
const int N = 200500;
struct Line
{
 int a,b;
}xu[N];
struct Info
{
 Line line;
 long double start;
}maxv[N],minv[N];
int maxc=0,minc=0;
int n;
int getint()
{
 int res=0;
 char ch=getchar();
 while((ch<'0' || ch>'9') && ch!='-')
  ch=getchar();
 bool fan=0;
 if(ch=='-')
 {
  fan=1;
  ch=getchar();
 }
 while('0'<=ch && ch<='9')
 {
  res=res*10+ch-'0';
  ch=getchar();
 }
 if(fan)
  res=-res;
 return res;
}
void GetData()
{
 int i;
 n=getint();
 n++;
 xu[1].a=0;
 xu[1].b=0;
 for(i=2;i<=n;i++)
 {
  xu[i].a=i-1;
  xu[i].b=xu[i-1].b+getint();
 }
}
void DoIt()
{
 int i;
 for(i=1;i<=n;i++)
 {
  while(maxc)
  {
   long double yuan=maxv[maxc].line.a;
   yuan*=maxv[maxc].start;
   yuan+=maxv[maxc].line.b;
   long double now=xu[i].a;
   now*=maxv[maxc].start;
   now+=xu[i].b;
   if(yuan>now)
    break;
   else
    maxc--;
  }
  maxc++;
  maxv[maxc].line=xu[i];
  if(maxc==1)
   maxv[maxc].start=-20000;
  else
  {
   maxv[maxc].start=maxv[maxc].line.b-maxv[maxc-1].line.b;
   maxv[maxc].start/=maxv[maxc-1].line.a-maxv[maxc].line.a;
  }
 }
 for(i=n;i>=1;i--)
 {
  while(minc)
  {
   long double yuan=minv[minc].line.a;
   yuan*=minv[minc].start;
   yuan+=minv[minc].line.b;
   long double now=xu[i].a;
   now*=minv[minc].start;
   now+=xu[i].b;
   if(yuan<now)
    break;
   else
    minc--;
  }
  minc++;
  minv[minc].line=xu[i];
  if(minc==1)
   minv[minc].start=-20000;
  else
  {
   minv[minc].start=minv[minc].line.b-minv[minc-1].line.b;
   minv[minc].start/=minv[minc-1].line.a-minv[minc].line.a;
  }
 }
 int pointa=1,pointb=1;
 long double ans=1e50;
 while(pointa<=maxc && pointb<=minc)
 {
  long double wei,now;
  if(maxv[pointa].start>minv[pointb].start)
   wei=maxv[pointa].start;
  else
   wei=minv[pointb].start;
  now=(wei*maxv[pointa].line.a+maxv[pointa].line.b);
  now-=(wei*minv[pointb].line.a+minv[pointb].line.b);
  ans=min(ans,now);
  if(pointa==maxc)
   pointb++;
  else if(pointb==minc)
   pointa++;
  else if(maxv[pointa+1].start>minv[pointb+1].start)
   pointb++;
  else
   pointa++;
 }
 double ansx=ans;
 printf("%.10f\n",ansx);
}
int main()
{
 GetData();
 DoIt();
 return 0;
}