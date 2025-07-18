#include <bits/stdc++.h>

using namespace std;



int main()
{
    long long a,b,c,i,j,k,l,m,n,d,h;
        cin>>n>>m;
        if(m>n)
        {
            c=n-m+n+1;
            if(c<0)
            {
                cout<<"0"<<endl;
                return 0;
            }
            b=c/2;
            cout<<b<<endl;
            return 0;

        }
        else
        {
            c=m-1;
            b=c/2;
            cout<<b<<endl;
            return 0;
        }

    return 0;
}