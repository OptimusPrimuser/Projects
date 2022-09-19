import requests
import json
from flask import Flask

app = Flask(__name__)


@app.route('/node/groups')
def getNodes():
    jsonOUT=json.loads(requests.get("http://localhost:8080/v1/api/pod.json").text)
    nodeDict={}
    for item in jsonOUT['items']:
        node=item['spec']['nodeName']
        app_name=item['metadata']['name'][0]
        nodeDict[node]=nodeDict.get(node,set())
        nodeDict[node].add(app_name)
    for key in nodeDict:
        nodeDict[key]=list(nodeDict[key])
    tempSet={}
    for appVal in nodeDict.values():
        for app_n in appVal:
            tempSet[app_n]=tempSet.get(app_n,tuple())
            if len(tempSet[app_n])<len(appVal):
                tempSet[app_n]=tuple(appVal)

    nodeGroups={}
    for key in nodeDict:
        ngKey=tempSet[nodeDict[key][0]]
        nodeDict[key]=ngKey
        nodeGroups[ngKey]=nodeGroups.get(ngKey,list())
        nodeGroups[ngKey].append(key)

    resOut={}
    ngKeys=list(nodeGroups.keys())
    resOut["groups"]=[]
    for i in range(len(ngKeys)):
        data={}
        data["name"]="G{}".format(i+1)
        data["nodes"]=nodeGroups[ngKeys[i]]
        data["apps"]=ngKeys[i]
        resOut["groups"].append(data)
    resOut["message"]="success"
    resOut["groupCount"]=len(ngKeys)
  
    return resOut


app.run(debug=True, port=9000, host='0.0.0.0')
