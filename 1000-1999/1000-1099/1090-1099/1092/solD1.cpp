#include <bits/stdc++.h>
using namespace std;
int n, a[200005];
int main()
{
    while(~scanf("%d",&n)) {
        stack <int> s;
        for(int i=0;i<n;i++) {
            int m; scanf("%d",&m);
            m%=2;
            if(!s.empty()) {
                if(s.top()==m) s.pop();
                else s.push(m);
            } else s.push(m);
        }
        if(s.size()>1) printf("NO\n");
        else printf("YES\n");
    }

}