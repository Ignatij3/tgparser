package tgparser

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// TODO: add support for group chats

// Conversation holds information present in telegram personal chat.
type Conversation struct {
	FirstPersonName string
	FirstPersonID   string
	PartnerName     string
	PartnerID       string
	msgs            []Message
	// calls []Service TO BE IMPLEMENTED
}

// ParseConversation accepts unmodified json chat export and returns parsed conversation if successful.
func ParseConversation(data []byte) (Conversation, error) {
	conv := Conversation{
		msgs: make([]Message, 0),
	}

	messages := make([]message, 0)

	buffer := bytes.NewBuffer([]byte{})
	json.Compact(buffer, data)
	decoder := json.NewDecoder(buffer)

	decoder.Token() // {

	decoder.Token() // "name"
	if err := decoder.Decode(&conv.PartnerName); err != nil {
		return Conversation{}, err
	}

	decoder.Token()     // "type"
	decoder.Decode(nil) // "personal_chat"

	decoder.Token() // "id"
	prID := 0
	if err := decoder.Decode(&prID); err != nil {
		return Conversation{}, err
	}
	conv.PartnerID = "user" + strconv.Itoa(prID)

	decoder.Token() // "messages"
	if err := decoder.Decode(&messages); err != nil {
		return Conversation{}, err
	}

	decoder.Token()

	foundFirstPerson := false
	for i := range messages {
		if messages[i].Type != "service" { // ignoring calls, for now
			if !foundFirstPerson && messages[i].FromID != conv.PartnerID {
				conv.FirstPersonID = messages[i].FromID
				conv.FirstPersonName = messages[i].From
				foundFirstPerson = true
			}
			conv.msgs = append(conv.msgs,
				Message{
					info: &messages[i],
				})
		}
	}

	return conv, nil
}

// Messages returns copy of array containing all messages present in chat.
func (c Conversation) Messages() []Message {
	messages := make([]Message, len(c.msgs))
	copy(messages, c.msgs)
	return messages
}
