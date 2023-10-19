package mq

import "github.com/ibm-messaging/mq-golang/v5/ibmmq"

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
