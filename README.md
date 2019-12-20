# WinSCP.ini decrypt password
Парсим файл конфигурации и расшифровываем сохраненные пароли
Для расшифровки используем проект [anoopengineer/winscppasswd](https://github.com/anoopengineer/winscppasswd)

```sh
winscp-decrypt.exe --config-path=WinSCP.ini --decrypt-path=WinSCP_decrypt.ini
```
