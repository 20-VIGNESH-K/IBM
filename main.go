package main

import (
	"fmt"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

func QueueConnection(qmgrName string) (ibmmq.MQQueueManager, error) {
	// Set the IBM MQ queue manager and queue information

	// IBM MQ server connection details
	channelName := "DEV.ADMIN.SVRCONN"
	connName := "localhost(1414)"
	userName := "admin"
	password := "passw0rd"

	//QUEUE CONNECTION

	gocno := ibmmq.NewMQCNO()
	var csp *ibmmq.MQCSP = ibmmq.NewMQCSP()
	csp.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
	csp.UserId = userName
	csp.Password = password
	gocno.SecurityParms = csp
	cd := ibmmq.NewMQCD()
	cd.ChannelName = channelName
	cd.ConnectionName = connName
	gocno.ClientConn = cd
	gocno.Options = ibmmq.MQCNO_CLIENT_BINDING
	return ibmmq.Connx(qmgrName, gocno)
}

func OpenPutQueue(queueName string, qMgr ibmmq.MQQueueManager) (ibmmq.MQObject, error) {
	openOptions := ibmmq.MQOO_OUTPUT
	od := ibmmq.NewMQOD()
	od.ObjectName = queueName
	od.ObjectType = ibmmq.MQOT_Q
	return qMgr.Open(od, openOptions)
}

func PutMessage(message string, qObj ibmmq.MQObject) error {
	pmd := ibmmq.NewMQMD()
	pmo := ibmmq.NewMQPMO()
	pmd.Format = ibmmq.MQFMT_STRING
	pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT | ibmmq.MQPMO_NEW_MSG_ID
	return qObj.Put(pmd, pmo, []byte(message))

}

func OpenGetQueue(queueName string, qMgr ibmmq.MQQueueManager) (ibmmq.MQObject, error) {
	openOptions := ibmmq.MQOO_INPUT_SHARED
	od := ibmmq.NewMQOD()
	od.ObjectName = queueName
	od.ObjectType = ibmmq.MQOT_Q
	return qMgr.Open(od, openOptions)

}

func GetMessage(qObj_read ibmmq.MQObject, buffer []byte) (int, error) {
	pmdr := ibmmq.NewMQMD()
	pmdr.Format = ibmmq.MQFMT_STRING
	gmo := ibmmq.NewMQGMO()
	return qObj_read.Get(pmdr, gmo, buffer)

}

func handlePanic() {
	if a := recover(); a != nil {
		fmt.Println("Error occured in get message")
	}
}

func main() {

	//qmgrName := "QM1"
	queueName := "DEV.DEAD.LETTER.QUEUE"
	message := "Hello world!!!"

	qMgr, err := QueueConnection("QM1")
	if err != nil {
		fmt.Println("error : ", err)
		panic(err)
	}
	fmt.Println("Queue connected!!!")
	defer qMgr.Disc()

	//OPEN_PUT_QUEUE

	qObj, err := OpenPutQueue(queueName, qMgr)
	if err != nil {
		fmt.Errorf("Error : ", err)
		panic(err)
	}
	fmt.Println("Write Queue Opened")
	defer qObj.Close(0)

	// //OPEN_PUT_MESSAGE

	err = PutMessage(message, qObj)
	if err != nil {
		fmt.Errorf("Error : ", err)
		panic(err)
	}
	fmt.Println("Message sent successfully!")

	// OPEN_GET_QUEUE

	qObj_read, err := OpenGetQueue(queueName, qMgr)
	if err != nil {
		fmt.Errorf("Error : ", err)
		panic(err)
	}
	fmt.Println("Read Queue Opened")
	defer qObj_read.Close(0)

	//OPEN_GET_MESSAGE
	buffer := make([]byte, 100)
	for {
		defer handlePanic()
		len, err := GetMessage(qObj_read, buffer)

		if err != nil {
			fmt.Errorf("Error : ", err)
			panic(err.Error())
		}
		fmt.Println("msg length : ", len)
		fmt.Println("msg : ", string(buffer))
	}
	fmt.Println("Message received successfully!")

}
