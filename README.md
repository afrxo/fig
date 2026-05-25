Git submodules and subtrees can't selectively import a subdirectory from a remote repo and remap it into your project. Submodules always pull the entire repo, and subtrees can prefix but not re-root a subfolder. fig solves this by cloning the repo, extracting only the subdirectory you point at, and copying it into your project. Each package is pinned to a tag, branch, or commit SHA in a lockfile. Run fig sync to pull updates.

## Install with Rokit
```
[tools]
fig = "afrxo/fig@0.3.0"
```

## Build locally

```
git clone https://github.com/afrxo/fig.git
cd fig
go build -o fig .
```

## Auth

Log in before pulling from private repos.

```bash
fig login    # authenticate with your Git host
fig user     # check who you're logged in as
fig logout   # clear stored credentials
```

## Usage

```bash
# Add a package
fig pick physics https://github.com/org/engine.git -p packages/physics -r v2.1.0

# Add another from a different repo
fig pick analytics https://github.com/org/tools.git -p src/analytics

# Sync everything in your lockfile
fig sync

# Sync specific packages
fig sync physics analytics

# Update a package to a new ref
fig pick physics https://github.com/org/engine.git -p packages/physics -r v3.0.0

# Remove a package
fig remove physics

# See what's vendored
fig list
```

Packages are copied into `src/` by default. Use `-o` to change the destination.

## Lockfile

`fig` writes a `fig-lock.yml` file tracking each package's repo, path, ref, and resolved commit SHA. Commit this alongside your vendored code.
