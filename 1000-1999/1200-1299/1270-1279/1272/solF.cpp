#include<iostream>
#include<vector>
#include<algorithm>

using namespace std;

#define int int16_t

const int bm=10;

const int INF=1e4;

struct query
{
    int l;
    int r;
    int q;
};

int dp[400+5][200+5][610+5];

query p[400+5][200+5][610+5];

string get_ans(query answ)
{
    string ns="";
    int bal_prev=answ.q-bm;
    while (answ.l!=0||answ.r!=0)
    {
        if (p[answ.l][answ.r][answ.q].q<answ.q) ns+="(";
        else ns+=")";
        answ=p[answ.l][answ.r][answ.q];
    }
    for (int i=0; i<-bal_prev; i++)
    {
        ns+="(";
    }
    reverse(ns.begin(), ns.end());
    int balance=0;
    for (int i=0; i<(int)ns.size(); i++)
    {
        if (ns[i]=='(') balance++;
        else balance--;
    }
    for (int i=0; i<balance; i++)
    {
        ns+=")";
    }
    return ns;
}

int greedy(string a, string b)
{
    int bal=0;
    int l=0;
    int r=0;
    int k=INF;
    while (l!=a.size()||r!=b.size())
    {
        k=min(k,bal);
        if (l==a.size())
        {
            if (b[r]=='(') bal++;
            else bal--;
            r++;
            continue;
        }
        if (r==b.size())
        {
            if (a[l]=='(') bal++;
            else bal--;
            l++;
            continue;
        }
        if (a[l]=='(')
        {
            bal++;
            l++;
            continue;
        }
        if (b[r]=='(')
        {
            bal++;
            r++;
            continue;
        }
        if (a[l]==b[r]&&a[l]==')')
        {
            bal--;
            l++;
            r++;
            continue;
        }
        if (a[l]==')')
        {
            bal--;
            l++;
            continue;
        }
        if (b[r]==')')
        {
            bal--;
            r++;
            continue;
        }
    }
    k=min(k,bal);
    if (k>=0) return 0;
    else return -k;
}

int32_t main()
{
    string a,b;
    cin>>a>>b;
    int n,m;
    int k=greedy(a,b);
    string t="";
    for (int i=0; i<k; i++)
    {
        t+="(";
    }
    a=t+a;
    n=a.size();
    m=b.size();
    for (int sum=0; sum<=n+m; sum++)
    {
        for (int l=0; l<=n; l++)
        {
            int r=sum-l;
            if (r<0||r>m) continue;
            for (int bal=bm; bal<=n+m+bm; bal++)
            {
                dp[l][r][bal]=INF;
            }
        }
    }
    dp[0][0][0+bm]=0;
    p[0][0][0]={-1,-1,-1};
    a+="!";
    b+="!";
    for (int sum=0; sum<n+m; sum++)
    {
        for (int l=0; l<=n; l++)
        {
            int r=sum-l;
            if (r<0||r>m) continue;
            //cout<<l<<" "<<r<<endl;
            for (int bal=bm; bal<=n+m+bm; bal++)
            {
                if (dp[l][r][bal]==INF) continue;
                //cout<<l<<" "<<r<<" "<<bal-bm<<endl;
                if (a[l]=='(')
                {
                    if (dp[l][r][bal]+1<dp[l+1][r][bal+1])
                    {
                        p[l+1][r][bal+1]={l,r,bal};
                        dp[l+1][r][bal+1]=dp[l][r][bal]+1;
                    }
                }
                if (a[l]==')')
                {
                    if (dp[l][r][bal]+1<dp[l+1][r][bal-1])
                    {
                        p[l+1][r][bal-1]={l,r,bal};
                        dp[l+1][r][bal-1]=dp[l][r][bal]+1;
                    }
                }
                if (b[r]=='(')
                {
                    if (dp[l][r][bal]+1<dp[l][r+1][bal+1])
                    {
                        p[l][r+1][bal+1]={l,r,bal};
                        dp[l][r+1][bal+1]=dp[l][r][bal]+1;
                    }
                }
                if (b[r]==')')
                {
                    if (dp[l][r][bal]+1<dp[l][r+1][bal-1])
                    {
                        p[l][r+1][bal-1]={l,r,bal};
                        dp[l][r+1][bal-1]=dp[l][r][bal]+1;
                    }
                }
                if (a[l]==b[r]&&a[l]=='(')
                {
                    if (dp[l][r][bal]+1<dp[l+1][r+1][bal+1])
                    {
                        p[l+1][r+1][bal+1]={l,r,bal};
                        dp[l+1][r+1][bal+1]=dp[l][r][bal]+1;
                    }
                }
                if (a[l]==b[r]&&a[l]==')')
                {
                    if (dp[l][r][bal]+1<dp[l+1][r+1][bal-1])
                    {
                        p[l+1][r+1][bal-1]={l,r,bal};
                        dp[l+1][r+1][bal-1]=dp[l][r][bal]+1;
                    }
                }
            }
        }
    }
    int ans=INF;
    string ns="";
    for (int bal=bm; bal<=n+m+bm; bal++)
    {
        if (dp[n][m][bal]==INF) continue;
        query answ={n,m,bal};
        string y=get_ans(answ);
        if ((int)y.size()<ans)
        {
            ans=y.size();
            ns=y;
            //cout<<bal-bm<<endl;
        }
    }
    cout<<ns;
}