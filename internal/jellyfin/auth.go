package jellyfin

import "fmt"

// Ping hits /System/Info/Public without a session token.
func (c *Client) Ping() (PublicInfo, error) {
	var info PublicInfo
	err := c.get(pathPublicInfo, nil, &info)
	return info, err
}

// Login authenticates by name and stores the session token on c.
func (c *Client) Login(user, pass string) error {
	var out AuthResult
	err := c.post(pathAuth, nil, map[string]string{
		"Username": user,
		"Pw":       pass,
		"Password": pass,
	}, &out)
	if err != nil {
		return err
	}
	if out.AccessToken == "" || out.User.ID == "" {
		return fmt.Errorf("login returned empty session")
	}
	c.Token = out.AccessToken
	c.UserID = out.User.ID
	c.UserName = out.User.Name
	c.ServerID = out.ServerID
	return nil
}

// Me loads the current user and refreshes UserID / UserName.
func (c *Client) Me() (User, error) {
	var u User
	err := c.get(pathMe, nil, &u)
	if err == nil && u.ID != "" {
		c.UserID = u.ID
		c.UserName = u.Name
	}
	return u, err
}
