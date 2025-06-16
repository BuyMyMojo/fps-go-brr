## 0.5.1 (2025-06-17)

### Fix

- **csv output**: add missing frame width and height headers

## 0.5.0 (2025-06-16)

### Feat

- **main**: add a ton of features - progress bar - frame res measurements - verbosity flag
- add gox build script

### Refactor

- **README**: add resdet info
- **README**: add new lines before lists to fit markdown standard

## 0.4.0 (2025-06-15)

## 0.3.0 (2025-06-15)

### Refactor

- change variable and function names to fit with golang conventions
- remove pointless else statments
- add package string
- **gitignore**: ignore compact build binary

## 0.2.0 (2025-06-14)

### Feat

- **LICENSE**: add dual licenses

## 0.1.1 (2025-06-14)

### Feat

- **workflow**: run on normal commits too

### Fix

- **workflow**: fix artifact upload
- **workflow**: remove sudo again
- **workflow**: fix the upx install
- **workflow**: fix release workflow

## 0.1.0 (2025-06-14)

### Feat

- **forgejo**: add a workflow for building binaries
- **build-compact.sh**: add build script that outputs a heavily compacted binary
- massive additions to output csv
- Get claude AI to generate and debug FPS measurement functionality
- allow selecting diff checker type and min diff
- first wave of tests

### Fix

- **gitignore**: add compiled binary to gitignore
