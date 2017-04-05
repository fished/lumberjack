package textfile

import "fmt"

func Writer(logpath string, messages chan *lumberjack.Message) {
	msg := <-messages
	fmt.Println(msg)
}
