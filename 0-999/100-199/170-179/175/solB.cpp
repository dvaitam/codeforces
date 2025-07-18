#include <vector>
#include <list>
#include <map>
#include <set>
#include <deque>
#include <stack>
#include <bitset>
#include <algorithm>
#include <functional>
#include <numeric>
#include <utility>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <cstdio>
#include <cmath>
#include <cstdlib>
#include <ctime>
#include <cstring>
using namespace std;
typedef long long ll;
typedef unsigned long long ull;
char name[1200][20];
int score[1200],n;
int idx[1200];
int cmp( const void *arg1, const void *arg2 ){
   return strcmp(name[*(int *)arg1],name[*(int *)arg2]);
}
int cmp2( const void *arg1, const void *arg2 ){
   return score[*(int *)arg1]-score[*(int *)arg2];
}
int cmp3(int arg1,int arg2){
    return score[arg1]-score[arg2];
}
int *upper(int i){
    int s=i;
    for(;i<n&&score[idx[i]]==score[idx[s]];i++);
    return idx+i;
}
double bt(int s){
    int *ptr=upper(s);
    int p=ptr-idx;
    return (double)(n-1-p+1)/(double)n;
}
double nwt(int s){
    int *ptr=upper(s);
    ptr--;
    int p=ptr-idx;
    return (double)(p+1)/(double)n;
}
int main(){
#ifdef DEBUG
	freopen("2.txt","r",stdin);
#endif
    int i,j;
    scanf("%d",&n);
    for(i=0;i<n;i++) idx[i]=i;
    for(i=0;i<n;i++)scanf(" %s %d ",name[i],&score[i]);
    qsort(idx,n,sizeof(idx[0]),cmp);

    for(i=1,j=0;i<n;i++)
        if(strcmp(name[idx[i-1]],name[idx[i]])==0){
            if(score[idx[i]]>score[idx[j]]) idx[j]=idx[i];
        }
        else{
            j++;
            idx[j]=idx[i];
        }
    n=j+1;
    qsort(idx,n,sizeof(idx[0]),cmp2);
    printf("%d\n",n);
    for(i=0;i<n;i++){
        double nwtr=nwt(i);
        double btr=bt(i);
        if(nwtr>=0.99)printf("%s pro\n",name[idx[i]]);
        else if(nwtr>=0.9&&btr>0.01)printf("%s hardcore\n",name[idx[i]]);
        else if(nwtr>=0.8&&btr>0.1)printf("%s average\n",name[idx[i]]);
        else if(nwtr>=0.5&&btr>0.2)printf("%s random\n",name[idx[i]]);
        else if(btr>0.5)printf("%s noob\n",name[idx[i]]);
    }
	return 0;
}