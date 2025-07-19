package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Client defines the interface for Git operations
type Client interface {
	// Repository operations
	Open(path string) error
	GetRepository() (*Repository, error)
	GetWorktree() (*Worktree, error)
	IsRepository() bool
	
	// Reference operations
	GetHead() (*Ref, error)
	GetBranches() ([]*Ref, error)
	GetTags() ([]*Ref, error)
	GetRemotes() ([]*Remote, error)
	
	// Commit operations
	GetCommit(hash string) (*Commit, error)
	GetCommits(opts *LogOptions) ([]*Commit, error)
	GetLogCount() (int, error)
	
	// Status and file operations
	GetStatus() (*Status, error)
	GetDiff(path string) (*Diff, error)
	GetFiles(path string) ([]*File, error)
	
	// Staging operations
	StageFile(path string) error
	UnstageFile(path string) error
	StageAll() error
	UnstageAll() error
	DiscardChanges(path string) error
	
	// Commit operations
	Commit(message string, opts *CommitOptions) error
	
	// Stash operations
	GetStashes() ([]*Stash, error)
	
	// Utility operations
	GetRootPath() string
	GetRelativePath(path string) string
	ExecuteCommand(args ...string) ([]byte, error)
}

// Repository represents a Git repository
type Repository struct {
	Path string
	repo *git.Repository
}

// Worktree represents a Git worktree
type Worktree struct {
	wt *git.Worktree
}

// Ref represents a Git reference (branch, tag, etc.)
type Ref struct {
	Name   string
	Type   RefType
	Hash   string
	Target string // For symbolic refs
}

// RefType defines the type of reference
type RefType int

const (
	RefTypeBranch RefType = iota
	RefTypeTag
	RefTypeRemote
	RefTypeHEAD
	RefTypeOther
)

// Commit represents a Git commit
type Commit struct {
	Hash      string
	Author    Signature
	Committer Signature
	Message   string
	Summary   string
	Body      string
	Parents   []string
	Tree      string
	Stats     *DiffStats
}

// Signature represents author/committer information
type Signature struct {
	Name  string
	Email string
	Time  time.Time
}

// DiffStats represents diff statistics
type DiffStats struct {
	FilesChanged int
	Insertions   int
	Deletions    int
}

// Status represents the working directory status
type Status struct {
	Branch    string
	Ahead     int
	Behind    int
	Staged    []FileStatus
	Modified  []FileStatus
	Untracked []FileStatus
	Conflict  []FileStatus
}

// FileStatus represents the status of a single file
type FileStatus struct {
	Path        string
	X           string // Staged status
	Y           string // Unstaged status
	From        string // For renames
	IsNew       bool
	IsDeleted   bool
	IsModified  bool
	IsRenamed   bool
	IsCopied    bool
	IsConflict  bool
	IsUntracked bool
	IsAdded     bool
}

// File represents a file in the repository
type File struct {
	Path    string
	Mode    os.FileMode
	Size    int64
	IsDir   bool
	IsBinary bool
}

// Stash represents a Git stash
type Stash struct {
	Index   int
	Message string
	Branch  string
	Commit  *Commit
}

// LogOptions represents options for log queries
type LogOptions struct {
	MaxCount int
	Skip     int
	Branch   string
	Path     string
	All      bool
	Reverse  bool
}

// DiffOptions represents options for diff operations
type DiffOptions struct {
	ContextLines int
	IgnoreSpace  bool
	IgnoreCase   bool
	Paths        []string
}

// CommitOptions represents options for commit operations
type CommitOptions struct {
	All      bool // Automatically stage modified/deleted files
	Amend    bool // Amend the previous commit
	Signoff  bool // Add Signed-off-by line
	Author   *Signature // Override commit author
	Committer *Signature // Override commit committer
}

// GoGitClient implements the Client interface using go-git
type GoGitClient struct {
	path string
	repo *git.Repository
}

// NewClient creates a new Git client
func NewClient() Client {
	return &GoGitClient{}
}

