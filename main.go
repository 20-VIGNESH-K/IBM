// package main

// import (
//     "fmt"
//     "log"
//     "github.com/streadway/amqp"
// )

// func main() {
//     // Connection details
//     amqpURI := "amqp://guest:guest@localhost:5672"

//     // Connect to RabbitMQ
//     conn, err := amqp.Dial(amqpURI)
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer conn.Close()

//     // Create a channel
//     ch, err := conn.Channel()
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer ch.Close()

//     // Declare a queue
//     queueName := "myQueue"
//     _, err = ch.QueueDeclare(
//         queueName, // Queue name
//         true,      // Durable
//         false,     // Delete when unused
//         false,     // Exclusive
//         false,     // No-wait
//         nil,       // Arguments
//     )
//     if err != nil {
//         log.Fatal(err)
//     }

//     // Define the message
//     message := "Hello, RabbitMQ"

//     // Publish the message to the queue
//     err = ch.Publish(
//         "",        // Exchange
//         queueName, // Routing key
//         false,     // Mandatory
//         false,     // Immediate
//         amqp.Publishing{
//             ContentType: "text/plain",
//             Body:        []byte(message),
//         })
//     if err != nil {
//         log.Fatal(err)
//     }

//     fmt.Println("Message sent successfully!")
// }

// package main

// import (
// 	"fmt"

// 	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
// )

// func main() {
// 	// Set the IBM MQ connection information
// 	qmgrName := "QM1"            // Queue Manager name
// 	queueName := " DEV.DEAD.LETTER.QUEUE"     // Name of the target queue
// 	//connName := "localhost(1414)" // Host and port where IBM MQ is running
// 	//channelName := "YOUR_CHANNEL" // Channel name
// 	user := "YOUR_USER"         // Your username
// 	password := "YOUR_PASSWORD" // Your password

// 	// Connect to the IBM MQ queue manager
// 	// qMgr, err := ibmmq.Connx(qmgrName, connName, ibmmq.MQCNO_CLIENT_BINDING,ibmmq.MQCNO
// 	// 	ibmmq.MQCNO_HANDLE_SHARE_BLOCK)
// 	MQCNO := ibmmq.NewMQCNO()
// 	qMgr, err := ibmmq.Connx(qmgrName, MQCNO)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer qMgr.Disc()

// 	// Set the connection authentication information
// 	auth := ibmmq.NewMQCSP()
// 	auth.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
// 	auth.UserId = user
// 	auth.Password = password

// 	// // Authenticate the connection
// 	// err = qMgr.SetClientConn(&auth)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Open the target queue for putting messages
// 	putqObject, err := qMgr.OpenQueue(queueName, ibmmq.MQOO_OUTPUT)

// 	if err != nil {
// 		panic(err)
// 	}
// 	defer putqObject.Close(0)

// 	// Create a message to send
// 	putMessage := ibmmq.NewMQMD()
// 	putMessage.Format = ibmmq.MQFMT_STRING
// 	messageData := []byte("Hello, IBM MQ!")

// 	// Put the message to the queue
// 	err = putqObject.Put(putMessage, nil, messageData)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Message sent successfully!")
// }

package main

import (
	"fmt"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

func main() {
	// Set the IBM MQ queue manager and queue information
	qmgrName := "QM1"
	fmt.Println(qmgrName)
	queueName := "DEV.DEAD.LETTER.QUEUE"
	channelName := "DEV.ADMIN.SVRCONN"
	fmt.Println(queueName)
	connName := "localhost(1414)" // IBM MQ server connection details
	fmt.Println(connName)
	userName := "admin"
	password := "passw0rd"

	message := "Hello world2!!!"

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

	// Create a PUT message options
	// pmo := ibmmq.NewMQPMO()

	// // Set the persistence and put message options
	// pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT
	// pmo.Options |= ibmmq.MQPMO_NEW_MSG_ID
	// pmo.Options |= ibmmq.MQPMO_NEW_CORREL_ID

	// Create a message to send
	// msg := ibmmq.NewMQMD()

	// //msg.Format = ibmmq.MQFMT_STRING
	// // msg.ReplyToQ = queueName
	// // Set this to the appropriate reply-to queue
	// messageData := []byte("Hello, IBM MQ!")

	// // Put the message to the queue
	// //putmqmd := ibmmq.NewMQMD()
	// good := ibmmq.NewMQOD()
	// err = qMgr.Put1(good, msg, pmo, messageData)
	// if err != nil {
	// 	panic(err)
	// }

	// Clean up and disconnect from the queue manager
	qMgr.Disc()
}
