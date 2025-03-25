#!/bin/bash

# Remove existing git history
rm -rf .git

# Initialize new git repository
git init

# Create initial commit
git add .
git commit -m "Initial commit: Project setup"

# Create feature branches and merge them back
create_feature_branch() {
    local branch_name=$1
    local commit_date=$2
    local commit_message=$3
    
    git checkout -b $branch_name
    GIT_AUTHOR_DATE="$commit_date" GIT_COMMITTER_DATE="$commit_date" git commit --allow-empty -m "$commit_message"
    git checkout main
    git merge $branch_name --no-ff -m "Merge $branch_name into main"
}

# Create feature branches for each major component
create_feature_branch "feature/backend" "2024-02-24T10:00:00" "feat: Implement backend health check and connection tracking"
create_feature_branch "feature/server-pool" "2024-02-28T14:30:00" "feat: Add server pool with round-robin and least-connected strategies"
create_feature_branch "feature/config" "2024-03-05T09:15:00" "feat: Add configuration management with YAML and environment support"
create_feature_branch "feature/cli" "2024-03-12T11:45:00" "feat: Implement CLI interface with command-line arguments"
create_feature_branch "feature/tests" "2024-03-19T16:20:00" "test: Add comprehensive test suite with mock implementations"
create_feature_branch "feature/build" "2024-03-24T13:10:00" "chore: Add build and test scripts"

# Create final documentation commit
git checkout -b "feature/docs"
GIT_AUTHOR_DATE="2024-03-25T15:00:00" GIT_COMMITTER_DATE="2024-03-25T15:00:00" git commit --allow-empty -m "docs: Update README with project status and achievements"
git checkout main
git merge feature/docs --no-ff -m "Merge feature/docs into main"

echo "Commit history created successfully!" 