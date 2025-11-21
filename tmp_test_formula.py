import random,sys,subprocess,os

def run_go(src,input_data):
    p = subprocess.run(['bash','-lc',src],input=input_data.encode(),stdout=subprocess.PIPE,stderr=subprocess.PIPE,cwd='/home/ubuntu/codeforces')
    return p.stdout.decode(), p.stderr.decode()
