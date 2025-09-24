package converter

import (
	"encoding/json"

	"post/internal/model"
	"post/pkg/logger"
)

func UserFromPayload(bs []byte) model.User {
	const op = "converter.UserFromPayload"
	var payload model.User

	if err := json.Unmarshal(bs, &payload); err != nil {
		logger.Warn(op, "err", err)

		return model.User{}
	}

	return model.User{
		ID: payload.ID,
	}
}