// Open opens a Git repository at the given path
func (c *GoGitClient) Open(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	repo, err := git.PlainOpen(absPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	c.path = absPath
	c.repo = repo
	return nil
}

// GetRepository returns the underlying git repository
func (c *GoGitClient) GetRepository() (*Repository, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}
	return &Repository{Path: c.path, repo: c.repo}, nil
}

// GetWorktree returns the worktree
func (c *GoGitClient) GetWorktree() (*Worktree, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}
	wt, err := c.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}
	return &Worktree{wt: wt}, nil
}

// IsRepository checks if the current directory is a Git repository
func (c *GoGitClient) IsRepository() bool {
	_, err := git.PlainOpen("")
	return err == nil
}

// GetHead returns the HEAD reference
func (c *GoGitClient) GetHead() (*Ref, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	head, err := c.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	return &Ref{
		Name: head.Name().String(),
		Type: RefTypeHEAD,
		Hash: head.Hash().String(),
	}, nil
}

// GetBranches returns all branches
func (c *GoGitClient) GetBranches() ([]*Ref, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	branches, err := c.repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	var result []*Ref
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		result = append(result, &Ref{
			Name: ref.Name().String(),
			Type: RefTypeBranch,
			Hash: ref.Hash().String(),
		})
		return nil
	})

	return result, err
}

// GetTags returns all tags
func (c *GoGitClient) GetTags() ([]*Ref, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	tags, err := c.repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	var result []*Ref
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		result = append(result, &Ref{
			Name: ref.Name().String(),
			Type: RefTypeTag,
			Hash: ref.Hash().String(),
		})
		return nil
	})

	return result, err
}

// GetRemotes returns all remotes
func (c *GoGitClient) GetRemotes() ([]*Remote, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	remotes, err := c.repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("failed to get remotes: %w", err)
	}

	var result []*Remote
	for _, remote := range remotes {
		result = append(result, &Remote{
			Name: remote.Config().Name,
			URLs: remote.Config().URLs,
		})
	}

	return result, nil
}

// Remote represents a Git remote
type Remote struct {
	Name string
	URLs []string
}

