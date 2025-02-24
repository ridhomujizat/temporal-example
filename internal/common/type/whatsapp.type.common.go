package types

type MessageTypeWhatsapp string

const (
	TEXT        MessageTypeWhatsapp = "text"
	IMAGE       MessageTypeWhatsapp = "image"
	CONTACTS    MessageTypeWhatsapp = "contacts"
	DOCUMENT    MessageTypeWhatsapp = "document"
	INTERACTIVE MessageTypeWhatsapp = "interactive"
	BUTTON      MessageTypeWhatsapp = "button"
	LOCATION    MessageTypeWhatsapp = "location"
	VIDEO       MessageTypeWhatsapp = "video"
	STICKER     MessageTypeWhatsapp = "sticker"
	ORDER       MessageTypeWhatsapp = "order"
	UNKNOWN     MessageTypeWhatsapp = "unknown"
	VOICE       MessageTypeWhatsapp = "voice"
	EPHEMERAL   MessageTypeWhatsapp = "ephemeral"
)

type IncomingWhatspp struct {
	TenantId  string             `json:"tenant_id" validate:"required"`
	AccountId string             `json:"account_id" validate:"required"`
	TunnelUrl string             `json:"tunnel_url" validate:"required,url"`
	Contacts  []ContactsWhatsapp `json:"contacts" validate:"omitempty,dive"`
	Messages  []MessagesWhatsapp `json:"messages" validate:"omitempty,dive"`
}

type ContactsWhatsapp struct {
	WaId    string          `json:"wa_id" validate:"required"`
	Profile ProfileWhatsapp `json:"profile" validate:"omitempty"`
}

type ProfileWhatsapp struct {
	Name string `json:"name" validate:"required"`
}

type MessagesWhatsapp struct {
	Timestamp   string                     `json:"timestamp" validate:"required"`
	Type        MessageTypeWhatsapp        `json:"type" validate:"required"`
	From        string                     `json:"from" validate:"required"`
	Text        TextWhatsapp               `json:"text" validate:"required"`
	Interactive InteractiveMessageWhatsapp `json:"interactive" validate:"required"`
	Location    LocationWhatsapp           `json:"location" validate:"omitempty"`
}

type LocationWhatsapp struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type TextWhatsapp struct {
	Body string `json:"body" validate:"required"`
}
type BodyWhatsapp struct {
	Text string `json:"text" validate:"required"`
}
type InteractiveMessageWhatsapp struct {
	Type         string                       `json:"type" validate:"omitempty,oneof='list button'"`
	List_reply   InteractiveListReplyWhatsapp `json:"list_reply" validate:"omitempty"`
	Button_reply InteractiveListReplyWhatsapp `json:"button_reply" validate:"omitempty"`
	NVM_reply    NVMReplyWhatsapp             `json:"nvm_reply" validate:"omitempty"`
}

type NVMReplyWhatsapp struct {
	Body        string `json:"body" validate:"required"`
	Name        string `json:"name" validate:"required"`
	ReponseJSON string `json:"response_json" validate:"required"`
}

type InteractiveListReplyWhatsapp struct {
	Title string `json:"title" validate:"required"`
	Id    string `json:"id" validate:"required"`
}

// Line struct for outgoing whatsapp message =============
type OutgoingTextWhatsapp struct {
	RecipientType    string       `json:"recipient_type"`
	MessagingProduct string       `json:"messaging_product"`
	To               string       `json:"to"`
	Type             string       `json:"type"`
	Text             TextWhatsapp `json:"text"`
}
