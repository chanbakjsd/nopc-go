# How do I contribute to this project?
Here's a quick rundown on how to start contributing to this project!
1. Find something to do/fix.
This can be anything. Maybe an issue from the issue tracker or a thing you want added to the language. Make sure to coordinate with others so they agree with what you're doing AND aren't already working on it to make sure your effort doesn't go to waste.
2. Fork the repository.
3. Pull from latest and start writing your code.
4. Create a branch for your fix.
It's quite important to create a branch as you may have multiple things you want to add at once. You want to have one branch for each thing you want to add. 
5. Work away at your fix.
Make sure to follow the formatting guidelines! Also, add some tests to make sure your code works as intended before pushing it out.
6. Add the relavant files that has been changed.
7. Commit!
Your code's done, you've staged the code. Now just commit it. We require the naming of the commit to be in the format of `<module>: your changes` (e.g. `core: add foo`) to make it easier to navigate in the future. Furthermore, we require commit signing. Read more about commit signatures [here](https://help.github.com/en/articles/signing-commits).
8. Create a pull request.
Let people know you have done your thing and let them review your code. Go to the pull request page of this repository and create a new pull request. Follow the instructions shown!
9. Wait for a review.
It might take a while but eventually, a maintainer will come and read your code to tell you how well you've done. Do fixes as needed by repeating step 5-8.
You'll also want to rebase your commits to make sure only one commit is created for everything that are relevant together. Please `git rebase -i` for that. After the rebase, you might need to add -f to your git push fork command to force git into pushing your commit. You can read more about rebasing [here](https://www.atlassian.com/git/tutorials/rewriting-history/git-rebase).
