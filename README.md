# Sovorem.am CLI

[Sovorem.am](https://sovorem.am/)-ի պաշտոնական command-line գործիքը՝ հայկական backend ծրագրավորման ուսումնական հարթակի համար։ Ուսանողները օգտագործում են էս CLI-ը՝ դասի test-երը լոկալ run անելու և ստորագրված արդյունքները Sovorem.am backend-ին submit անելու համար։

## Բովանդակություն

- [Installation (Տեղադրում)](#installation-տեղադրում)
  - [1. Install Go (Տեղադրել Go-ն)](#1-install-go-տեղադրել-go-ն)
  - [2. Install the Sovorem CLI (Տեղադրել Sovorem CLI-ը)](#2-install-the-sovorem-cli-տեղադրել-sovorem-cli-ը)
  - [3. Login to the CLI (Login լինել CLI-ում)](#3-login-to-the-cli-login-լինել-cli-ում)
- [Usage (Օգտագործում)](#usage-օգտագործում)
- [Configuration (Configuration)](#configuration-configuration)
  - [Base URL for HTTP tests (Base URL HTTP test-երի համար)](#base-url-for-http-tests-base-url-http-test-երի-համար)
  - [CLI colors (CLI-ի գույները)](#cli-colors-cli-ի-գույները)
  - [Troubleshooting the Config (Config-ի խնդիրների լուծում)](#troubleshooting-the-config-config-ի-խնդիրների-լուծում)
- [Upgrading (Update անելը)](#upgrading-update-անելը)
  - [Troubleshooting Upgrading (Update-ի խնդիրների լուծում)](#troubleshooting-upgrading-update-ի-խնդիրների-լուծում)

## Installation (Տեղադրում)

### 1. Install Go (Տեղադրել Go-ն)

Sovorem CLI-ից օգտվելու համար քեզ պետք ա համակարգչիդ վրա տեղադրված Golang-ի թարմ version։

Դասընթացների մեծ մասը նախատեսված են Linux-ի կամ macOS-ի համար (կամ Linux-ը Windows-ում՝ WSL-ի միջոցով)։ Եթե Windows ես օգտագործում, սովորաբար ավելի լավ ա ընտրել WSL-ը և հետևել ներքևի Linux-ի հրահանգներին։ Որոշ դասեր նաև Windows/PowerShell-ով են աշխատում։

**Տարբերակ 1 (Linux/WSL/macOS).** [Webi installer-ը](https://webinstall.dev/golang/) ամենապարզ ձևն ա.

```sh
curl -sS https://webi.sh/golang | sh
```

**Տարբերակ 2 (ցանկացած ՕՀ, ներառյալ Windows/PowerShell).** Օգտվիր [Golang-ի պաշտոնական տեղադրման հրահանգներից](https://go.dev/doc/install)։

Go-ն տեղադրելուց հետո բացիր նոր terminal-ի պատուհան ու run արա `go version`՝ համոզվելու համար, որ ամեն ինչ աշխատում ա։

### 2. Install the Sovorem CLI (Տեղադրել Sovorem CLI-ը)

```sh
go install github.com/sovorem/sovorem@latest
```

Run արա `sovorem --version`՝ տեղադրումը ստուգելու համար։

Եթե ստանում ես "command not found" error-ը, ապա ավելացրու Go-ի bin directory-ն քո `PATH`-ում (սովորաբար `$HOME/go/bin`-ն ա).

```sh
# Linux/WSL
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

```sh
# macOS
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.zshrc
```

```sh
# fish
fish_add_path $HOME/go/bin
```

### 3. Login to the CLI (Login լինել CLI-ում)

Run արա `sovorem login`՝ քո Sovorem.am account-ով մուտք գործելու (authenticate լինելու) համար։ Դրանից հետո արդեն պատրաստ ես run ու submit անելու դասերը։

## Usage (Օգտագործում)

| Command | Description (Նկարագրություն) |
|---------|-------------|
| `sovorem login` | Մուտք գործել քո Sovorem.am account-ով (authenticate լինել) |
| `sovorem logout` | Logout լինել CLI-ից (անջատել account-ը) |
| `sovorem status` | Ցույց տալ login-ի և version-ի status-ը |
| `sovorem run UUID` | Run անել դասի test-երը՝ առանց submit անելու |
| `sovorem run UUID -s` | Run անել test-երը և միանգամից submit անել |
| `sovorem submit UUID` | Run անել test-երը և արդյունքները submit անել Sovorem.am |
| `sovorem config base_url URL` | Override անել HTTP test-երի base URL-ը |
| `sovorem upgrade` | Տեղադրել CLI-ի ամենավերջին version-ը |

Դասի UUID-ները երևում են [sovorem.am](https://sovorem.am)-ի համապատասխան դասի էջում։

## Configuration (Configuration)

CLI-ն պահում ա իր settings-ը `~/.sovorem.yaml` ֆայլում, կամ `$XDG_CONFIG_HOME/sovorem/config.yaml`-ում, եթե `XDG_CONFIG_HOME`-ը սահմանված ա։

Բոլոր command-երը աջակցում են `-h`/`--help` flag-երը։

### Base URL for HTTP tests (Base URL HTTP test-երի համար)

HTTP test-եր ունեցող դասերի համար կարող ես սահմանել base URL, որը կփոխարինի (override կանի) դասի default արժեքին։ Սա պետք ա գալիս, երբ քո լոկալ server-ը աշխատում ա ուրիշ port-ի վրա։

```sh
sovorem config base_url http://localhost:8080/
sovorem config base_url
sovorem config base_url --reset
```

URL-ի մեջ ներառիր նաև protocol scheme-ը (`http://` կամ `https://`)։

### CLI colors (CLI-ի գույները)

Կարող ես սահմանել քո նախընտրած գույները terminal-ի output-ի համար (հաջողված, սխալ, լրացուցիչ տեքստ).

```sh
sovorem config colors --red VALUE --green VALUE --gray VALUE
sovorem config colors
sovorem config colors --reset
```

Որպես `VALUE` օգտագործիր [ANSI color code](https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit) կամ hex տող։

### Config-ի խնդիրների լուծում (Troubleshooting)

Configuration-ը ամբողջությամբ reset անելու համար ջնջիր config ֆայլը։ CLI-ն հաջորդ run-ի ժամանակ ավտոմատ կստեղծի նոր ու թարմ ֆայլ։ Դրանից հետո նորից run արա `sovorem login`։

## Upgrading (Update անելը)

CLI-ն ստուգում ա update-ների առկայությունը login լինելիս և մինչև login պահանջող command-եր run անելը։

```sh
sovorem upgrade
```

Կամ կարող ես տեղադրել կոնկրետ version.

```sh
go install github.com/sovorem/sovorem@v0.1.0
```

### Update-ի խնդիրների լուծում (Troubleshooting)

**Bypass արա proxy-ն**, եթե անընդհատ տեսնում ես նույն upgrade message-ը.

```sh
GOPROXY=direct go install github.com/sovorem/sovorem@latest
```

**Նորից տեղադրիր (Reinstall)**, եթե դա չօգնեց.

```sh
rm "$(which sovorem)"
go install github.com/sovorem/sovorem@latest
sovorem login
```

## Development (Մշակում)

```sh
git clone https://github.com/sovorem/sovorem.git
cd sovorem
go test ./...
go build -o sovorem .
```

## License (Լիցենզիա)

Տես [LICENSE](LICENSE)-ը։
