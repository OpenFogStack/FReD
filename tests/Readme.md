## Log Level

In order to change the log level, create a `.env`-file in this folder with the following content:

```env
LOG_LEVEL=info #debug,info,warn,error,fatal,panic
LOG_LEVEL_STORE=error #debug,info,warn,error,fatal,panic
```

## Debugging

It is possible to debug nodeB of the 3NodeTest:

- Create some breakpoints
- Run the configuration "3NodeTest: Run Tests & Debug nodeB (start "Debug nodeB" immediately after!)" (equals to `make debug-nodeB` in 3NodeTest)
- Run the configuration "Debug nodeB"
