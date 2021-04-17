## Remitly CLI

### Prerequisites
```bash
export REMITLY_PATH=$HOME/.remitly
export REMITLY_PROFILE=default
```

### Building and usage

```bash
make build

./remitly initialize -n $REMITLY_PROFILE --url http://cloud.remitly.io/ --username XXX
./remitly deploy --help # for more flag information
./remitly deploy -a app_name --revision 1.0.0
```

### Assumptions
- I assume, the app was built before and through the `--application` flag, we are passing some pointer to the artifact.
- I assume, the CLI will be installed on user's OS by some package manager i.e. `brew`
- I assume, the exercise's input naming was fluent, because `--version` flag is kind of reserved.

### Known issues
- CLI doesn't allow scaling an already deployed version of the app. As of current state, the CLI will block the deployment.
- Some things that use `os` package are not unit tested. I did not wanted to add complexity by mock wrapping it.
- I'm **not sure** if the CLI works on a non unix operating systems, file paths may be problematic.

### Features (and improvements), that could be added
- Already deployed version scaling.
- Homebrew tap and formula.
- Improve orchestration, right now we create instances then orchestrate them. It could be improved to be more K8s like.
- Concurrency adds complexity, so I wanted to avoid that for now, but it is a good feature to add. (linked to the point above)
- Profile and context managing could be done 100 times better. I did not want to spend too much time on that.
- A flag for optional Load Balancer creation, could be a nice feature
- Support for custom subcommands, like: 
  ```
  remitly install github.com/foo/bar
  remitly foo-bar --some value
  ```
- Logging can always be improved + log level steering.
- I'm not 100% sure about the project structure, never did an CLI before.
- Integration tests.