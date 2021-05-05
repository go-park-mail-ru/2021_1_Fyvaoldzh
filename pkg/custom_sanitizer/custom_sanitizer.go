package custom_sanitizer

import (
	"github.com/microcosm-cc/bluemonday"
	"kudago/application/models"
)

type CustomSanitizer struct {
	sanitizer *bluemonday.Policy
}

func NewCustomSanitizer(sz *bluemonday.Policy) *CustomSanitizer {
	customSanitizer := CustomSanitizer{sanitizer: sz}
	return &customSanitizer
}

func (cs *CustomSanitizer) SanitizeOwnProfile(profile *models.UserOwnProfile) {
	profile.Name = cs.sanitizer.Sanitize(profile.Name)
	profile.Avatar = cs.sanitizer.Sanitize(profile.Avatar)
	profile.Login = cs.sanitizer.Sanitize(profile.Login)
	profile.NewPassword = cs.sanitizer.Sanitize(profile.NewPassword)
	profile.OldPassword = cs.sanitizer.Sanitize(profile.OldPassword)
	profile.City = cs.sanitizer.Sanitize(profile.City)
	profile.Email = cs.sanitizer.Sanitize(profile.Email)
	profile.About = cs.sanitizer.Sanitize(profile.About)
}

func (cs *CustomSanitizer) SanitizeOtherProfile(profile *models.OtherUserProfile) {
	profile.Name = cs.sanitizer.Sanitize(profile.Name)
	profile.Avatar = cs.sanitizer.Sanitize(profile.Avatar)
	profile.City = cs.sanitizer.Sanitize(profile.City)
	profile.About = cs.sanitizer.Sanitize(profile.About)
}

func (cs *CustomSanitizer) SanitizeEventCards(events models.EventCards) models.EventCards {
	var newEvents models.EventCards
	for _, elem := range events {
		elem.Place = cs.sanitizer.Sanitize(elem.Place)
		elem.Description = cs.sanitizer.Sanitize(elem.Description)
		elem.Title = cs.sanitizer.Sanitize(elem.Title)
		newEvents = append(newEvents, elem)
	}

	return newEvents
}

func (cs *CustomSanitizer) SanitizeUsersOnEvent(users models.UsersOnEvent) models.UsersOnEvent {
	var newUsers models.UsersOnEvent
	for _, elem := range users {
		elem.Avatar = cs.sanitizer.Sanitize(elem.Avatar)
		elem.Name = cs.sanitizer.Sanitize(elem.Name)
		newUsers = append(newUsers, elem)
	}

	return newUsers
}

func (cs *CustomSanitizer) SanitizeUserCards(users models.UserCards) models.UserCards {
	var newUsers models.UserCards
	for _, elem := range users {
		elem.City = cs.sanitizer.Sanitize(elem.City)
		elem.Avatar = cs.sanitizer.Sanitize(elem.Avatar)
		elem.Name = cs.sanitizer.Sanitize(elem.Name)
		newUsers = append(newUsers, elem)
	}

	return newUsers
}

func (cs *CustomSanitizer) SanitizeEvent(elem *models.Event) {
	elem.Title = cs.sanitizer.Sanitize(elem.Title)
	elem.Description = cs.sanitizer.Sanitize(elem.Description)
	elem.Street = cs.sanitizer.Sanitize(elem.Street)
	elem.Image = cs.sanitizer.Sanitize(elem.Image)
	elem.Subway = cs.sanitizer.Sanitize(elem.Subway)
	elem.Place = cs.sanitizer.Sanitize(elem.Place)
	elem.Category = cs.sanitizer.Sanitize(elem.Category)
	elem.Followers = cs.SanitizeUsersOnEvent(elem.Followers)
	elem.Tags = cs.SanitizeTags(elem.Tags)
}

func (cs *CustomSanitizer) SanitizeTags(tags models.Tags) models.Tags {
	var newTags models.Tags
	for _, elem := range tags {
		elem.Name = cs.sanitizer.Sanitize(elem.Name)
		newTags = append(newTags, elem)
	}
	return newTags
}

func (cs *CustomSanitizer) SanitizeActions(actions models.ActionCards) models.ActionCards {
	var newActions models.ActionCards
	for _, elem := range actions {
		elem.Name1 = cs.sanitizer.Sanitize(elem.Name1)
		elem.Name2 = cs.sanitizer.Sanitize(elem.Name2)
		newActions = append(newActions, elem)
	}
	return newActions
}
func (cs *CustomSanitizer) SanitizeDialogueCards(dialogues models.DialogueCards) models.DialogueCards {
	var sanitizeDialogues models.DialogueCards
	for _, elem := range dialogues {
		elem.Interlocutor.Avatar = cs.sanitizer.Sanitize(elem.Interlocutor.Avatar)
		elem.Interlocutor.Name = cs.sanitizer.Sanitize(elem.Interlocutor.Name)
		elem.LastMessage.Text = cs.sanitizer.Sanitize(elem.LastMessage.Text)
		sanitizeDialogues = append(sanitizeDialogues, elem)
	}

	return sanitizeDialogues
}

func (cs *CustomSanitizer) SanitizeMessages(messages models.Messages) models.Messages {
	var sanitizeMessages models.Messages
	for _, elem := range messages {
		elem.Text = cs.sanitizer.Sanitize(elem.Text)
		sanitizeMessages = append(sanitizeMessages, elem)
	}

	return sanitizeMessages
}

func (cs *CustomSanitizer) SanitizeDialogue(dialogue *models.Dialogue) {
	dialogue.Interlocutor.Name = cs.sanitizer.Sanitize(dialogue.Interlocutor.Name)
	dialogue.Interlocutor.Avatar = cs.sanitizer.Sanitize(dialogue.Interlocutor.Avatar)
	dialogue.DialogMessages = cs.SanitizeMessages(dialogue.DialogMessages)
}
