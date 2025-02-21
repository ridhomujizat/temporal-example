package accountTypes

import (
	"onx-outgoing-go/internal/common/enum"
)

type SettingAccount struct {
	Account             string  `json:"account"`
	AccountName         string  `json:"account_name"`
	BlastTemplateID     *int    `json:"blast_template_id"`
	ChannelID           int64   `json:"channel_id"`
	CloseTemplateID     *int    `json:"close_template_id"`
	CreatedBy           int64   `json:"created_by"`
	CsatTemplateID      *int    `json:"csat_template_id"`
	EndpointAudio       string  `json:"endpoint_audio"`
	EndpointBlast       string  `json:"endpoint_blast"`
	EndpointEndSession  string  `json:"endpoint_endSession"`
	EndpointFile        string  `json:"endpoint_file"`
	EndpointImage       string  `json:"endpoint_image"`
	EndpointText        string  `json:"endpoint_text"`
	EndpointVideo       string  `json:"endpoint_video"`
	GreetingTemplateID  int64   `json:"greeting_template_id"`
	ID                  int64   `json:"id"`
	IsBlast             bool    `json:"is_blast"`
	IsCsatManual        bool    `json:"is_csat_manual"`
	Key                 string  `json:"key"`
	ParentID            string  `json:"parent_id"`
	Platform            string  `json:"platform"`
	ScheduleTemplateID  *int    `json:"schedule_template_id"`
	SettingAutoresponse string  `json:"setting_autoresponse"`
	SocioconnectToken   *string `json:"socioconnect_token"`
	UpdatedBy           int64   `json:"updated_by"`
}

type AccountEmail struct {
	Host       string             `json:"host"`
	Port       int                `json:"port"`
	Secure     bool               `json:"secure"`
	RequireTLS bool               `json:"require_tls"`
	Username   string             `json:"username"`
	Password   string             `json:"password"`
	Alias      string             `json:"alias"`
	Account    string             `json:"account"`
	Platform   enum.PlatformEmail `json:"platform"`
}
