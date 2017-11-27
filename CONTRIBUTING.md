# Guidelines for contribution

+ Open an issue if some feature is to be added. Don't drop a pull request without a parent issue.
+ Follow the Template for filing an issue.
+ Patience is key. :)
+ The commit message have to follow a certain pattern :
  - The commit should always be in a patter as issued below. `git commit --amend` is useful for editing commit messages.
    ```
    fileName(s) : short description
     + long description : item 1
     + long description : item 2
    
    Ref: ISSUE#{issuenum}
    ```
    The long description is optional, but should be a list of changes done. 
    There is a blank line separating the body and the Issue Reference line.
  - Possibly try to do the commits in the order of dependency. Say a newly added function is needed by a `file1`
    which is in `file2`, always commit `file2` first and then `file1`.
  - If the above is not possible, say changes in every file is associated with a single feature, the feature should be added in
    the short description, followed by description of changes in each file.
+ `lint`, `fmt`, `vet`, `ineffassign`, `gocyclo`, `misspell` the code you wrote.