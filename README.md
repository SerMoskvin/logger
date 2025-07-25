# Package "**Logger**"

> [!NOTE]
> This package use "*go.uber.org/zap*" and "*gopkg.in/natefinch/lumberjack.v2*".

Provides logging of specific actions with different logging levels (debug, info, warn,error)

## Work start

### Use default config or create your own:
- In your custom config for every level write:
    - filepath, where logs of this level will be stored;
    - max size of that storage file;
    - max number of backups for this file;
    - max age (in days) for logs to storage in file before rotation; 
    - compress logs files or not (true/false).

- To load default config from the package ("default_config.yml"), you can use this code:
`cfg, err := logger.LoadDefaultConfig()`