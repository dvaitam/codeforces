#include <bits/stdc++.h>

using namespace std;

int main()
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    long long m,n,k,l;
    cin>>n>>m>>k>>l;
    if(n<m||l>n||k>=n)
        cout<<-1<<endl;
    else
    {
        long long mi=l+k;
        if(mi%m==0)
        {
            if(m*(mi/m)<=n)
                cout<<mi/m<<endl;
            else
                cout<<-1<<endl;
        }
        else
        {
            if(m*(mi/m+1)<=n)
                cout<<mi/m+1<<endl;
            else
                cout<<-1<<endl;
            
        }
    }
    return 0;
}