// GetCommit returns a single commit by hash
func (c *GoGitClient) GetCommit(hash string) (*Commit, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	commitHash := plumbing.NewHash(hash)
	commit, err := c.repo.CommitObject(commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	return c.commitToModel(commit)
}

// GetCommits returns commits based on the given options
func (c *GoGitClient) GetCommits(opts *LogOptions) ([]*Commit, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	var head plumbing.Hash
	if opts.Branch != "" {
		ref, err := c.repo.Reference(plumbing.ReferenceName(opts.Branch), true)
		if err != nil {
			return nil, fmt.Errorf("failed to get branch reference: %w", err)
		}
		head = ref.Hash()
	} else {
		ref, err := c.repo.Head()
		if err != nil {
			return nil, fmt.Errorf("failed to get HEAD: %w", err)
		}
		head = ref.Hash()
	}

	logOptions := &git.LogOptions{
		From:  head,
		Order: git.LogOrderCommitterTime,
	}

	if opts.Path != "" {
		logOptions.PathFilter = func(path string) bool {
			return path == opts.Path
		}
	}

	commits, err := c.repo.Log(logOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	var result []*Commit
	count := 0
	err = commits.ForEach(func(commit *object.Commit) error {
		if opts.Skip > 0 {
			opts.Skip--
			return nil
		}

		if opts.MaxCount > 0 && count >= opts.MaxCount {
			return fmt.Errorf("max count reached")
		}

		commitModel, err := c.commitToModel(commit)
		if err != nil {
			return err
		}

		result = append(result, commitModel)
		count++
		return nil
	})

	if err != nil && err.Error() != "max count reached" {
		return nil, err
	}

	if opts.Reverse {
		// Reverse the slice
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}
	}

	return result, nil
}

// GetLogCount returns the total number of commits
func (c *GoGitClient) GetLogCount() (int, error) {
	if c.repo == nil {
		return 0, fmt.Errorf("repository not opened")
	}

	ref, err := c.repo.Head()
	if err != nil {
		return 0, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commits, err := c.repo.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get commits: %w", err)
	}

	count := 0
	err = commits.ForEach(func(*object.Commit) error {
		count++
		return nil
	})

	return count, err
}

// GetStatus returns the working directory status
func (c *GoGitClient) GetStatus() (*Status, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	worktree, err := c.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	// Get branch information
	branch := ""
	head, err := c.repo.Head()
	if err == nil {
		branch = head.Name().Short()
	}

	// Calculate ahead/behind
	ahead, behind := 0, 0
	if branch != "" && branch != "HEAD" {
		remoteRef, err := c.repo.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", branch)), true)
		if err == nil {
			// This is a simplified calculation - in a real implementation,
			// you'd want to use more sophisticated logic
			_ = remoteRef
		}
	}

	result := &Status{
		Branch: branch,
		Ahead:  ahead,
		Behind: behind,
	}

	for path, fileStatus := range status {
		file := FileStatus{
			Path: path,
			X:    string(fileStatus.Staging),
			Y:    string(fileStatus.Worktree),
		}

		// Set flags based on status
		switch {
		case fileStatus.Staging == 'A':
			file.IsNew = true
			result.Staged = append(result.Staged, file)
		case fileStatus.Staging == 'D':
			file.IsDeleted = true
			result.Staged = append(result.Staged, file)
		case fileStatus.Staging == 'M':
			file.IsModified = true
			result.Staged = append(result.Staged, file)
		case fileStatus.Staging == 'R':
			file.IsRenamed = true
			result.Staged = append(result.Staged, file)
		case fileStatus.Staging == 'C':
			file.IsCopied = true
			result.Staged = append(result.Staged, file)
		}

		switch {
		case fileStatus.Worktree == '?':
			result.Untracked = append(result.Untracked, file)
		case fileStatus.Worktree == 'M':
			file.IsModified = true
			result.Modified = append(result.Modified, file)
		case fileStatus.Worktree == 'D':
			file.IsDeleted = true
			result.Modified = append(result.Modified, file)
		case fileStatus.Worktree == 'A':
			file.IsNew = true
			result.Modified = append(result.Modified, file)
		}

		// Handle conflicts
		if fileStatus.Staging == 'U' || fileStatus.Worktree == 'U' {
			file.IsConflict = true
			result.Conflict = append(result.Conflict, file)
		}
	}

	return result, nil
}

// GetDiff returns the diff for the given path
func (c *GoGitClient) GetDiff(path string) (*Diff, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	worktree, err := c.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	// This is a simplified implementation
	// In a real implementation, you'd use git.Diff
	_ = worktree
	_ = path

	return &Diff{}, nil
}

// GetFiles returns files in the given path
func (c *GoGitClient) GetFiles(path string) ([]*File, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	// This is a simplified implementation
	return []*File{}, nil
}

// GetStashes returns all stashes
func (c *GoGitClient) GetStashes() ([]*Stash, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("repository not opened")
	}

	// This is a simplified implementation
	return []*Stash{}, nil
}

// StageFile stages a single file
func (c *GoGitClient) StageFile(path string) error {
	if c.repo == nil {
		return fmt.Errorf("repository not opened")
	}

	worktree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	_, err = worktree.Add(path)
	if err != nil {
		return fmt.Errorf("failed to stage file %s: %w", path, err)
	}

	return nil
}

// UnstageFile unstages a single file
func (c *GoGitClient) UnstageFile(path string) error {
	if c.repo == nil {
		return fmt.Errorf("repository not opened")
	}

	worktree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Use reset to unstage files
	_, err = worktree.Add(path)
	if err != nil {
		return fmt.Errorf("failed to unstage file %s: %w", path, err)
	}

	return nil
}

// StageAll stages all changes
func (c *GoGitClient) StageAll() error {
	if c.repo == nil {
		return fmt.Errorf("repository not opened")
	}

	worktree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	_, err = worktree.Add(".")
	if err != nil {
		return fmt.Errorf("failed to stage all files: %w", err)
	}

	return nil
}

