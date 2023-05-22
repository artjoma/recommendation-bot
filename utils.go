package main

func panicOnErr(err error, msg string) {
	if err != nil {
		panic("msg:" + msg + " err:" + err.Error())
	}
}
