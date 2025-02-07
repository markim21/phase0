package server

import "testing"

func TestCopy(t *testing.T) {

	// test successful copy response
	t.Run("test a successful copy", func(t *testing.T) {
		var response1, response2, response Response
		var request1, request2, request Request 
		request1 = {"key": "1", "val": "abc", "op": "CREATE"}
		doCreate(request1, response1)
		request2 = {"key": "2", "val": "xyz", "op": "CREATE"}
		doCreate(request2, response2)

		var request Request = {"src_key": "1", "op": "COPY", "dst_key":"2"}
		var response Response
		response.Status = FAIL

		doCopyVal(request, response)
		
		// test status
		if response.Status != OK {
			t.Errorf("got %q want %q", response.Status, OK)
		}
	})

	// test successful copy value
	//t.Run("test dst_key is equal to src_key", func(t *testing.T) {	})

	// fail - src_key doesn't exist
	// fail - dst_key doesn't exist
	// fail - dst_key == src_key (?fail?)
}

func TestLogin(t *testing.T) { // assume every user u has uid string

	// test begin session
	// user u at client begins a session by issuing 
	// {"uid": "fbs", "op": "LOGIN"} where fbs is value for uid

	// test multiple sessions (don't exist, can't exist)

	// during a session 
	// every request submitted should echo uid of user u (inserted automatically by client)
	// every response sent by server should echo uid in Request.
}

func TestLogout(t *testing.T) {

	// test end current session
	// user u ends current session by issuing 
	// {"op": "LOGOUT"}

}
/*
func setupKeys(t testing.TB, request1 *Request, request2 *Request) {
	t.Helper()
	var response1, response2 Response
	// src_key
	request1 = {"key": "1", "val": "abc", "op": "CREATE"}
	doCreate(request1, response1)

	// dst_key
	request2 = {"key": "2", "val": "xyz", "op": "CREATE"}
	doCreate(request2, response2)
}

func clearKeys(t testing.TB){
	t.Helper()
	var request1, request2 Request
	var response1, response2 Response
	request1 = {"key" : "1", "op": "DELETE"}
	doDelete(request1, response1)
	request2 = {"key" : "2", "op": "DELETE"}
	doDelete(request2, response2)
}
*/
// go mod init ____ in each new folder before running commands like go test
