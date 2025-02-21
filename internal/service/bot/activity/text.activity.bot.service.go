package botactivity

import types "onx-outgoing-go/internal/common/type"

func (a *ActivityBotService) Text(payload types.PayloadBot) (string, error) {
	return "Hello " + payload.MetaData.UniqueId, nil
}