// UnstageAll unstages all changes
func (c *GoGitClient) UnstageAll() error {
	if c.repo == nil {
		return fmt.Errorf("repository not opened")
	}

	// Use git reset to unstage all files
	_, err := c.ExecuteCommand("reset", "HEAD", ".")
	return err
}

// DiscardChanges discards changes to a file
func (c *GoGitClient) DiscardChanges(path string) error {
	if c.repo == nil {
		return fmt.Errorf("repository not opened")
	}

	// Use git checkout to discard changes
	_, err := c.ExecuteCommand("checkout", "--", path)
	return err
}

// Commit creates a new commit
func (c *GoGitClient) Commit(message string, opts *CommitOptions) error {
	if c.repo == nil {
		return fmt.Errorf("repository not opened")
	}

	worktree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	commitOptions := &git.CommitOptions{}
	if opts != nil {
		commitOptions.All = opts.All
		// Note: go-git doesn't support amend, signoff, author override directly
		// These would need to be implemented separately
	}

	_, err = worktree.Commit(message, commitOptions)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// GetRootPath returns the repository root path
func (c *GoGitClient) GetRootPath() string {
	return c.path
}

// GetRelativePath returns the relative path from the repository root
func (c *GoGitClient) GetRelativePath(path string) string {
	rel, _ := filepath.Rel(c.path, path)
	return rel
}

// ExecuteCommand executes a Git command
func (c *GoGitClient) ExecuteCommand(args ...string) ([]byte, error) {
	if c.path == "" {
		return nil, fmt.Errorf("repository path not set")
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = c.path
	return cmd.Output()
}

// commitToModel converts a go-git commit to our Commit model
func (c *GoGitClient) commitToModel(commit *object.Commit) (*Commit, error) {
	// Split message into summary and body
	message := commit.Message
	summary := message
	body := ""
	
	if idx := strings.Index(message, "\n"); idx >= 0 {
		summary = message[:idx]
		body = strings.TrimSpace(message[idx+1:])
	}

	// Calculate diff stats (simplified)
	stats := &DiffStats{}

	commitModel := &Commit{
		Hash:    commit.Hash.String(),
		Message: message,
		Summary: summary,
		Body:    body,
		Tree:    commit.TreeHash.String(),
		Stats:   stats,
		Author: Signature{
			Name:  commit.Author.Name,
			Email: commit.Author.Email,
			Time:  commit.Author.When,
		},
		Committer: Signature{
			Name:  commit.Committer.Name,
			Email: commit.Committer.Email,
			Time:  commit.Committer.When,
		},
	}

	// Add parents
	for _, parent := range commit.ParentHashes {
		commitModel.Parents = append(commitModel.Parents, parent.String())
	}

	return commitModel, nil
}

// Diff represents a diff

type Diff struct {
	Files []*DiffFile
}

// DiffFile represents a file diff
type DiffFile struct {
	OldPath string
	NewPath string
	OldMode os.FileMode
	NewMode os.FileMode
	IsNew   bool
	IsDeleted bool
	IsRenamed bool
	IsCopied  bool
	IsBinary  bool
	Hunks     []*DiffHunk
}

// DiffHunk represents a diff hunk
type DiffHunk struct {
	OldStart int
	OldLines int
	NewStart int
	NewLines int
	Lines    []*DiffLine
}

// DiffLine represents a diff line
type DiffLine struct {
	Type    DiffLineType
	Content string
	OldLine int
	NewLine int
}

// DiffLineType represents the type of diff line
type DiffLineType int

const (
	DiffLineContext DiffLineType = iota
	DiffLineAddition
	DiffLineDeletion
)

// Helper functions
func convertCommit(commit *object.Commit) (*Commit, error) {
	stats := &DiffStats{
		FilesChanged: 0, // Placeholder
		Insertions:   0, // Placeholder
		Deletions:    0, // Placeholder
	}

	commitModel := &Commit{
		Hash:    commit.Hash.String(),
		Message: commit.Message,
		Summary: strings.Split(commit.Message, "\n")[0],
		Body:    strings.Join(strings.Split(commit.Message, "\n")[1:], "\n"),
		Parents: []string{},
		Tree:    commit.TreeHash.String(),
		Stats:   stats,
		Author: Signature{
			Name:  commit.Author.Name,
			Email: commit.Author.Email,
			Time:  commit.Author.When,
		},
		Committer: Signature{
			Name:  commit.Committer.Name,
			Email: commit.Committer.Email,
			Time:  commit.Committer.When,
		},
	}

	// Add parents
	for _, parent := range commit.ParentHashes {
		commitModel.Parents = append(commitModel.Parents, parent.String())
	}

	return commitModel, nil
}

func convertStatus(gitStatus git.Status) (*Status, error) {
	status := &Status{
		Branch:    "main", // Placeholder
		Staged:    []FileStatus{},
		Modified:  []FileStatus{},
		Untracked: []FileStatus{},
		Conflict:  []FileStatus{},
	}

	for path, file := range gitStatus {
		fileStatus := FileStatus{
			Path: path,
			X:    string(file.Staging),
			Y:    string(file.Worktree),
		}

		// Determine file status flags
		if file.Worktree == 63 { // Untracked
			fileStatus.IsUntracked = true
			status.Untracked = append(status.Untracked, fileStatus)
		} else if file.Worktree == 77 { // Modified
			fileStatus.IsModified = true
			status.Modified = append(status.Modified, fileStatus)
		}

		// Add to appropriate category based on staging
		if file.Staging != 32 { // Unmodified
			status.Staged = append(status.Staged, fileStatus)
		}
	}

	return status, nil
}

// Repository methods
func (r *Repository) GetCommits(opts *LogOptions) ([]*Commit, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repository not available")
	}
	
	// Return sample commits for now
	return []*Commit{
		{
			Hash:    "abc123def456",
			Message: "Initial commit\n\nThis is the initial commit",
			Summary: "Initial commit",
			Body:    "This is the initial commit",
			Parents: []string{},
			Tree:    "tree123",
			Stats: &DiffStats{
				FilesChanged: 1,
				Insertions:   10,
				Deletions:    0,
			},
			Author: Signature{
				Name:  "Test User",
				Email: "test@example.com",
				Time:  time.Now(),
			},
			Committer: Signature{
				Name:  "Test User",
				Email: "test@example.com",
				Time:  time.Now(),
			},
		},
	}, nil
}

func (r *Repository) GetCommit(hash string) (*Commit, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repository not available")
	}
	
	// Return sample commit
	return &Commit{
		Hash:    hash,
		Message: "Sample commit\n\nThis is a sample commit",
		Summary: "Sample commit",
		Body:    "This is a sample commit",
		Parents: []string{"parent123"},
		Tree:    "tree456",
		Stats: &DiffStats{
			FilesChanged: 1,
			Insertions:   5,
			Deletions:    2,
		},
		Author: Signature{
			Name:  "Sample User",
			Email: "sample@example.com",
			Time:  time.Now(),
		},
		Committer: Signature{
			Name:  "Sample User",
			Email: "sample@example.com",
			Time:  time.Now(),
		},
	}, nil
}

func (r *Repository) GetStatus() (*Status, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repository not available")
	}
	
	// Return sample status
	return &Status{
		Branch: "main",
		Staged: []FileStatus{
			{Path: "README.md", X: "M", Y: " "},
		},
		Modified: []FileStatus{
			{Path: "main.go", X: " ", Y: "M"},
		},
		Untracked: []FileStatus{
			{Path: "newfile.txt", X: "?", Y: "?"},
		},
		Conflict: []FileStatus{},
	}, nil
}

func (r *Repository) GetCommitDiff(commitHash string) (string, error) {
	if r.repo == nil {
		return "", fmt.Errorf("repository not available")
	}
	
	// Return sample diff
	return fmt.Sprintf(`diff --git a/file.txt b/file.txt
index 0000000..1111111 100644
--- a/file.txt
+++ b/file.txt
@@ -1,3 +1,5 @@
 line 1
-line 2
+line 2 modified
 line 3
+line 4 added
+line 5 added

Commit: %s`, commitHash), nil
}

