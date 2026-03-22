package territory

import (
	"a-star-is-born/internal/model"
	"path"
	"strings"
	"time"
)

func BuildContributorTerritory(snapshot model.Snapshot) (*model.Node, error) {
	root := &model.Node{
		Name:          "root",
		Contributions: map[string]int{},
	}

	for _, c := range snapshot.Contributors {
		for _, commit := range snapshot.Commits {
			if commit.Author == nil || commit.Author.Login != c.Login {
				continue
			}

			for _, f := range commit.Files {
				addPath(root, f.Filename, c.Login, commit.Date)
			}
		}
	}

	calcDominance(root)

	return root, nil
}

func addPath(node *model.Node, filePath, login string, commitDate time.Time) {
	parts := strings.Split(filePath, "/")
	if len(parts) == 0 {
		return
	}

	current := node
	for i, part := range parts {
		var child *model.Node
		for _, ch := range current.Children {
			if ch.Name == part {
				child = ch
				break
			}
		}

		if child == nil {
			child = &model.Node{
				Name:          part,
				Path:          path.Join(parts[:i+1]...),
				Contributions: map[string]int{},
			}
			current.Children = append(current.Children, child)
		}

		child.Contributions[login]++
		if child.LastCommitAt == nil || commitDate.After(*child.LastCommitAt) {
			child.LastCommitAt = &commitDate
		}

		current = child
	}
}

func calcDominance(node *model.Node) {
	if len(node.Contributions) > 0 {
		total := 0
		for _, v := range node.Contributions {
			total += v
		}

		maxCount := 0
		dominant := ""
		for login, v := range node.Contributions {
			if v > maxCount {
				maxCount = v
				dominant = login
			}
		}

		node.Dominant = dominant
		node.Dominance = float64(maxCount) / float64(total)
	}

	for _, ch := range node.Children {
		calcDominance(ch)
	}
}
