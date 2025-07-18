#include<bits/stdc++.h>
using namespace std;
int main()
{
    long long n,m,a,b,c,d,i,j,k,l,p,q,r,s[1000];
    string x,y,z;
    cin>>n>>m;
    d=m;
    l=0;
    k=0;
    for(i=n;i>=1;i--)
    {
        if(((d/i)>0))
        {
            if(d/i>1)
            {
                l=l+d/i;
            }
            else
            {

            l++;
            }
            d=d%i;
            k=k+i;
            if(k==m)
            {
                break;
            }
            //cout<<d<<"***"<<" "<<k<<endl;
        }
       // cout<<d<<" "<<k<<endl;
    }
    cout<<l;
}