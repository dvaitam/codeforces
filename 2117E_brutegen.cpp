#include <bits/stdc++.h>
using namespace std;
int brute(vector<int>a, vector<int>b){
    int n=a.size();
    int best=0;
    for(int del=-1; del<n; ++del){
        vector<int> aa,bb;
        for(int i=0;i<n;i++) if(i!=del) aa.push_back(a[i]);
        for(int i=0;i<n;i++) if(i!=del) bb.push_back(b[i]);
        int m=aa.size();
        for(int mask=0; mask<(1<<(m-1)); ++mask){
            bool good=true;
            for(int i=0;i<m-1;i++){
                if((mask>>i)&1){
                    if(aa[i]!=bb[i+1]){good=false;break;}
                }else{
                    if(bb[i]!=aa[i+1]){good=false;break;}
                }
            }
            if(!good) continue;
            int cnt=0;
            for(int i=0;i<m;i++) cnt+=aa[i]==bb[i];
            best=max(best,cnt);
        }
    }
    return best;
}
int solve(vector<int>a, vector<int>b){
    map<int,int> cnt;
    for(int x:a) cnt[x]++;
    for(int x:b) cnt[x]++;
    int ans=0;
    for(auto [k,v]:cnt) ans=max(ans,v);
    return ans;
}
int main(){
    mt19937 rng(123);
    for(int n=2;n<=10;n++){
        vector<int> values(n);
        iota(values.begin(), values.end(), 1);
        vector<int> a(n), b(n);
        int LIMIT=10000;
        for(int iter=0; iter<LIMIT; ++iter){
            shuffle(values.begin(), values.end(), rng);
            for(int i=0;i<n;i++) a[i]=values[i];
            shuffle(values.begin(), values.end(), rng);
            for(int i=0;i<n;i++) b[i]=values[i];
            int br=brute(a,b);
            int sol=solve(a,b);
            if(br!=sol){
                cerr<<"Mismatch n="<<n<<" br="<<br<<" sol="<<sol<<"\n";
                cerr<<"a:";
                for(int x:a) cerr<<x<<" ";
                cerr<<"\nb:";
                for(int x:b) cerr<<x<<" ";
                cerr<<"\n";
                return 0;
            }
        }
    }
    cerr<<"All good!\n";
}
