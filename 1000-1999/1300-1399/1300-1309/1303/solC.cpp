#include <bits/stdc++.h>
using namespace std;
#define sws ios_base::sync_with_stdio(false);cout.tie(NULL);cin.tie(NULL);
#define mp make_pair
#define pb push_back
#define rep(i, a, b) for (int i = a; i < b; i++)
#define dbg(msg,x) cout<<msg<<" "<<x<<endl;
#define output(x) for(auto c:x){cout<<c<<" ";}cout<<" ";
#define ll long long 
#define ff first
#define ss second
#define pq priority_queue
typedef vector<int> vi;
typedef vector<bool> vb;
typedef pair<int, int> pii;
typedef vector<pair<int,int> > vpp;


int charToInt(char a){
    return a-(int)'a';
}

char intToChar(int a){
    return (char)(a+'a');
}

int main(){
    sws
    int t;cin>>t;
    while(t--){
        string s;cin>>s;
        bool adj[26][26];
        memset(adj,false,sizeof(adj));
        int n = s.size();

        rep(i,0,n){
            if(i>0){
                char bf = s[i-1];
                adj[charToInt(bf)][charToInt(s[i])]=1;
                adj[charToInt(s[i])][charToInt(bf)]=1;
            }
            if(i<n-1){
                char bf = s[i+1];
                adj[charToInt(bf)][charToInt(s[i])]=1;
                adj[charToInt(s[i])][charToInt(bf)]=1;     
            }
        }   

        bool visited[26];
        int connections[26];
        memset(visited,false,sizeof(visited));
        memset(connections,0,sizeof(connections));

        string ans="";
        bool poss = true;

        rep(i,0,26){
            int counter=0;
            rep(j,0,26){
                counter+=adj[i][j];
            }

            if(counter>2){
                poss=false;
            }
            connections[i]=counter;
        }

        if(!poss){
            cout<<"NO\n";
            continue;
        }
        
        rep(i,0,26){
            if(!visited[i] && connections[i]!=2){
                stack<pii>s;
                s.push(mp(i,-1));

                while(!s.empty()){
                    int cur = s.top().ff, p  = s.top().ss;
                    s.pop();
                    visited[cur]=true;
                    ans+=intToChar(cur);

                    rep(c,0,26){
                        if(adj[cur][c]==1){
                            if(c==p)continue;
                            if(visited[c]){
                                poss=false;
                                break;
                            }else{
                                s.push(mp(c,cur));
                            }
                        }
                    }
                }
            }
        }

        if(!poss || ans.size()!=26){
            cout<<"NO\n";
            continue;
        }

        cout<<"YES\n"<<ans<<"\n";

    }

}