package internal

import (
	"fmt"

	"github.com/afrxo/fig/auth"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/client"
	githttp "github.com/go-git/go-git/v6/plumbing/transport/http"
)

func CloneRepository(url string, dest string, creds auth.Credentials) (*git.Repository, error) {
	return git.PlainClone(dest, &git.CloneOptions{
		URL: url,
		ClientOptions: []client.Option{
			client.WithHTTPAuth(&githttp.BasicAuth{
				Username: creds.Username,
				Password: creds.Token,
			}),
		},
	})
}

func CheckoutRepository(r *git.Repository, ref string) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	_, err = r.Reference(plumbing.NewBranchReferenceName(ref), false)
	if err == nil {
		return w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(ref),
		})
	}

	_, err = r.Reference(plumbing.NewTagReferenceName(ref), false)
	if err == nil {
		return w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewTagReferenceName(ref),
		})
	}

	hash, err := r.ResolveRevision(plumbing.Revision(ref))
	if err == nil {
		return w.Checkout(&git.CheckoutOptions{
			Hash: *hash,
		})
	}

	return fmt.Errorf("ref not found: %s", ref)
}
