package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)
	assert.Implements(t, (*Client)(nil), client)
}

func TestGoGitClient(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	
	// Initialize a git repository
	client := NewClient()
	
	// Test IsRepository with non-git directory
	assert.False(t, client.IsRepository())
	
	// Test with current directory (should work if it's a git repo)
	if client.IsRepository() {
		repo, err := client.GetRepository()
		assert.NoError(t, err)
		assert.NotNil(t, repo)
	}
}

func TestRefType(t *testing.T) {
	assert.Equal(t, RefTypeBranch, RefType(0))
	assert.Equal(t, RefTypeTag, RefType(1))
	assert.Equal(t, RefTypeRemote, RefType(2))
	assert.Equal(t, RefTypeHEAD, RefType(3))
}

func TestSignature(t *testing.T) {
	sig := Signature{
		Name:  "Test User",
		Email: "test@example.com",
		Time:  time.Now(),
	}
	
	assert.Equal(t, "Test User", sig.Name)
	assert.Equal(t, "test@example.com", sig.Email)
	assert.False(t, sig.Time.IsZero())
}

func TestCommit(t *testing.T) {
	commit := &Commit{
		Hash:    "abc123",
		Message: "Test commit message\n\nThis is the body",
		Summary: "Test commit message",
		Body:    "This is the body",
		Parents: []string{"parent1", "parent2"},
		Tree:    "tree123",
		Stats:   &DiffStats{FilesChanged: 1, Insertions: 10, Deletions: 5},
		Author: Signature{
			Name:  "Author",
			Email: "author@example.com",
			Time:  time.Now(),
		},
		Committer: Signature{
			Name:  "Committer",
			Email: "committer@example.com",
			Time:  time.Now(),
		},
	}
	
	assert.Equal(t, "abc123", commit.Hash)
	assert.Equal(t, "Test commit message", commit.Summary)
	assert.Equal(t, "This is the body", commit.Body)
	assert.Len(t, commit.Parents, 2)
	assert.Equal(t, "tree123", commit.Tree)
	assert.Equal(t, 1, commit.Stats.FilesChanged)
	assert.Equal(t, 10, commit.Stats.Insertions)
	assert.Equal(t, 5, commit.Stats.Deletions)
}

func TestStatus(t *testing.T) {
	status := &Status{
		Branch: "main",
		Ahead:  2,
		Behind: 1,
		Staged: []FileStatus{
			{Path: "file1.txt", X: "M", Y: " ", IsModified: true},
		},
		Modified: []FileStatus{
			{Path: "file2.txt", X: " ", Y: "M", IsModified: true},
		},
		Untracked: []FileStatus{
			{Path: "file3.txt", X: "?", Y: "?"},
		},
		Conflict: []FileStatus{
			{Path: "file4.txt", X: "U", Y: "U", IsConflict: true},
		},
	}
	
	assert.Equal(t, "main", status.Branch)
	assert.Equal(t, 2, status.Ahead)
	assert.Equal(t, 1, status.Behind)
	assert.Len(t, status.Staged, 1)
	assert.Len(t, status.Modified, 1)
	assert.Len(t, status.Untracked, 1)
	assert.Len(t, status.Conflict, 1)
}

func TestFileStatus(t *testing.T) {
	file := FileStatus{
		Path:       "test.go",
		X:          "M",
		Y:          "D",
		IsModified: true,
		IsDeleted:  true,
	}
	
	assert.Equal(t, "test.go", file.Path)
	assert.Equal(t, "M", file.X)
	assert.Equal(t, "D", file.Y)
	assert.True(t, file.IsModified)
	assert.True(t, file.IsDeleted)
}

func TestLogOptions(t *testing.T) {
	opts := &LogOptions{
		MaxCount: 10,
		Skip:     5,
		Branch:   "main",
		Path:     "src/",
		All:      true,
		Reverse:  false,
	}
	
	assert.Equal(t, 10, opts.MaxCount)
	assert.Equal(t, 5, opts.Skip)
	assert.Equal(t, "main", opts.Branch)
	assert.Equal(t, "src/", opts.Path)
	assert.True(t, opts.All)
	assert.False(t, opts.Reverse)
}

func TestDiffOptions(t *testing.T) {
	opts := &DiffOptions{
		ContextLines: 3,
		IgnoreSpace:  true,
		IgnoreCase:   false,
		Paths:        []string{"file1.txt", "file2.txt"},
	}
	
	assert.Equal(t, 3, opts.ContextLines)
	assert.True(t, opts.IgnoreSpace)
	assert.False(t, opts.IgnoreCase)
	assert.Len(t, opts.Paths, 2)
}

func TestRemote(t *testing.T) {
	remote := &Remote{
		Name: "origin",
		URLs: []string{"https://github.com/user/repo.git"},
	}
	
	assert.Equal(t, "origin", remote.Name)
	assert.Len(t, remote.URLs, 1)
	assert.Equal(t, "https://github.com/user/repo.git", remote.URLs[0])
}

func TestGetRelativePath(t *testing.T) {
	client := &GoGitClient{path: "/home/user/project"}
	
	// Test when path is within repo
	rel := client.GetRelativePath("/home/user/project/src/main.go")
	assert.Equal(t, "src/main.go", rel)
	
	// Test when path is same as repo
	rel = client.GetRelativePath("/home/user/project")
	assert.Equal(t, ".", rel)
	
	// Test when path is outside repo
	rel = client.GetRelativePath("/home/user/other/file.txt")
	assert.Equal(t, "../other/file.txt", rel)
}

func TestExecuteCommand(t *testing.T) {
	// Skip this test in CI environments
	if testing.Short() {
		t.Skip("skipping command execution test in short mode")
	}
	
	client := &GoGitClient{path: "."}
	
	// Test with simple command
	output, err := client.ExecuteCommand("version")
	if err != nil {
		// Git might not be available
		t.Skipf("git not available: %v", err)
	}
	
	assert.Contains(t, string(output), "git version")
}

// TestOpenWithNonGitDirectory tests opening a non-git directory
func TestOpenWithNonGitDirectory(t *testing.T) {
	client := &GoGitClient{}
	
	tempDir := t.TempDir()
	err := client.Open(tempDir)
	assert.Error(t, err)
}

// TestErrorCases tests various error cases
func TestErrorCases(t *testing.T) {
	client := &GoGitClient{}
	
	// Test operations without repository
	_, err := client.GetRepository()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository not opened")
	
	_, err = client.GetBranches()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository not opened")
	
	_, err = client.GetCommits(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository not opened")
	
	_, err = client.GetStatus()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository not opened")
}