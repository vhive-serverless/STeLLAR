import random

ALLOWEDCHARS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
UNIQUEPAYLOADBYTES = 1024 * 1024

def InitializeGlobalRandomPayload():

	for i in range (0,UNIQUEPAYLOADBYTES):
		uniquePayload.append(random.randint(0,len(ALLOWEDCHARS)-1))


	GlobalRandomPayload = ''.join(str(e) for e in uniquePayload)

	length= len(GlobalRandomPayload)
	for i in range(0,3):
		length *= 2
		GlobalRandomPayload =GeneratePayloadFromGlobalRandom(length)

def GeneratePayloadFromGlobalRandom(payloadLengthBytes):
	# INCOMPLETE: Not sure where the GlobalRandomPayload variable comes from

	while len(repeatedRandomPayload) < payloadLengthBytes:
		repeatedRandomPayload = 

def extractJSONTimestampChain(responsePayload):
	# assuming it's StreamingBody
	reply = json.load(responsePayload)
	# or reply = json.loads(payloadResponse.read().decode("utf-8"))
	return reply["body"]["TimeStampChain"]

def AppendTimestampToChain(timestampChain):
	return timestampChain.append(time.utcnow())

def StringArrayToArrayOfString(stri):
	stri1 = stri[1:]
	stri1 = stri1[:-1]
	return stri1.split()


