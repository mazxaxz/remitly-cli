## TBD

```bash
export REMITLY_PATH=$HOME/.remitly
export REMITLY_PROFILE=default

make build

./remitly initialize -n default --url http://cloud.remitly.io/ --username XXX
./remitly deploy -a app_name --revision 1.0.0
```
