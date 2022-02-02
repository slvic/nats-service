package types

type Message struct {
	Id   string `field:"id"`
	Data []byte `field:"data"`
}
