package twitch

type ChannelInfo struct {
	DisplayName     string `json:"display_name"`
	Login           string `json:"login"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ID              string `json:"id"`
}
