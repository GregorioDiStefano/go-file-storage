from flask import Flask, request
import requests
import simplejson as json
import sys
import subprocess
from urlparse import urlparse
app = Flask(__name__)


@app.route("/update", methods=['POST'])
def update():
    try:
        jsonData = json.loads(request.data)
        callback_url = jsonData["callback_url"]
    except:
        sys.stderr.write("Error parsing JSON")
        return "Fail"

    if urlparse(callback_url).hostname != "registry.hub.docker.com":
        sys.stderr.write("Invalid callback_url: ", callback_url, "\n")
        return "Fail"

    reponse = requests.post(callback_url, '{"state": "success"}')

    if reponse.status_code == 200:
        subprocess.Popen([sys.argv[1]], shell=True)
    else:
        sys.stderr.write("Failed to validate hook.\n")
        return "Fail"

    return "OK"

if __name__ == "__main__":
    app.debug = True
    print "Will run: ", sys.argv[1], " on successfull webhook call."
    app.run(host="0.0.0.0", port=8080)
