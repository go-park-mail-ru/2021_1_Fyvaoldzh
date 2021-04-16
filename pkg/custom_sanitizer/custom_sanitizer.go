package custom_sanitizer

import (
	"kudago/application/models"

	"github.com/microcosm-cc/bluemonday"
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
	profile.Planning = cs.SanitizeEventCards(profile.Planning)
	profile.Visited = cs.SanitizeEventCards(profile.Visited)
}

func (cs *CustomSanitizer) SanitizeOtherProfile(profile *models.OtherUserProfile) {
	profile.Name = cs.sanitizer.Sanitize(profile.Name)
	profile.Avatar = cs.sanitizer.Sanitize(profile.Avatar)
	profile.City = cs.sanitizer.Sanitize(profile.City)
	profile.About = cs.sanitizer.Sanitize(profile.About)
	profile.Planning = cs.SanitizeEventCards(profile.Planning)
	profile.Visited = cs.SanitizeEventCards(profile.Visited)
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
