package nbiot

import "fmt"

// Team represents a team.
type Team struct {
	TeamID  string            `json:"teamId"`
	Members []Member          `json:"members,omitempty"`
	Tags    map[string]string `json:"tags,omitempty"`
}

// Member is the member element of the Team type
type Member struct {
	UserID        *string `json:"userId"`
	Role          *string `json:"role"`
	Name          *string `json:"name"`
	Email         *string `json:"email"`
	Phone         *string `json:"phone"`
	VerifiedEmail *bool   `json:"verifiedEmail"`
	VerifiedPhone *bool   `json:"verifiedPhone"`
	ConnectID     *string `json:"connectId"`
	GitHubLogin   *string `json:"gitHubLogin"`
	AuthType      *string `json:"authType"`
	AvatarURL     *string `json:"avatarUrl"`
}

// Team gets a team.
func (c *Client) Team(id string) (Team, error) {
	var team Team
	err := c.get("/teams/"+id, &team)
	return team, err
}

// Teams gets all teams that the user belongs to.
func (c *Client) Teams() ([]Team, error) {
	var teams struct {
		Teams []Team `json:"teams"`
	}
	err := c.get("/teams", &teams)
	return teams.Teams, err
}

// CreateTeam creates a team.
func (c *Client) CreateTeam(team Team) (Team, error) {
	err := c.create("/teams", &team)
	return team, err
}

// UpdateTeam updates a team, but not its members.
// No tags are deleted, only added or updated.
func (c *Client) UpdateTeam(team Team) (Team, error) {
	err := c.update("/teams/"+team.TeamID, &team)
	return team, err
}

// UpdateTeamMemberRole updates the role of a team member.
func (c *Client) UpdateTeamMemberRole(teamID, userID, role string) error {
	return c.update(fmt.Sprintf("/teams/%s/members/%s", teamID, userID), Member{Role: &role})
}

// DeleteTeamMember deletes a team member.
func (c *Client) DeleteTeamMember(teamID, userID string) error {
	return c.delete(fmt.Sprintf("/teams/%s/members/%s", teamID, userID))
}

// DeleteTeamTag deletes a tag from a team.
func (c *Client) DeleteTeamTag(id, name string) error {
	return c.delete(fmt.Sprintf("/teams/%s/tags/%s", id, name))
}

// DeleteTeam deletes a team.
func (c *Client) DeleteTeam(id string) error {
	return c.delete("/teams/" + id)
}

// Invite is an invitation to a team.
type Invite struct {
	Code      string `json:"code"`
	CreatedAt int    `json:"createdAt"` // time stamp (in ms)
}

// Invite gets an invite.
func (c *Client) Invite(teamID, code string) (Invite, error) {
	var invite Invite
	err := c.get(fmt.Sprintf("/teams/%s/invites/%s", teamID, code), &invite)
	return invite, err
}

// Invites gets all invites for the team.
func (c *Client) Invites(teamID string) ([]Invite, error) {
	var invites struct {
		Invites []Invite `json:"invites"`
	}
	err := c.get(fmt.Sprintf("/teams/%s/invites", teamID), &invites)
	return invites.Invites, err
}

// CreateInvite creates a invite.
func (c *Client) CreateInvite(teamID string) (Invite, error) {
	var invite Invite
	err := c.create(fmt.Sprintf("/teams/%s/invites", teamID), &invite)
	return invite, err
}

// AcceptInvite accepts a invite.
func (c *Client) AcceptInvite(code string) error {
	return c.create("/teams/accept", &Invite{Code: code})
}

// DeleteInvite deletes an invite.
func (c *Client) DeleteInvite(teamID, code string) error {
	return c.delete(fmt.Sprintf("/teams/%s/invites/%s", teamID, code))
}
