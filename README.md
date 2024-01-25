## GOBrute

A brazingly fast Wordpress login bruteforce tool written in Go. 

### Changelog v2.1:
1. Bug fixes.
2. Performance improvement.
3. Added better user experience.

### Info: 
1. Find usernames **https://target.com/wp-json/wp/v2/users**.
2. Generate passwords related to username using [this tool](https://github.com/anonymouseshrivastav/password-generator).
3. Now you can use this tool.
4. Give a star if its useful.

### Supported in:
Linux <br/>
Termux <br/>
Mac Os <br/>
Windows <br/>


### Installation:

```bash
pkg install golang git
```

```bash
git clone https://github.com/anonymouseshrivastav/gobrute
```

```bash
cd gobrute
```

```bash
go build
```

### Usage:
```bash
gobrute --url <login page url> -u <username> -p <password list> -t <threads>
```
```bash
gobrute --url https://target.com/wp-login.php -u admin -p passwords.txt -t 40
```


### Access the tool from any directory:

```bash
cp gobrute /data/data/com.termux/files/usr/bin/gobrute
```

### More Info:
1. Adding more features.
2. Auther: Anon Shrivastav.
3. Telegram: https://t.me/anonshrivastavofficial.
