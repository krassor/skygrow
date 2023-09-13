package redisBroker

import (
	"context"
	"fmt"

	"github.com/krassor/skygrow/tg-gpt-bot/internal/dto"
)

func (c *RedisClient) SubscribeOpenAiRequest(ctx context.Context, msg dto.OpenaiMsg) error  {}

go func() {
	log.Println("starting subscriber...")
	sub = client.Subscribe(channel)
	messages := sub.Channel()
	for message := range messages {
		from := strings.Split(message.Payload, ":")[0]
		//send to all websocket sessions/peers
		for user, peer := range Peers {
			if from != user { //don't recieve your own messages
				peer.WriteMessage(websocket.TextMessage, []byte(message.Payload))
			}
		}
	}
}()

}