#include<iostream>
#define MOD 998244353
using namespace std;
int n,i,rez,sum,sum1,Dp[2][2][205],A[100005],cr,j;

void f(int &a, int b)
{
    a=a+b;
    if(a>=MOD)
        a-=MOD;
}

int main()
{
    ios::sync_with_stdio(false);
    cin>>n;
    for(i=1; i<=n; i++)
        cin>>A[i];
    if(A[1]!=-1)
        Dp[1][1][A[1]]=1;
    else
        for(i=1; i<=200; i++)
            Dp[1][1][i]=1;
    cr=1;
    for(i=2; i<=n; i++)
    {
        cr=1-cr;
        if(A[i]!=-1)
        {
            for(j=1; j<=200; j++)
            {
                if(A[i]<j)
                    f(Dp[cr][0][A[i]],Dp[1-cr][0][j]);
                if(A[i]==j)
                {
                    f(Dp[cr][0][A[i]],Dp[1-cr][0][j]);
                    f(Dp[cr][0][A[i]],Dp[1-cr][1][j]);
                }
                if(A[i]>j)
                {
                    f(Dp[cr][1][A[i]],Dp[1-cr][0][j]);
                    f(Dp[cr][1][A[i]],Dp[1-cr][1][j]);
                }
                Dp[1-cr][1][j]=Dp[1-cr][0][j]=0;
            }
        }
        else
        {
            sum=0;
            for(j=1; j<=200; j++)
                f(sum,Dp[1-cr][0][j]);
            sum1=0;
            for(j=1; j<=200; j++)
            {
                f(Dp[cr][0][j],sum);
                f(Dp[cr][0][j],Dp[1-cr][1][j]);
                f(sum,MOD-Dp[1-cr][0][j]);
                f(Dp[cr][1][j],sum1);
                f(sum1,Dp[1-cr][1][j]);
                f(sum1,Dp[1-cr][0][j]);
                Dp[1-cr][0][j]=Dp[1-cr][1][j]=0;
            }
        }
    }
    for(i=1; i<=200; i++)
        f(rez,Dp[cr][0][i]);
    cout<<rez<<"\n";
    return 0;
}