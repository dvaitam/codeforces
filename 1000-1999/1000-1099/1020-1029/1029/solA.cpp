#include<bits/stdc++.h>
using namespace std;
#define ll long long int
const int N=2e6+7;
vector <pair<ll,ll>> v,u;

int main()
{
    ios_base ::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL);
    ll m=0,n,i=0,q=5,x=0,k=0,y=1,z=1,j=1;
    string s,p="";
    cin>>n>>m>>s;
    k=n-2,i=n-2,j=n-1;
    for(;i>=0;)
    {
        if(s[j]==s[i])
        {
            if(i==0)
            break;
            j--,i--;
        }
        else
        {
            k--;
            i=k;
            j=n-1;
        }
    }
    if(i>=0)
    x=n-j;
    for(i=x;i<n;i++)
    p+=s[i];
    cout<<s;
    for(i=0;i<m-1;i++)
    cout<<p;
    return 0;
}