package otxapi



type OTXPulseDetailService struct {
	client *Client
}

// Pulse represents an OTX Pulse
type PulseDetail struct {
	ID               *string                 `json:"id"`
	Author           *string                 `json:"author_name"`
	Name             *string                 `json:"name"`
	Description      *string                 `json:"description,omitempty"`
	CreatedAt        *string                 `json:"created,omitempty"`
	ModifiedAt       *string                 `json:"modified"`
	References       []string                `json:"references,omitempty"`
	Tags             []string                `json:"tags,omitempty"`
	Indicators       []struct {
		ID               *string                 `json:"_id"`
		Indicator        *string                 `json:"indicator"`
		Type             *string                 `json:"type"`
		Description      *string                 `json:"description,omitempty"`
	}                                        `json:"indicators,omitempty"`
	Revision         *float32                `json:"revision,omitempty"`
}

func (r PulseDetail) String() string {
	return Stringify(r)
}

type OTXThreatIntelFeedService struct {
	client *Client
}

type ThreatIntelFeed struct {
	Pulses      []PulseDetail     `json:"results"`
	// These fields provide the page values for paginating through a set of
	// results.  Any or all of these may be set to the zero value for
	// responses that are not part of a paginated set, or for which there
	// are no additional pages.
	//NextPageNum  int   Coming soon
	//PrevPageNum  int   Coming soon
	NextPageString  *string	     `json:"next"`
	PrevPageString  *string      `json:"prev"`
	Count           int          `json:"count"`
}

func (r ThreatIntelFeed) String() string {
	return Stringify(r)
}
