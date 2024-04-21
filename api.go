package api

import "time"

// A Decision is the top level configuration that everything else falls from. It captures the
// prompt being decided and contains all of the possible choices for that prompt. An instantiation
// of a Decision is Bracket.
type Decision struct {
	ID      ID       `json:"id"`
	OwnerID ID       `json:"ownerId"`
	Public  bool     `json:"public"`
	Prompt  string   `json:"prompt"`
	Choices []Choice `json:"choices"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

// A Choice represents a single, discrete, possible answer to a decision. These are configured once
// and then exist for reference as a choice bank for future Decisions
type Choice struct {
	ID       ID     `json:"id"`
	OwnerID  ID     `json:"ownerId"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// A Bracket is an instantiation of a Decision. Lots of decisions have to be made over and over
// again (e.g what are we having for dinner?), and that's where Brackets come in. A Bracket will
// eventually have a winning Choice and results can be referenced later on
type Bracket struct {
	ID         ID `json:"id"`
	OwnerID    ID `json:"ownerId"`
	DecisionID ID `json:"decisionId"`
	WinnerID   ID `json:"winnerId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// A Match represents a single pairing of two choices for a given Bracket. When a Choice is selected
// the WinnerID is set appropriately
type Match struct {
	ID        ID `json:"id"`
	BracketID ID `json:"bracketId"`

	Left     Choice `json:"left"`
	Right    Choice `json:"right"`
	WinnerID ID     `json:"winnerId"`
	Round    int    `json:"round"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NewMatch struct {
	Round   int `json:"round"`
	LeftID  ID  `json:"leftId"`
	RightID ID  `json:"rightId"`
}
