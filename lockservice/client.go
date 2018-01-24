package lockservice

import (
	"crypto/rand"
	"math/big"
	"net/rpc"
	"strconv"
)

// the lockservice Clerk lives in the client
// and maintains a little state.
type Clerk struct {
	servers       [2]string // primary port, backup port
	primaryActive bool
	UUID          string
}

func MakeClerk(primary string, backup string) *Clerk {
	ck := new(Clerk)
	ck.servers[0] = primary
	ck.servers[1] = backup
	// Your initialization code here.
	ck.primaryActive = true
	ck.UUID = strconv.Itoa(int(nrand()))
	return ck
}

// call() sends an RPC to the rpcname handler on server srv
// with arguments args, waits for the reply, and leaves the
// reply in reply. the reply argument should be the address
// of a reply structure.
//
// call() returns true if the server responded, and false
// if call() was not able to contact the server. in particular,
// reply's contents are valid if and only if call() returned true.
//
// you should assume that call() will time out and return an
// error after a while if it doesn't get a reply from the server.
//
// please use call() to send all RPCs, in client.go and server.go.
func call(srv string, rpcname string,
	args interface{}, reply interface{}) bool {
	c, errx := rpc.Dial("unix", srv)
	if errx != nil {
		return false
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	return false
}

// generate numbers that have a high probability of being unique
func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

// generate uuid for a request
func getReqUuid(ck *Clerk) string {
	return ck.UUID + "-" + strconv.Itoa(int(nrand()))
}

// ask the lock service for a lock.
// returns true if the lock service
// granted the lock, false otherwise.
func (ck *Clerk) Lock(lockname string) bool {

	// prepare the arguments.
	args := &LockArgs{}
	args.Lockname = lockname
	args.UUID = getReqUuid(ck)
	var reply LockReply

	var ok bool
	// send an RPC request, wait for the reply.
	if ck.primaryActive == true {
		ok = call(ck.servers[0], "LockServer.Lock", args, &reply)
		ck.primaryActive = ok
	}
	// if P fails, contact B
	if ok == false {
		ok = call(ck.servers[1], "LockServer.Lock", args, &reply)
	}

	if ok == false {
		return false
	}
	return reply.OK
}

// ask the lock service to unlock a lock.
// returns true if the lock was previously held,
// false otherwise.
func (ck *Clerk) Unlock(lockname string) bool {

	// prepare the arguments.
	args := &UnlockArgs{}
	args.Lockname = lockname
	args.UUID = getReqUuid(ck)
	var reply UnlockReply

	var ok bool
	// send an RPC request, wait for the reply.
	if ck.primaryActive == true {
		ok = call(ck.servers[0], "LockServer.Unlock", args, &reply)
		ck.primaryActive = ok
	}
	// if P fails, contact B
	if ok == false {
		ok = call(ck.servers[1], "LockServer.Unlock", args, &reply)
	}

	if ok == false {
		return false
	}
	return reply.OK
}
