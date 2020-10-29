# Status

MVP

## Installation

```bash
curl -L https://github.com/oleg-docler/br/releases/download/0.0.6/br -o br && chmod +x br && ./br
```

## Usage
Show your 'open', 'in progress', 'under review' Jira issues:
```bash
br
```
Checkout to a branch or create a new one based on issues title:

```bash
br [num]
```
[num] - Jira issue number, e.g. JASMIN-1293, 1293 - issue number

If something wrong with the credentials, you can remove the ~/.br/config.json file.
After several incorrect login attempts, you need to login to your jira site first.

## License
[MIT](https://choosealicense.com/licenses/mit/)