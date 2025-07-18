#include <bits/stdc++.h>

using namespace std;

int cnt[28], id, fst, snd;

char temp[28], res[2][14];

int main()

{

    //freopen("in.txt", "r", stdin);

    scanf("%s", temp);

    for (int i=0; i<26; ++i)

        if (temp[i]==temp[i+1]) {

            puts("Impossible");

            return 0;

        }

    for (int i=0; i<27; ++i) {

        ++cnt[temp[i]-'A'];

        if (cnt[temp[i]-'A']==2)

            id=temp[i]-'A';

    }

    fst=snd=-1;

    for (int i=0; i<27; ++i)

        if (temp[i]==id+'A') {

            if (fst==-1)

                fst=i;

            else

                snd=i;

        }

    res[0][12-((snd-fst-1)>>1)]=id+'A';

    bool flag=false;

    int ptr=12-((snd-fst-1)>>1);

    ++ptr;

    for (int i=fst+1; i<=snd-1; ++i) {

        if (ptr==13) {

            --ptr;

            flag=true;

        }

        if (!flag)

            res[flag][ptr++]=temp[i];

        if (flag)

            res[flag][ptr--]=temp[i];

    }

    flag=false;

    ptr=12-((snd-fst-1)>>1);

    --ptr;

    for (int i=fst-1; i>=0; --i) {

        if (ptr<0) {

            ++ptr;

            flag=true;

        }

        if (!flag)

            res[flag][ptr--]=temp[i];

        if (flag)

            res[flag][ptr++]=temp[i];

    }

    ptr=12-((snd-fst-1)>>1);

    if ((snd-fst-1)&1)

        --ptr;

    flag=true;

    for (int i=snd+1; i<27; ++i) {

        if (ptr<0) {

            ++ptr;

            flag=false;

        }

        if (flag)

            res[flag][ptr--]=temp[i];

        if (!flag)

            res[flag][ptr++]=temp[i];

    }

    puts(res[0]);

    puts(res[1]);

    return 0;

}