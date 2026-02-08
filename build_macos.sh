wails build -clean -platform darwin/universal

wails build \
  -clean \
  -platform darwin/universal \
  -ldflags "-s -w" \
  -trimpath
