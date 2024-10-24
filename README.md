### Agora - Terminal based API client

Agora is a simple and basic app that lets you build API requests on the terminal. 
Inspired by [lazygit](https://github.com/jesseduffield/lazygit) and [Postman](https://www.postman.com). 
Powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).

This tool is intended for backend developers who are prototyping API servers and want to
quickly and repeatedly test out their endpoints, without ever leaving their terminal or IDE.

If you need a more robust and full feature API client, you may want to use
something like Postman or Bruno.

All data are saved locally in your machine.

![](./assets/demo.gif)

### Installation

```shell
go install github.com/gabrielfu/agora
```

Then, you can launch Agora with the command `agora`.

### Usage

Run `agora` anywhere. Data will be saved in the default directory `$HOME/.agora`.

```shell
agora
```

Alternatively, you can save your data in your project directory to have an isolate environment.

```shell
# this will read and write data to `./.agora`
agora .

# this will read and write data to `/another/project/.agora`
agora /another/project
```

### Features

- [X] Send HTTP requests (only JSON body supported)
- [X] Multiple collections
- [X] All data saved locally
- [X] Support Linux, MacOS and Windows

#### Coming Soon

- [ ] Response history
- [ ] Authentication helper
- [ ] Environments
- [ ] Request timeout
- [ ] File upload
- [ ] Non JSON body