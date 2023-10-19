package main

import (
	"fmt"
	"ibm/mq"
)

func ErrorMessage(err error) {
	if err != nil {
		fmt.Println("error : ", err)
		panic(err)
	}
}

func handlePanic() {
	if a := recover(); a != nil {
		fmt.Println("All messages are read and No new messages to read now")
	}
}
func main() {

	qmgrName := "QM1"
	queueName := "DEV.DEAD.LETTER.QUEUE"
	message := "Hello world!"

	//QUEUE_CONNECTION
	qMgr, err := mq.QueueConnection(qmgrName)
	ErrorMessage(err)
	fmt.Println("Queue connected!!!")
	defer qMgr.Disc()

	//OPEN_PUT_QUEUE
	qObj, err := mq.OpenPutQueue(queueName, qMgr)
	ErrorMessage(err)
	fmt.Println("Write Queue Opened")
	defer qObj.Close(0)

    //OPEN_PUT_MESSAGE
	err = mq.PutMessage(message, qObj)
	ErrorMessage(err)
	fmt.Println("Message sent successfully!")

	// OPEN_GET_QUEUE
	qObj_read, err := mq.OpenGetQueue(queueName, qMgr)
	ErrorMessage(err)
	fmt.Println("Read Queue Opened")
	defer qObj_read.Close(0)

	//OPEN_GET_MESSAGE
	buffer := make([]byte, 100)
	for {
		defer handlePanic()
		len, err := mq.GetMessage(qObj_read, buffer)
		ErrorMessage(err)
		fmt.Println("Message received successfully!")
		fmt.Println("msg length : ", len)
		fmt.Println("msg : ", string(buffer))	
	}
}
