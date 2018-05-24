# tgbot-golang
Telegram bot that allows to compile Golang source code with output

### NaCl
This bot uses [Google Native Client](https://developer.chrome.com/native-client) to sandbox compiled binaries.
To use bot, you should compile NaCl and place `sel_ldr_x86_64` binary to the directory with the compiled bot.

### Package goimports
Bot also uses `goimports` command to automatically add imports of packages used in the source code.

To install this package execute following command:
```
$ go get golang.org/x/tools/cmd/goimports
```

### Configuration
To configure bot you need just to copy contained in the repo `config.example.json` file as `config.json`
and change `<YOUR_API_TOKEN>` to your Telegram Bot API token.

Example configuration:
```
{
    "token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
}
```