package tgparser

import "time"

type MediaType string

const (
	Animation    MediaType = "animation"
	AudioFile    MediaType = "audio_file"
	StickerMsg   MediaType = "sticker"
	VoiceMessage MediaType = "voice_message"
	VideoMessage MediaType = "video_message"
	VideoFile    MediaType = "video_file"
)

type EntityType string

const (
	Bold   EntityType = "bold"
	Code   EntityType = "code"
	Italic EntityType = "italic"
	Link   EntityType = "link"
	Phone  EntityType = "phone"
	Plain  EntityType = "plain"
	Strkth EntityType = "strikethrough"
)

const (
	VideoMsgDimensions = 384
	DateFormat         = time.DateOnly + "T" + time.TimeOnly
)
