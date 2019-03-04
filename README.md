# GoZix GoRedis

## Dependencies

* [viper](https://github.com/gozix/viper)

## Configuration example

```json
{
  "redis": {
    "host": "127.0.0.1",
    "port": "6379",
    "db": 0,
    "max_retiries": 2,   
    "read_timeout": "2s",
    "write_timeout": "2s",
    "idle_timeout": "1m"
  }
}
```
