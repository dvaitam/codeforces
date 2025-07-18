#include<bits/stdc++.h>
using namespace std;
int main()
{
    int n,x=0,b,cnt=0;
    vector<int>v;
    cin>>n;
    int a[n];
    for(int i=0;i<n;i++)
    {
        cin>>a[i];
        x+=a[i];
    }
    b=a[0];
    if(b>x/2)
    {
        cout<<1<<endl<<1;
        return 0;
    }
    else
    {
        for(int i=1;i<n;i++)
        {
            if(a[i]*2<=a[0])
                {
                    b+=a[i];
                    cnt++;
                    v.push_back(i+1);
                }
                if(b>x/2)
            { cout<<cnt+1<<endl;
            cout<<1<<" ";
            for(int j=0;j<v.size();j++)
                cout<<v[j]<<" ";
                return 0;
            }


        }

    }
    cout<<0;
}