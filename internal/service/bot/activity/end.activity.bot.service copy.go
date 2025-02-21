package botactivity

import types "onx-outgoing-go/internal/common/type"

func (a *ActivityBotService) End(payload types.PayloadBot, resutl string) (string, error) {
	return "Goodbye " + resutl, nil
}
