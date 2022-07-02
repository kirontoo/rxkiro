# TODOs

- [ ] UPDATE DOCS FOR ENVIRONTMENT VARS

- [ ] TEST EVERYTHING
- [ ] timed messages?
- [ ] command to add a command
- [ ] migrate to `pq` from `postgrest`?
- [ ] BUG: not returning an error message when a command does not exist

## Command Variables
- [ ] `mention` and `user` vars - add a `@` check for users
- [x] error handling when there is no input for cmd var
- [x] ${random 1 100} - random number from some range
- [x] ${random 200} - random number from 0 to 200
- [x] ${random} - any random number (default from 0 to 100000)
- [ ] !song - [spotify api: now playing](https://developer.spotify.com/documentation/web-api/reference/#/operations/get-the-users-currently-playing-track)
- [ ] !cmd options
    - [ ] add
    - [ ] delete
    - [ ] edit
- [ ] !help
- [ ] !timer
    - [ ] start
    - [ ] pause
    - [ ] continue
    - [ ] stop
- [ ] !random list of strings here  - a command that picks a random word from the list

## Misc
- [ ] regex to check all messages for certain terms (i.e wos -> "did someone say wos?")
