package main

import (
	"fmt"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

func main() {
	// Set the IBM MQ queue manager and queue information

	// IBM MQ server connection details

	qmgrName := "QM1"
	queueName := "DEV.DEAD.LETTER.QUEUE"
	channelName := "DEV.ADMIN.SVRCONN"
	connName := "localhost(1414)"
	userName := "admin"
	password := "passw0rd"
	message := "Hello world2!!!"

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
	fmt.Println("111")
	qMgr, err := ibmmq.Connx(qmgrName, gocno)
	if err != nil {
		fmt.Println("error : ", err)
		panic(err)
	}
	defer qMgr.Disc()
	fmt.Println("Queue connected!!!")

	//OPEN_PUT_QUEUE

	openOptions := ibmmq.MQOO_OUTPUT
	od := ibmmq.NewMQOD()
	od.ObjectName = queueName
	od.ObjectType = ibmmq.MQOT_Q
	qObj, err := qMgr.Open(od, openOptions)
	if err != nil {
		fmt.Errorf("Error : ", err)
		panic(err)
	}
	defer qObj.Close(0)

	//OPEN_PUT_MESSAGE

	pmd := ibmmq.NewMQMD()
	pmo := ibmmq.NewMQPMO()
	pmd.Format = ibmmq.MQFMT_STRING
	pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT | ibmmq.MQPMO_NEW_MSG_ID
	err = qObj.Put(pmd, pmo, []byte(message))
	if err != nil {
		fmt.Errorf("Error : ", err)
		panic(err)
	}
	fmt.Println("Message sent successfully!")

	//OPEN_GET_QUEUE

	openOptions = ibmmq.MQOO_INPUT_SHARED
	od = ibmmq.NewMQOD()
	od.ObjectName = queueName
	od.ObjectType = ibmmq.MQOT_Q
	qObj_read, err := qMgr.Open(od, openOptions)
	if err != nil {
		fmt.Errorf("Error : ", err)
		panic(err)
	}
	defer qObj_read.Close(0)

	//OPEN_GET_MESSAGE

	for {
		pmdr := ibmmq.NewMQMD()
		pmdr.Format = ibmmq.MQFMT_STRING
		gmo := ibmmq.NewMQGMO()
		buffer := make([]byte, 100)
		len, _ := qObj_read.Get(pmdr, gmo, buffer)
		// if err != nil {
		// 	fmt.Errorf("Error : ", err)
		// 	panic(err)
		// }
		fmt.Println("msg length : ", len)
		fmt.Println("msg : ", string(buffer))
	}
}
