# warden

A CLI for creating encrypted backups. See [architecture docs](/docs/architecture.md) for more info.

## Design

![Crypto Workflow](docs/img/crypto-workflow.png)

## TODO

- [ ] create backup command

## Commands

| Command      | Description                                                   |
| ------------ | ------------------------------------------------------------- |
| init         | Create a new encrypted backup store                           |
| show         | Print resource information (see appendix for valid resources) |
| backup <dir> | Create a new backup of a directory                            |

### Appendix

- valid resources: masterkey, config
