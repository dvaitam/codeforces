//program 314D

#include<iostream>
#include<cstdio>
#include<cstdlib>
#include<cstring>
#include<algorithm>

using namespace std;

struct Point
{
  long long X,Y;
};

bool operator <(Point A,Point B)
{
  return A.X<B.X;
}

int Get()
{
  char c;
  while(c=getchar(),(c<'0'||c>'9')&&(c!='-'));
  bool Flag=(c=='-');
  if(Flag)
    c=getchar();
  int X=0;
  while(c>='0'&&c<='9')
    {
      X=X*10+c-48;
      c=getchar();
    }
  return Flag?-X:X;
}

void Output(int X)
{
  if(X<0)
    {
      putchar('-');
      X=-X;
    }
  int Len=0,Data[10];
  while(X)
    {
      Data[Len++]=X%10;
      X/=10;
    }
  if(!Len)
    Data[Len++]=0;
  while(Len--)
    putchar(Data[Len]+48);
  putchar('\n');
}

const long long INF=1000000000000000LL;

Point A[100000];
long long LMax[100001],LMin[100001],RMax[100001],RMin[100001];

long long P(int L,int R)
{
  return A[R-1].X-A[L].X;
}

long long Q(int L,int R)
{
  return max(LMax[L],RMax[R])-min(LMin[L],RMin[R]);
}

int main()
{
  int N=Get();
  for(int i=0;i<N;i++)
    {
      int P=Get(),Q=Get();
      A[i].X=P+Q;
      A[i].Y=P-Q;
    }
  sort(A,A+N);
  LMax[0]=-INF;
  LMin[0]=INF;
  for(int i=0;i<N;i++)
    {
      LMax[i+1]=max(LMax[i],A[i].Y);
      LMin[i+1]=min(LMin[i],A[i].Y);
    }
  RMax[N]=-INF;
  RMin[N]=INF;
  for(int i=N-1;i>=0;i--)
    {
      RMax[i]=max(RMax[i+1],A[i].Y);
      RMin[i]=min(RMin[i+1],A[i].Y);
    }
  long long Ans=INF;
  for(int i=0,j=1;i<N;i++)
    {
      Ans=min(Ans,max(P(i,j),Q(i,j)));
      while(j<N&&P(i,j)<Q(i,j))
        {
          j++;
          Ans=min(Ans,max(P(i,j),Q(i,j)));
        }
    }
  double Real=(double)Ans/2;
  printf("%0.10lf\n",Real);
  return 0;
}