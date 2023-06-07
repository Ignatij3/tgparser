package tgparser

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
)

type MsgText string

type (
	// Message represents single chat message in telegram.
	Message struct {
		info *message
	}

	// message is same as "Message", but for internal use.
	message struct {
		General
		MediaInfo
		MediaType MediaType `json:"media_type"`
		MimeType  string    `json:"mime_type"`
	}

	// MediaInfo contains information about media present in message.
	// If message does not contain any media, fields here are set to zero values.
	MediaInfo struct {
		File      string `json:"file"`
		Thumbnail string `json:"thumbnail"`
		Duration  int    `json:"duration_seconds"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
	}

	// General contains information which may be present in any type of message.
	General struct {
		ID       int      `json:"id"`
		Type     string   `json:"type"`
		Date     string   `json:"date"`
		DateUnix string   `json:"date_unixtime"`
		From     string   `json:"from"`
		FromID   string   `json:"from_id"`
		Reply    int      `json:"reply_to_message_id"`
		Edited   string   `json:"edited"`
		EditUnix string   `json:"edited_unixtime"`
		Text     MsgText  `json:"text"`
		TextEnts []Entity `json:"text_entities"`
	}

	Entity struct {
		Type EntityType `json:"type"`
		Text string     `json:"text"`
	}
)

func (mt *MsgText) UnmarshalJSON(text []byte) error {
	t := strings.TrimSpace(string(text))
	if strings.HasPrefix(t, "[") {
		ents := make([]Entity, 0)
		err := json.Unmarshal(text, &ents)
		for _, entity := range ents {
			*mt += MsgText(entity.Text)
		}

		if err != nil {
			parsedText := parseComplicatedText([]byte(t))
			*mt = MsgText(parsedText)
			err = nil
		}

		return err
	}

	var data string
	err := json.Unmarshal(text, &data)
	*mt = MsgText(data)

	return err
}

// parseComplicatedText parses text field containing either plain text and text entities.
func parseComplicatedText(text []byte) string {
	var res string
	decoder := json.NewDecoder(bytes.NewReader(text))

	decoder.Token() // [
	for {
		t, _ := decoder.Token()
		if bracket, ok := t.(json.Delim); ok && bracket.String() == "{" {
			decoder.Token() // "type"
			decoder.Token() // ${type} (link, most likely)
			decoder.Token() // "text"

			t, _ = decoder.Token()
			res += t.(string)
			decoder.Token() // }

		} else if bracket, ok := t.(json.Delim); ok && bracket.String() == "]" {
			break
		} else if data, ok := t.(string); ok {
			res += data
		}
	}

	return res
}

// ParseMessage unmashals single telegram message and returns it.
// Message must be export in json format, otherwise error will occur.
func ParseMessage(data []byte) (Message, error) {
	var info message

	if err := json.Unmarshal(data, &info); err != nil {
		return Message{}, err
	}

	return Message{&info}, nil
}

// GetMediaType returns media type of the media (if the message has one) or nothing.
func (m Message) GetMediaType() MediaType {
	return m.info.MediaType
}

// GetMimeType returns mime type of the media (if the message has one) or nothing.
func (m Message) GetMimeType() string {
	return m.info.MimeType
}

// GeneralInfo returns information common for any message.
// Exceptions are that "Edited" and "EditUnix" are never present in round_video_message message type and therefore are set to zero value.
func (m Message) GeneralInfo() General {
	return m.info.General
}

// MediaInfo returns media information for message if present.
// If media is not present in the message, all fields returned contain no information.
func (m Message) MediaInfo() MediaInfo {
	return m.info.MediaInfo
}

// Date returns timestamp of message converted to time.Time for convenience.
func (m Message) Date() time.Time {
	date, _ := time.Parse(DateFormat, m.info.General.Date)
	return date
}

// Text returns plain text present in message.
func (m Message) Text() string {
	return string(m.info.General.Text)
}

// IsMedia returns whether message contains media or not.
func (m Message) IsMedia() bool {
	return m.info.MediaType != ""
}
