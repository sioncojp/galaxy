package galaxy

import "fmt"

// SelectAllCommits ...select * from commits in DB.galaxy
func (config *Config) SelectAllCommits() ([]Commits, error) {
	var commits []Commits

	db, err := config.DBConnect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	db.Find(&commits)
	return commits, nil
}

// ShowUrlString ...Output Url String
func ShowUrlString(cn, url string) string {
	return fmt.Sprintf("http://%s-%s/", cn[:7], url)
}
