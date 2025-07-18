#include<bits/stdc++.h>
using namespace std;
#define maxx 1000000

int ara[maxx];

main()
{
    int n;
    cin >> n;

    if(n==3 || n==5)
    {
        cout << "NO" << endl;
    }
    else
    {
        if(n%2==0)
        {
            int p[]={-1,1,-1,-2,2,1};
            for(int i=0;i<n;i++)
            {
                ara[i]=p[i%6];
            }
        }
        else
        {
            int p[]={2,-2,-3,3,1,-1};
            int s[]={3,3,-2,-1,1,-1,2};
            for(int i=0;i<n;i++)
            {
                if(i<7)
                {
                    ara[i]=s[i];
                }
                else
                {
                    ara[i]=p[i%6];
                }
            }
        }
        cout << "YES" << endl;
        for(int i=0;i<n;i++)
        {
            cout << ara[i];
            if(i==(n-1))
            {
                cout << endl;
            }
            else
            {
                cout << " ";
            }
        }
    }

    return 0;
}