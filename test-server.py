import subprocess
from flask import Flask, request
import time
import threading

lock = threading.Lock()
app = Flask(__name__)

@app.route("/")
def start():
    global lock
    if lock.locked():
        return "starting or stoping"
    print("start test")
    lock.acquire()

    env = request.args.get("env", "spb")
    size = request.args.get("size", "large")
    period = request.args.get("period", "25")

    cmd = "bash scripts/test-all.sh {} {} {}".format(env, size, period)
    print(cmd)

    res = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    while res.poll() is None:
        line = res.stdout.readline()
        line = line.strip()
        if line:
            print('Subprogram output: [{}]'.format(line))
    if res.returncode == 0:
        print('Subprogram success')
    else:
        print('Subprogram failed')

    time.sleep(10)

    lock.release()
    return "OK"