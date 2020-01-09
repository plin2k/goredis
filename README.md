# GoZix GoRedis

## Dependencies

* [viper](https://github.com/gozix/viper)

## Configuration example

```json
{
  "redis": {
    "default": {
      "host": "127.0.0.1",
      "port": "6379",
      "db": 0,
      "password": "somepassword",
      "max_retiries": 2,
      "read_timeout": "2s",
      "write_timeout": "2s",
      "idle_timeout": "1m"
    }
  }
}
```
"password" field is optional and ignored if empty
"db" field is optional. Default is 0