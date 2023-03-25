package validation

import (
	"errors"
	"github.com/Ghytro/galleryapp/internal/view/polls"
)

func ValidateCreatedPoll(model *polls.NewPollRequest) error {
	if model.Topic == "" {
		return errors.New("тема опроса пуста")
	}
	if len(model.Topic) > 100 {
		return errors.New("тема опроса не может быть длиннее 100 символов")
	}
	if len(model.Options) == 0 {
		return errors.New("у опроса нет вариантов ответа")
	}
	for _, opt := range model.Options {
		if len(opt) > 100 {
			return errors.New("опция не может быть длиннее 100 символов")
		}
	}
	return nil
}
