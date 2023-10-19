package mq

import (
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

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
