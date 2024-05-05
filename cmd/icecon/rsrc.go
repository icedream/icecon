package main

//go:generate go run -mod=mod github.com/josephspurrier/goversioninfo/cmd/goversioninfo -manifest "../../rsrc/app.manifest" -icon "../../rsrc/app.ico" -o "rsrc_windows_386.syso"
//go:generate go run -mod=mod github.com/josephspurrier/goversioninfo/cmd/goversioninfo -manifest "../../rsrc/app.manifest" -icon "../../rsrc/app.ico" -o "rsrc_windows_amd64.syso" -64
