#include <bits/stdc++.h>
using namespace std;
#define faster ios_base:: sync_with_stdio(false),cin.tie(0),cout.tie(0)


int main()
{
    long long int a[101],l[101],r[101],b[101][2],bb=0,ll=0,rr=0;
    int n,i,j;

    cin>>n;
    for(i=0; i<n; i++)
        cin>>a[i];
    sort(a,a+n);

    long long int cnt=0;
    for(i=0; i<n; i++)
    {
        if(a[i]==a[i+1])
            cnt++;
        else
        {
            cnt++;
            b[bb][0]=a[i];
            b[bb][1]=cnt;
            bb++;
            cnt=0;
        }
    }

    int flag=0;
    for(i=bb-2; i>=0; i--)
    {
        if(b[i][1]%2==0)
        {
            for(j=0; j<b[i][1]/2; j++)
            {
                l[ll++]=b[i][0];
                r[rr++]=b[i][0];
            }
        }
        else
        {
            for(j=0; j<b[i][1]/2; j++)
            {
                l[ll++]=b[i][0];
                r[rr++]=b[i][0];
            }
            if(flag==0)
            {
                l[ll++]=b[i][0];
                flag=1;
            }
            else if(flag==1)
            {
                r[rr++]=b[i][0];
                flag=0;
            }
        }

    }

    for(i=ll-1; i>=0; i--)
        cout<<l[i]<<' ';
    for(j=0; j<b[bb-1][1]; j++)
        cout<<b[bb-1][0]<<' ';
    for(i=0; i<rr; i++)
        cout<<r[i]<<' ';

    return 0;
}