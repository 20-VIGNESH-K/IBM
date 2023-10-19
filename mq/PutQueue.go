package mq

import "github.com/ibm-messaging/mq-golang/v5/ibmmq"

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
