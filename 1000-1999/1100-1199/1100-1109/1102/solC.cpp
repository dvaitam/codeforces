#include<bits/stdc++.h>

using namespace std;
typedef long long LL;

int main() {
    /* Enter your code here. Read input from STDIN. Print output to STDOUT */   
        LL n,a,b;
        cin>>n>>a>>b;
        LL arr[n],vis[n];
        for(int i=0;i<n;i++)
        {
                cin>>arr[i];
                vis[i]=0;
        }
        int cnt1=0,cnt2=0;
        if(a>b)
        {
                cout<<n<<endl;
                exit(0);
        }
        int ans=0;
        //sort(arr,arr+n)
        int turn=0;
        for(int i=0;i<n;i++)
        {
                if(arr[i]-a<=0)
                        ans++;
        }
        if(ans%2)
        {
                ans=ans/2+1;
        }
        else
                ans=ans/2;
        cout<<ans<<endl;
    return 0;
}