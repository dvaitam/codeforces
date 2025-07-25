#include <algorithm>
#include <stdio.h>
#include <string.h>
#include <cstdio>
#include <cmath>
#include <vector>
#define maxn 70
#define maxlen 2000007
using namespace std;
 
typedef vector<double> VD;
typedef vector<VD> VVD;
typedef vector<int> VI;
 
int cost[maxn][maxn];
int Lmate[maxn],Rmate[maxn];
int u[maxn],v[maxn];
int dist[maxn];
int dad[maxn];
int seen[maxn];
int n,m;
char str1[maxlen],str2[maxlen];
int get_id(char ch){
    if(ch>='a' & ch<='z') return ch-'a';
    else return ch-'A'+26;
}
char get_ch(int aa){
    if(aa>=26) return aa+'A'-26;
    return 'a'+aa;
}
int MinCostMatching() {
    for (int i = 0; i < n; i++) {
        u[i] = cost[i][0];
        for (int j = 1; j < n; j++) u[i] = min(u[i], cost[i][j]);
    }
    for (int j = 0; j < n; j++) {
        v[j] = cost[0][j] - u[0];
        for (int i = 1; i < n; i++) v[j] = min(v[j], cost[i][j] - u[i]);
    }
    memset(Rmate,-1,sizeof Rmate);
    memset(Lmate,-1,sizeof Lmate);
    int mated = 0;
    for (int i = 0; i < n; i++) {
        for (int j = 0; j < n; j++) {
            if (Rmate[j] != -1) continue;
            if (cost[i][j] - u[i] - v[j]==0) {
                Lmate[i] = j;
                Rmate[j] = i;
                mated++;
                break;
            }
        }
    }
        while (mated < n) {
 
        int s = 0;
        while (Lmate[s] != -1) s++;
        memset(dad,-1,sizeof dad);
        memset(seen,0,sizeof seen);
        for (int k = 0; k < n; k++)  dist[k] = cost[s][k] - u[s] - v[k];
        int j = 0;
        while (true) {
 
            j = -1;
            for (int k = 0; k < n; k++) {
                if (seen[k]) continue;
                if (j == -1 || dist[k] < dist[j]) j = k;
            }
            seen[j] = 1;
            // termination condition
            if (Rmate[j] == -1) break;
            // relax neighbors
            const int i = Rmate[j];
            for (int k = 0; k < n; k++) {
                if (seen[k]) continue;
                double new_dist = dist[j] + cost[i][k] - u[i] - v[k];
                if (dist[k] > new_dist) {
                    dist[k] = new_dist;
                    dad[k] = j;
                }
            }
        }
 
        // update dual variables
        for (int k = 0; k < n; k++) {
            if (k == j || !seen[k]) continue;
            const int i = Rmate[k];
            v[k] += dist[k] - dist[j];
            u[i] -= dist[k] - dist[j];
        }
        u[s] += dist[j];
 
        // augment along path
        while (dad[j] >= 0) {
            const int d = dad[j];
            Rmate[j] = Rmate[d];
            Lmate[Rmate[j]] = j;
            j = d;
        }
        Rmate[j] = s;
        Lmate[s] = j;
 
        mated++;
    }
 
    int value = 0;
 
    for (int i = 0; i < n; i++)    value += cost[i][Lmate[i]];
    printf("%d\n",-value);
    for (int i = 0; i < n; i++)    printf("%c",get_ch(Lmate[i]));
    return value;
}
int main(){
    scanf("%d%d",&m,&n);
    scanf("%s%s",str1,str2);
    int i,j,k;
    for(i=0;i<m;i++){
        cost[get_id(str1[i])][get_id(str2[i])]--;
    }
    MinCostMatching();
 
 
    return 0;
}