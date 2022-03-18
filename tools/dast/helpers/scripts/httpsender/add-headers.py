# The sendingRequest and responseReceived functions will be called for all
# requests/responses sent/received by ZAP, including automated tools (e.g.
# active scanner, fuzzer, ...)

# Note that new HttpSender scripts will initially be disabled
# Right click the script in the Scripts tree and select "enable"

# 'initiator' is the component the initiated the request:
#      1   PROXY_INITIATOR
#      2   ACTIVE_SCANNER_INITIATOR
#      3   SPIDER_INITIATOR
#      4   FUZZER_INITIATOR
#      5   AUTHENTICATION_INITIATOR
#      6   MANUAL_REQUEST_INITIATOR
#      7   CHECK_FOR_UPDATES_INITIATOR
#      8   BEAN_SHELL_INITIATOR
#      9   ACCESS_CONTROL_SCANNER_INITIATOR
#     10   AJAX_SPIDER_INITIATOR
# For the latest list of values see the HttpSender class:
# https://github.com/zaproxy/zaproxy/blob/main/zap/src/main/java/org/parosproxy/paros/network/HttpSender.java
# 'helper' just has one method at the moment: helper.getHttpSender() which
# returns the HttpSender instance used to send the request.
#
# New requests can be made like this:
# msg2 = msg.cloneAll() # msg2 can then be safely changed without affecting msg
# helper.getHttpSender().sendAndReceive(msg2, false)
# print('msg2 response code =' + msg2.getResponseHeader().getStatusCode())
import json

def openFile(fileName):
    # Opening JSON file
    f = open(fileName)
    data = json.load(f)
    return data

def sendingRequest(msg, initiator, helper):
    # Debugging can be done using print like this
    # Don't add headers in check headers initiator
    if initiator != 7:
        headers = dict(openFile("/zap/helpers/headers.json"))
        for x in list(headers):
            msg.getRequestHeader().setHeader(x, headers[x]);
        print("Headers added to request")


def responseReceived(msg, initiator, helper):
    print("response received")
