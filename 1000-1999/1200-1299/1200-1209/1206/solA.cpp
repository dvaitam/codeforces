#include <bits/stdc++.h>
using namespace std;
#define ll long long
int main()
{
    int n,flag=0,t,q;cin>>n;
    int a[n];
    for(int i=0;i<n;i++)
    cin>>a[i];
    int m;cin>>m;
    int b[m];
    for(int i=0;i<m;i++)
    cin>>b[i];
    sort(a,a+n);sort(b,b+m);
    for(int i=0;i<n;i++)
    {
        for(int j=0;j<m;j++)
        {
            int f=(a[i]+b[j]);
            if(binary_search(a,a+n,f)||binary_search(b,b+m,f))
            flag=0;
            else 
            flag=1;
            if(flag==1)
            { q=a[i];t=b[j];break;}
        }
        if(flag==1)
        break;
    }
    cout<<q<<" "<<t;
    return 0;
}