package acp

import "encoding/json"

// ContentBlockType 定义内容类型。
type ContentBlockType string

const (
	ContentBlockTypeText         ContentBlockType = "text"
	ContentBlockTypeImage        ContentBlockType = "image"
	ContentBlockTypeAudio        ContentBlockType = "audio"
	ContentBlockTypeResourceLink ContentBlockType = "resource_link"
	ContentBlockTypeResource     ContentBlockType = "resource"
)

// ContentBlock 表示一段内容。
type ContentBlock struct {
	Type        ContentBlockType `json:"type"`
	Text        string           `json:"text,omitempty"`
	Data        string           `json:"data,omitempty"`
	MimeType    string           `json:"mimeType,omitempty"`
	URI         string           `json:"uri,omitempty"`
	Description string           `json:"description,omitempty"`
	Name        string           `json:"name,omitempty"`
	Meta        json.RawMessage  `json:"_meta,omitempty"`
}

// NewTextContentBlock 创建文本内容。
func NewTextContentBlock(text string) ContentBlock {
	return ContentBlock{Type: ContentBlockTypeText, Text: text}
}
