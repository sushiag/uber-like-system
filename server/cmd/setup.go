package ws

type Client struct {
	clients map[uint64][]*ClientConn,
}

func NewClient() *Client {
	return &Client{
		clients: make(map[uint64][]*ClientConn),
	}
	
}