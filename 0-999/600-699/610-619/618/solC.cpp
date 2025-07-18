#include<bits/stdc++.h>
#define ll long long
#define pii pair<int,int>
#define piii pair<pair<int,int>,int>
#define ff first
#define ss second
#define mp make_pair
#define pb push_back
#define SIZE 100006
#define MOD 1000000007
using namespace std;

inline ll getnum()
{
    char c = getchar();
    ll num,sign=1;
    for(;c<'0'||c>'9';c=getchar())if(c=='-')sign=-1;
    for(num=0;c>='0'&&c<='9';)
    {
        c-='0';
        num = num*10+c;
        c=getchar();
    }
    return num*sign;
}

piii P[6];

int main()
{
    ll n=getnum(),ans=0,ax,ay,bx,by,cx,cy;

    for(int i=1;i<=3;i++)
    {
        int x=getnum(),y=getnum();
        P[i]=mp(mp(x,y),i);
    }

    sort(P+1,P+4);

    for(int i=4;i<=n;i++)
    {
        ll x=getnum(),y=getnum();

        ax=P[1].ff.ff;
        ay=P[1].ff.ss;
        bx=P[2].ff.ff;
        by=P[2].ff.ss;
        cx=P[3].ff.ff;
        cy=P[3].ff.ss;

        if((bx - ax)*(cy - ay) - (by - ay)*(cx - ax)==0)
        {
            if((bx - ax)*(y - ay) - (by - ay)*(x - ax)==0)
            {
                P[4]=mp(mp(x,y),i);
                sort(P+1,P+5);
            }
            else
            {
                P[3]=mp(mp(x,y),i);
            }
            continue;
        }

        ll f1=(bx - ax)*(y - ay) - (by - ay)*(x - ax);
        ll f2=(cx - bx)*(y - by) - (cy - by)*(x - bx);
        ll f3=(ax - cx)*(y - cy) - (ay - cy)*(x - cx);

        if(f1<=0&&f2<=0&&f3<=0||f1>=0&&f2>=0&&f3>=0)
        {
            if(f1==0)P[2]=mp(mp(x,y),i);
            else P[3]=mp(mp(x,y),i);
        }
    }

    printf("%d %d %d",P[1].ss,P[2].ss,P[3].ss);
}