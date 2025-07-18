/*And I thought my jokes were bad*/
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <iostream>
#include <bits/stdc++.h>
using namespace std;
#define ll long long
#define dbg puts("Viva la vida");
#define CHECK(x) cout << (#x) << " is " << (x) << endl;
#define nl puts("");
typedef map<int,int> mii;
typedef map<string,int> msi;
typedef pair<int,int> pii;
typedef pair<int,pii > tii;
typedef vector <int> VI;
typedef vector<ll> VL;
typedef set<int> SI;
#define mp make_pair
#define pb push_back
#define IN(x) scanf("%d",&x);
#define INL(x) scanf("%lld",&x);
#define OUT(x) printf("%d",x);
#define OUTL(x) printf("%lld",x);
#define SP printf(" ");
#define X first
#define Y second
#define SZ(_a) (int)_a.size()
#define ALL(_a) _a.begin(),_a.end()
#define EPS 1e-9
#define PI acos(-1.0)
#define MAX 200005
#define MOD 1000000007
#define INF (1 << 31)
/*
Like memories in cold decay
Transmissions echoing away
Far from the world of you and I
Where oceans bleed into the sky
*/
char s[MAX];
int main()
{
//    ios::sync_with_stdio(false);
//    cin.tie(0);
 
 
    int i,j,k,l,m,n,t,x,y,a,b,cnt;
    cin>>n;
    scanf("%s",s);
    a=b=cnt=0;
    for(i=1;i<n;i+=2)
    {
        if(s[i]==s[i-1])
        {
            if(s[i-1]=='a')
                s[i]='b';
            else
                s[i]='a';
cnt++;
        }
    }
    cout<<cnt<<endl;
    printf("%s",s);
    return 0;
}