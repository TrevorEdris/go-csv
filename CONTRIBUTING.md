# Basic Contribution Instructions

_This document has been adapted from [spacetelescope/jswt](https://github.com/spacetelescope/jwst/blob/main/CONTRIBUTING.md)_.

## Reporting bugs / requesting a new feature

If you would like to report a bug or request a new feature, this can be done by
opening a new [issue](https://github.com/TrevorEdris/go-csv/issues).

## Contributing code

If you would like to contribute code, this is done by submitting a [pull request](https://github.com/TrevorEdris/go-csv/pulls)
to the `main` branch of `TrevorEdris/go-csv`. To do this, it is recommended to follow the
following workflow (which assumes you already have a Github account / command line tools).

### Step 1: Forking and cloning the `go-csv` repository

First, to clarify some terms that will be commonly used here:

* `upstream` refers to the main `go-csv` repository. this is where code from all
	contributors is ultimately merged into and where releases of the package will be made from.
* `origin` refers to the online fork you made of the upstream `go-csv` repository.
* `local` refers to the clone you made of the origin on your computer.
* The term `remote` in this context can refer to origin, or upstream. in general, it means anything hosted online on Github.

The first step is to create your own 'remote' (online) and 'local' (on your machine)
clones of the central `TrevorEdris/go-csv` repository. You will make code changes
on your machine to your 'local' clone, push these to 'origin' (your online fork),
and finally, open a pull request to the 'main' branch of `TrevorEdris/go-csv`.

1. On the 'TrevorEdris/go-csv' Github repository page, 'fork' the go-csv repository
to your own account space by clicking the appropriate button on the upper right-hand
side. This will create an online clone of the main go-csv repository but instead of
being under the 'TrevorEdris' organization, it will be under your personal
account space. It is necessary to fork the repository so that many contributors
aren't making branches directly on 'TrevorEdris/go-csv'. 

2. Now that you have remotely forked `go-csv`, it needs to be downloaded
to your machine. To create this 'local' clone, choose an area on your file system
and use the `git clone` command to dowload your remote fork on to your machine.

		>> cd directory
		>> git clone git@github.com:<your_username>/go-csv.git

3. Make sure that your references to 'origin' and 'upstream' are set correctly - you will
need this to keep everything in sync and push your changes online. While your inital
local clone will be an exact copy of your remote, which is an exact copy of the 'upstream'
`TrevorEdris/go-csv`, these all must be kept in sync manually (via git fetch/pull/push).

	To check the current status of these references:

		>> git remote -v

After your inital clone, you will likely be missing the reference to 'upstream'
(which is just the most commonly used name in git to refer to the main project repository - you
can call this whatever you want but the origin/upstream conventions are most commonly used) - to 
set this, use the `add` git command:

If you are using an SSH key to authenticate.

	>> git remote add upstream git@github.com:TrevorEdris/go-csv.git
Otherwise, you can simply set it to the repository URL but you will have to
enter your password every time you fetch/push

	>> git remote add upstream https://github.com/TrevorEdris/go-csv

If you ever want to reset these URLs, add references to other remote forks of
`go-csv` for collaboration, or change from URL to SSH, you can use the related
`git remote set-url` command.

### Step 2: Creating a branch for your changes

It is a standard practice in git to create a new 'branch' (off `upstream/main`)
for each new feature or bug fix. You can call this branch whatever you like - in
this example, we'll call it 'my_feature'. First, make sure you
have all recent changes to upstream by 'fetching' them:

		>> git fetch upstream

The following will create a new branch off local/main called 'my_feature', and automatically switch you over to your new branch.

		>> git checkout -b my_feature upstream/main

### Step 3: Making code changes

Now that you've forked, cloned, made a new branch for your feature, you are ready
to make changes to the code. As you make changes,
make sure to `git commit -m <"some message">` frequently
(in case you need to undo something by reverting back to a previous commit - you
cant do this if you commit everything at once!). After you've made your desired
changes, and committed these changes, you will need to push them online to your
'remote' fork of `go-csv`:

	>> git push origin my_feature

### Step 4: Opening a pull request

Now, you can open a pull request on the main branch of the upstream `go-csv` repository.

1. On the `TrevorEdris/go-csv` web page, after you push your changes you should
see a large green go-csv appear at the top prompting you to open a pull request
with your recently pushed changes. You can also open a pull request from the
[pull request tab](https://github.com/TrevorEdris/go-csv/pulls) on that page.
Select your fork and your 'my_feature' branch, and open a pull request against
the 'main' branch.

2. There is now a checklist of items that need to be done before your PR can be merged.
	* The continuous integration (CI) tests must complete and pass. The CI
	runs several different checks including running the unit tests and
    checking for code style issues. The CI runs upon opening
	a PR, and will re-run any time you push commits to that branch. 
	* Your PR will need to be reviewed and approved by the repo owner.
	They may require changes from you before your code can be merged, in which
	case you will need to go back and make these changes and push them (they will
		automatically appear in the PR when they're pushed to origin/my_feature).
