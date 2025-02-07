package server

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/google/uuid"
	"os"
	
	"crypto_utils"
	. "types"
)

var privateKey *rsa.PrivateKey 
var publicKey  *rsa.PublicKey

var name string 
var currentUid string
var kvstore map[string]interface{}
var Requests chan NetworkData
var Responses chan NetworkData

func init() {
	privateKey = crypto_utils.NewPrivateKey()
	publicKey = &privateKey.PublicKey
	publicKeyBytes := crypto_utils.PublicKeyToBytes(publicKey)
	if err := os.WriteFile("SERVER_PUBLICKEY", publicKeyBytes, 0666); err != nil {
		panic(err)
	}

	name = uuid.NewString()
	currentUid = ""
	kvstore = make(map[string]interface{})
	Requests = make(chan NetworkData)
	Responses = make(chan NetworkData)

	go receiveThenSend()
}

func receiveThenSend() {
	defer close(Responses)

	for request := range Requests {
		Responses <- process(request)
	}
}

// Input: a byte array representing a request from a client.
// Deserializes the byte array into a request and performs
// the corresponding operation. Returns the serialized
// response. This method is invoked by the network.
func process(requestData NetworkData) NetworkData {
	var request Request
	json.Unmarshal(requestData.Payload, &request)
	var response Response
	doOp(&request, &response)
	responseBytes, _ := json.Marshal(response)
	return NetworkData{Payload: responseBytes, Name: name}
}

// Input: request from a client. Returns a response.
// Parses request and handles a switch statement to
// return the corresponding response to the request's
// operation.
func doOp(request *Request, response *Response)  {
	response.Status = FAIL
	response.Uid = request.Uid
	switch request.Op {
	case NOOP:
		// NOTHING
	case CREATE:
		doCreate(request, response)
	case DELETE:
		doDelete(request, response)
	case READ:
		doReadVal(request, response)
	case WRITE:
		doWriteVal(request, response)
	case COPY:
		doCopyVal(request, response)
	case LOGIN:
		doLogin(request, response)
	case LOGOUT:
		doLogout(request, response)
	default:
		// struct already default initialized to
		// FAIL status
	}
}

/** begin operation methods **/
// Input: key k, value v, metaval m. Returns a response.
// Sets the value and metaval for key k in the
// key-value store to value v and metavalue m.
func doCreate(request *Request, response *Response) {
	if _, ok := kvstore[request.Key]; !ok {
		kvstore[request.Key] = request.Val
		response.Status = OK
	}
}

// Input: key k. Returns a response. Deletes key from
// key-value store. If key does not exist then take no
// action.
func doDelete(request *Request, response *Response) {
	if _, ok := kvstore[request.Key]; ok {
		delete(kvstore, request.Key)
		response.Status = OK
	}
}

// Input: key k. Returns a response with the value
// associated with key. If key does not exist
// then status is FAIL.
func doReadVal(request *Request, response *Response) {
	if v, ok := kvstore[request.Key]; ok {
		response.Val = v
		response.Status = OK
	}
}

// Input: key k and value v. Returns a response.
// Change value in the key-value store associated
// with key k to value v. If key does not exist
// then status is FAIL.
func doWriteVal(request *Request, response *Response) {
	if _, ok := kvstore[request.Key]; ok {
		kvstore[request.Key] = request.Val
		response.Status = OK
	}
}

// Input: key to copy value in src_key Returns a response.
// Value associated with src_key should be copied to dst_key.
// If either the src_key or dst_key does not exist, then status is FAIL.
func doCopyVal(request *Request, response *Response) {
	v, src_ok := kvstore[request.SrcKey]
	_, dst_ok := kvstore[request.DstKey]
	if src_ok && dst_ok {
		kvstore[request.DstKey] = v
		response.Status = OK
	}
}

// Input: uid to begin session. Returns a response.
// Begins a session with the uid in the request.
// If there is already a session active, do nothing. 
func doLogin(request *Request, response *Response) {
	if currentUid == "" {
		currentUid = request.Uid
		response.Status = OK
	} else {
		// Do not share the currently active uid if a second login is attempted
		response.Uid = ""
	}
}

// Input: None. Returns a response.
// Ends the current session.
func doLogout(request *Request, response *Response) {
	currentUid = ""
	response.Status = OK
	response.Uid = currentUid
}