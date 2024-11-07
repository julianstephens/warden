# warden

A CLI for creating encrypted backups. See [architecture docs](/docs/architecture.md) for more info.

## Design

![Crypto Workflow](docs/img/crypto-workflow.png)

## TODO

- [ ] add ability to create backup chunks/packs
- [ ] add fault-tolerant save for chunks/packs
- [ ] add cache for resuming backups

## Commands

| Command      | Description                                                   |
| ------------ | ------------------------------------------------------------- |
| init         | Create a new encrypted backup store                           |
| show         | Print resource information (see appendix for valid resources) |
| backup <dir> | Create a new backup of a directory                            |

### Appendix

- valid resources: masterkey, config
