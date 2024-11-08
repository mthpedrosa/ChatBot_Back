package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*type Mensagem struct {
	EnviadoDe string `json:"enviado_de" bson:"enviado_de"` // ID de quem enviou a mensagem (pode ser sua ID)
	Conteudo  string `json:"conteudo"`
	Timestamp int64  `json:"timestamp"`
}*/

type Conversation struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerId  string             `json:"customer_id" bson:"customer_id"`
	Messages    []Message          `json:"mensagens" bson:"mensagens,omitempty"`
	AssistantId string             `json:"assistant_id" bson:"assistant_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"`
	OtherFields []Fields           `json:"other_fields" bson:"other_fields"`
}
