package enum

type Platform string

const (
	SOCIO            Platform = "socioconnect"
	CAPIWHA          Platform = "capiwha"
	CHATAPI          Platform = "chatapi"
	WHATSAPPOFFICIAL Platform = "whatsappOfficial"
	PLATFORMDEX      Platform = "platformdex"
	DOLPHIN          Platform = "3dolphin"
	WAVECELL         Platform = "wavecell"
	WAPPIN           Platform = "wappin"
	OCA              Platform = "oca"
	TELEGRAMAPI      Platform = "telegramApi"
	OCTOPUSHCHAT     Platform = "octopushChat"
	KATAAI           Platform = "kata_ai"
	LIBRA            Platform = "libra"
	JATIS            Platform = "jatis"
	BOTIKA           Platform = "botika"
	QISCUSS          Platform = "qiscuss"
	INFOBIP          Platform = "infobip"
	NADYNE           Platform = "nadyne"
	SMSTURBO         Platform = "smsturbo"
	MAYTAPI          Platform = "maytapi"
	DIPS             Platform = "dips"
)

type PlatformEmail string

const (
	NODEMAILER PlatformEmail = "nodemailer"
	OUTLOOK    PlatformEmail = "outlook"
	ENGINE     PlatformEmail = "emailEngine"
	SENDINBLUE PlatformEmail = "sendinblue"
	MAILGUN    PlatformEmail = "mailgun"
)

type BotikaChannelType string

const (
	Facebook BotikaChannelType = "FBMESSENGER"
	LINE     BotikaChannelType = "LINE"
	WhatsApp BotikaChannelType = "OAWHATSAPPBOTIKA"
	Wa       BotikaChannelType = "OAWHATSAPPINFOMED"
	Chat     BotikaChannelType = "CHATBOTIKAWEBCHAT"
	Igdm     BotikaChannelType = "IGDM"
)
